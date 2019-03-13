package update

import (
    "path"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/blevesearch/bleve"
    "github.com/google/go-github/github"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/pkg/errors"

    "github.com/stkim1/service-backend/index/backend/config"
    "github.com/stkim1/service-backend/index/backend/control"
    "github.com/stkim1/service-backend/index/backend/model"
    pocketsearch "github.com/stkim1/service-backend/index/backend/search"
    "github.com/stkim1/service-backend/index/backend/util"
)

func UpdateRepoMeta(metaDB *gorm.DB, searchIndex bleve.Index, ctrl *control.Controller, repoModel *model.Repository) (*github.Response, error) {
    var (
        cfg *config.Config                   = ctrl.Config
        lastUpdate, updatedDate time.Time
        branch, wikiPage, authorType, login, profileURL, avatarURL string
        starCount, forkCount, watchCount int64
        contribAuthor *model.Author
        resp *github.Response
        repoData *github.Repository
        ctribs []*github.Contributor
        contribInfo *model.RepoContributor
        err error
    )

    // Do not update a repo within 24 hours from the last update
    if !repoModel.UpdatedAt.IsZero() && time.Now().Sub(repoModel.UpdatedAt) < (time.Minute * time.Duration(cfg.Update.MetaUpdateCycle)) {
        log.Infof("%s :: Updated already", repoModel.RepoPage)
        return nil, nil
    }

    // URL CHECK
    if len(repoModel.RepoPage) == 0 {
        return resp, errors.New("Cannot begin update a repo with empty URL")
    }

    // save when this is updated
    lastUpdate = repoModel.Updated

    /* ------------------------------------------- Handle Repository information ------------------------------------ */
    repoData, resp, err = ctrl.GetGithubRepoMeta(repoModel.RepoPage)
    if err != nil {
        return resp, err
    }

    branch         = util.SafeGetString(repoData.DefaultBranch)
    starCount      = int64(*repoData.StargazersCount)
    forkCount      = int64(*repoData.ForksCount)
    watchCount     = int64(*repoData.SubscribersCount)
    updatedDate    = repoData.UpdatedAt.Time
    wikiPage       = ""
    if util.SafeGetBool(repoData.HasWiki) {
        wikiPage   = repoModel.RepoPage + "/wiki"
    }

    repoModel.StarCount     = starCount
    repoModel.ForkCount     = forkCount
    repoModel.WatchCount    = watchCount
    repoModel.Updated       = updatedDate
    if repoModel.WikiPage != wikiPage {
        repoModel.WikiPage = wikiPage
    }
    if repoModel.Branch != branch {
        repoModel.Branch = branch
    }
    metaDB.Save(repoModel)

    /* ------------------------------------------- Handle Contributor information ----------------------------------- */
    // contributors
    ctribs, resp, err = ctrl.GetGithubContributors(repoModel.RepoPage)
    if err != nil {
        return resp, err
    }

    for _, cauthor := range ctribs {
        // contribution
        if cauthor == nil {
            errors.Errorf(repoModel.RepoPage + " : Null contribution data. WTF?")
            continue
        }
        // user id
        cid, err := util.SafeGetInt(cauthor.ID)
        if err != nil {
            errors.Errorf(repoModel.RepoPage + err.Error())
            continue
        }
        contribID := "gh" + strconv.Itoa(cid)

        // how many times this contributor has worked
        cfactor, err := util.SafeGetInt(cauthor.Contributions)
        if err != nil {
            errors.Errorf(repoModel.RepoPage + err.Error())
            continue
        }

        // find this user
        var users []model.Author
        metaDB.Where("author_id = ?", contribID).Find(&users)
        if len(users) == 0 {
            authorType    = strings.ToLower(util.SafeGetString(cauthor.Type))
            login         = util.SafeGetString(cauthor.Login)
            profileURL    = util.SafeGetString(cauthor.HTMLURL)
            avatarURL     = util.SafeGetString(cauthor.AvatarURL)

            contribAuthor = &model.Author {
                Service:        "github",
                Type:           authorType,
                AuthorId:       contribID,
                Login:          login,
                Name:           "",
                ProfileURL:     profileURL,
                AvatarURL:      avatarURL,
                Deceased:       false,
            }
            metaDB.Save(contribAuthor)
        }

        var repoContrib []model.RepoContributor
        metaDB.Where("repo_id = ? AND author_id = ?", repoModel.RepoId, contribID).Find(&repoContrib)
        if len(repoContrib) == 0 {
            contribInfo = &model.RepoContributor{
                RepoId:         repoModel.RepoId,
                AuthorId:       contribID,
                Contribution:   cfactor,
            }
            metaDB.Save(contribInfo)
        } else {
            repoContrib[0].Contribution = cfactor
            metaDB.Save(&repoContrib[0])
        }
    }

    if 0 < updatedDate.Sub(lastUpdate) || cfg.Update.ForceReadme {
        readme, err := util.GithubReadmeScrap(repoModel.RepoPage, path.Join(ctrl.Config.General.ReadmePath, repoModel.Slug + ".html"))
        if err != nil {
            log.Error(err)
        } else {
            sr := pocketsearch.NewSerachRepo(repoModel, &readme)
            sr.Index(searchIndex)
        }
    }
    return resp, nil
}

func UpdateAllRepoMeta(metaDB *gorm.DB, searchIndex bleve.Index, cfg *config.Config, metaWaiter *sync.WaitGroup, isUpdating *atomic.Value) {
    var (
        ctrl *control.Controller        = control.NewController(cfg)
        repos []model.Repository
    )
    metaWaiter.Add(1)
    isUpdating.Store(true)

    log.Info("Meta Update process started at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
    // update start
    metaDB.Find(&repos)
    for i, _ := range repos {
        resp, err := UpdateRepoMeta(metaDB, searchIndex, ctrl, &(repos[i]))
        if err != nil {
            log.Error(errors.WithStack(err))
        }
        if resp != nil && resp.Rate.Remaining < 100 {
            log.Info("HIT API LIMIT!!!")
            break
        }
    }
    log.Info("Meta Update process ended at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
    isUpdating.Store(false)
    metaWaiter.Done()
}