package main

import (
    "os"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/davecgh/go-spew/spew"
    "github.com/blevesearch/bleve"

    "github.com/stkim1/backend/model"
    "github.com/stkim1/backend/control"
    "github.com/stkim1/backend/config"
    "github.com/stkim1/backend/storage/boltbk"
    "github.com/stkim1/backend/storage"
    "github.com/stkim1/backend/update"
    pocketsearch "github.com/stkim1/backend/search"
)

func githubSortSupplementInfo(suppDB storage.Nosql, repoModel *model.Repository) error {
    var (
        repoID string                       = repoModel.RepoId
        repoURL string                      = repoModel.RepoPage

        repoSupp model.RepoSupplement
        //langs model.ListLanguage
        //releases model.ListRelease
        //tags model.ListTag
        err error
    )

    // URL CHECK
    if len(repoURL) == 0 {
        return errors.Errorf("Cannot begin update a repo with empty URL")
    }

    suppDB.AcquireLock(repoID, time.Second)
    err = suppDB.GetObj([]string{model.RepoSuppBucket}, repoID, &repoSupp)
    suppDB.ReleaseLock(repoID)
    if err != nil {
        return err
    }

    repoSupp.BuildRecentPublication(15)
    log.Info(spew.Sdump(repoSupp.RecentPublish))

    // save it to database
    suppDB.AcquireLock(repoID, time.Second)
    err = suppDB.UpsertObj([]string{model.RepoSuppBucket}, repoID, &repoSupp, storage.Forever)
    suppDB.ReleaseLock(repoID)
    if err != nil {
        log.Error(err.Error())
    }

    return nil
}

func main() {
    var (
        handleUpdate bool = true
    )
    // config
    cfgPath, ok := os.LookupEnv(config.EnvConfigFilePath)
    if !ok {
        cfgPath = "config.yaml"
    }
    cfg, err := config.NewConfig(cfgPath)

    // (META DB) database
    metaDB, err := gorm.Open(cfg.Database.DatabaseType, cfg.Database.DatabasePath)
    if err != nil {
        log.Fatal(errors.WithStack(err))
        return
    }
    defer metaDB.Close()

    // (BOLTDB) supplementary
    suppDB, err := boltbk.New(cfg.Supplement.DatabasePath)
    if err != nil {
        log.Fatal(errors.WithStack(err))
        return
    }
    defer suppDB.Close()

    // (SEARCH INDEX)
    sIndex, err := bleve.Open(cfg.Search.IndexStoragePath)
    if err != nil {
        m, err := pocketsearch.BuildIndexMapping()
        if err != nil {
            log.Fatal(err)
        }
        sIndex, err = bleve.New(cfg.Search.IndexStoragePath, m)
        if err != nil {
            log.Fatal(err)
        }
    }
    defer sIndex.Close()

    // controller
    ctrl := control.NewController(cfg)

    // update start
    log.Info("Update process started at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
    /*
        client := &http.Client{
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
                Dial: (&net.Dialer{
                    Timeout:   60 * time.Second,
                    KeepAlive: 60 * time.Second,
                }).Dial,
                // We use ABSURDLY large keys, and should probably not.
                TLSHandshakeTimeout:   60 * time.Second,
                ResponseHeaderTimeout: 60 * time.Second,
                ExpectContinueTimeout: 1 * time.Second,
            },
        }
    */
    var repos []model.Repository
    metaDB.Find(&repos)

    // update meta-data
    if false {
        for _, repo := range repos {
            _, err := update.UpdateRepoMeta(metaDB, sIndex, ctrl, &repo)
            if err != nil {
                log.Error(err.Error())
            }
        }
    } else if handleUpdate {
        var repoCount int = len(repos)
        for i, repo := range repos {
            log.Infof("%d / %d | %s - %s", i, repoCount, repo.RepoId, repo.RepoPage)
            resp, err := update.GithubSupplementInfo(suppDB, ctrl, &repo, cfg);
            if err != nil {
                log.Error(err.Error())
            }

            if resp != nil {
                log.Infof("Remaning API limit %d\n", resp.Rate.Remaining)
                if resp.Rate.Remaining < 100 {
                    log.Info("API limit is met")
                    break
                }
            }
        }
    } else {
        for _, repo := range repos {
            githubSortSupplementInfo(suppDB, &repo)
        }
    }
    log.Info("Update process ended at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
}
