package update

import (
    "time"
    "sync/atomic"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/google/go-github/github"

    "github.com/stkim1/backend/model"
    "github.com/stkim1/backend/control"
    "github.com/stkim1/backend/storage"
    "github.com/stkim1/backend/config"
)

func GithubSupplementInfo(suppDB storage.Nosql, ctrl *control.Controller, repoModel *model.Repository, cfg *config.Config) (*github.Response, error) {
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
        return nil, errors.New("Cannot begin update a repo with empty URL")
    }

    suppDB.AcquireLock(repoID, time.Second)
    err = suppDB.GetObj([]string{model.RepoSuppBucket}, repoID, &repoSupp)
    suppDB.ReleaseLock(repoID)
    if err != nil {
        repoSupp = model.RepoSupplement{RepoID:repoID}
    } else {
        // give a 30 minute headroom so this repo will be updated
        if !repoSupp.Updated.IsZero() && time.Now().Sub(repoSupp.Updated) < (time.Minute * time.Duration(cfg.Update.SuppUpdateCycle)) {
            log.Infof("%s :: Updated already", repoURL)
            return nil, nil
        }
    }

    // get languages
    langs, resp, err = ctrl.GetGithubRepoLanguages(repoURL)
    if err != nil {
        return resp, err
    } else {
        repoSupp.Languages = langs
    }

    // get releases
    releases, _, resp, err = ctrl.GetGithubAllReleases(repoURL, &repoSupp.Releases, cfg.Update.MaxReleaseCollect)
    if err != nil {
        return resp, err
    } else if len(releases) != 0 {
        repoSupp.Releases = releases
    }

    // get tags
    tags, _, resp, err = ctrl.GetGithubAllTags(repoURL, &repoSupp.Tags, cfg.Update.MaxReleaseCollect + 1)
    if err != nil {
        return resp, err
    } else if len(tags) != 0 {
        repoSupp.Tags = tags
    }

    repoSupp.BuildRecentPublication(cfg.Update.MaxReleaseRebuild)
    repoSupp.Updated = time.Now()
    //log.Info(spew.Sdump(repoSupp))

    // save it to database
    log.Infof("%s - %s :: Lang [%d], Releases [%d] Tags [%d]", repoID, repoURL, len(repoSupp.Languages), len(repoSupp.Releases), len(repoSupp.Tags))
    suppDB.AcquireLock(repoID, time.Second)
    err = suppDB.UpsertObj([]string{model.RepoSuppBucket}, repoID, &repoSupp, storage.Forever)
    suppDB.ReleaseLock(repoID)
    if err != nil {
        return resp, err
    }

    return resp, nil
}

func UpdateAllRepoSupplement(metaDB *gorm.DB, suppDB storage.Nosql, cfg *config.Config, suppWaiter *sync.WaitGroup, isUpdating *atomic.Value) {
    var (
        ctrl *control.Controller        = control.NewController(cfg)
        repos []model.Repository
    )
    suppWaiter.Add(1)
    isUpdating.Store(true)
    log.Info("Supplementary Update process started at " + time.Now().Format("Jan. 2 2006 3:04 PM"))

    metaDB.Find(&repos)
    for i, _ := range repos {
        resp, err := GithubSupplementInfo(suppDB, ctrl, &(repos[i]), cfg);
        if err != nil {
            log.Error(errors.WithStack(err))
        }
        if resp != nil && resp.Rate.Remaining < 100 {
            log.Info("HIT API LIMIT!!!")
            break
        }
    }
    log.Info("Supplementary Update process ended at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
    isUpdating.Store(false)
    suppWaiter.Done()
}
