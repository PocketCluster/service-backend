package main

import (
    "errors"
    "os"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/google/go-github/github"
    "github.com/davecgh/go-spew/spew"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/control"
    "github.com/stkim1/BACKEND/config"
    "github.com/stkim1/BACKEND/storage/boltbk"
    "github.com/stkim1/BACKEND/storage"
)

func GithubSupplementInfo(suppDB storage.Nosql, ctrl *control.Controller, repoModel *model.Repository) (*github.Response, error) {
    var (
        repoID string                       = repoModel.RepoId
        repoURL string                      = repoModel.RepoPage

        repoSupp model.RepoSupplement
        langs model.ListLanguage
        releases model.ListRelease
        tags model.ListTag

        resp *github.Response
        err error
    )

    // URL CHECK
    if len(repoURL) == 0 {
        return nil, trace.Wrap(errors.New("Cannot begin update a repo with empty URL"))
    }

    suppDB.AcquireLock(repoID, time.Second)
    err = suppDB.GetObj([]string{model.RepoSuppBucket}, repoID, &repoSupp)
    suppDB.ReleaseLock(repoID)
    if err != nil {
        // we don't work on an empty container
        repoSupp = model.RepoSupplement{RepoID:repoID}
    } else {
        //log.Info(spew.Sdump(repoSupp))
        if !repoSupp.Updated.IsZero() && time.Now().Sub(repoSupp.Updated) < (time.Hour * 6) {
            log.Infof("%s :: updated already", repoID)
            return nil, nil
        }
    }

    // get languages
    langs, resp, err = ctrl.GetGithubRepoLanguages(repoURL)
    if err != nil {
        return resp, trace.Wrap(err)
    } else {
        repoSupp.Languages = langs
    }

    // get releases
    releases, _, resp, err = ctrl.GetGithubAllReleases(repoURL, &repoSupp.Releases, 30)
    if err != nil {
        return resp, trace.Wrap(err)
    } else if len(releases) != 0 {
        repoSupp.Releases = releases
    }

    // get tags
    tags, _, resp, err = ctrl.GetGithubAllTags(repoURL, &repoSupp.Tags, 31)
    if err != nil {
        return resp, trace.Wrap(err)
    } else if len(tags) != 0 {
        repoSupp.Tags = tags
    }

    repoSupp.BuildRecentPublication(15)
    repoSupp.Updated = time.Now()
    //log.Info(spew.Sdump(repoSupp))

    // save it to database
    log.Infof("%s - %s :: Lang [%d], Releases [%d] Tags [%d]", repoID, repoURL, len(repoSupp.Languages), len(repoSupp.Releases), len(repoSupp.Tags))
    suppDB.AcquireLock(repoID, time.Second)
    err = suppDB.UpsertObj([]string{model.RepoSuppBucket}, repoID, &repoSupp, storage.Forever)
    suppDB.ReleaseLock(repoID)
    if err != nil {
        log.Error(err.Error())
    }

    return resp, nil
}


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
        return trace.Wrap(errors.New("Cannot begin update a repo with empty URL"))
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

    // database
    repoDB, err := gorm.Open(cfg.Database.DatabaseType, cfg.Database.DatabasePath)
    if err != nil {
        log.Fatal(trace.Wrap(err))
        return
    }
    // (BOLTDB) supplementary
    suppledb, err := boltbk.New(cfg.Supplement.DatabasePath)
    if err != nil {
        log.Fatal(trace.Wrap(err))
        return
    }

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
    repoDB.Find(&repos)

    if handleUpdate {
        var repoCount int = len(repos)
        for i, repo := range repos {
            log.Infof("%d / %d | %s - %s", i, repoCount, repo.RepoId, repo.RepoPage)
            resp, err := GithubSupplementInfo(suppledb, ctrl, &repo);
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
            githubSortSupplementInfo(suppledb, &repo);
        }
    }

    repoDB.Close()
    suppledb.Close()
    log.Info("Update process ended at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
}
