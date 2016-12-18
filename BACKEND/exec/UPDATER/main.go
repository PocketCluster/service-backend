package main

import (
    "errors"
    "os"
    "strconv"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/control"
    "github.com/stkim1/BACKEND/config"
)

func accessGithubAPI(repoDB *gorm.DB, ctrl *control.Controller, repoModel *model.Repository) error {
    // URL CHECK
    if len(repoModel.RepoPage) == 0 {
        return trace.Wrap(errors.New("Cannot begin update a repo with empty URL"))
    }

    // save when this is updated
    lastUpdate := repoModel.Updated

    /* ------------------------------------------- Handle Repository information ------------------------------------ */
    repoData, _, err := ctrl.GetGithubRepoMeta(repoModel.RepoPage)
    if err != nil {
        return trace.Wrap(err, "Cannot access repository data " + repoModel.RepoPage)
    }

    branch         := util.SafeGetString(repoData.DefaultBranch)
    starCount      := int64(*repoData.StargazersCount)
    forkCount      := int64(*repoData.ForksCount)
    watchCount     := int64(*repoData.SubscribersCount)
    updatedDate    := repoData.UpdatedAt.Time
    wikiPage       := repoModel.RepoPage + "/wiki"

    repoModel.StarCount     = starCount
    repoModel.ForkCount     = forkCount
    repoModel.WatchCount    = watchCount
    repoModel.Updated       = updatedDate
    if util.SafeGetBool(repoData.HasWiki) && repoModel.WikiPage != wikiPage {
        repoModel.WikiPage = wikiPage
    }
    if repoModel.Branch != branch {
        repoModel.Branch = branch
    }
    repoDB.Save(repoModel)

    /* ------------------------------------------- Handle Contributor information ----------------------------------- */
    // contributors
    ctribs, _, err := ctrl.GetGithubContributorsStat(repoModel.RepoPage)
    if err != nil {
        return trace.Wrap(err, "Cannot access contributors data " + repoModel.RepoPage + err.Error())
    }

    for _, ctrb := range ctribs {
        // contribution
        if ctrb == nil {
            log.Error(trace.Wrap(errors.New("Null contribution data. WTF?")))
            continue
        }

        // contributor
        cauthor := ctrb.Author
        if cauthor == nil {
            log.Error(trace.Wrap(errors.New("Null contributor info")))
            continue
        }

        // user id
        cid, err := util.SafeGetInt(cauthor.ID)
        if err != nil {
            continue
        }
        contribID := "gh" + strconv.Itoa(cid)

        // how many times this contributor has worked
        cfactor, err := util.SafeGetInt(ctrb.Total)
        if err != nil {
            continue
        }

        // find this user
        var users []model.Author
        repoDB.Where("author_id = ?", contribID).Find(&users)
        if len(users) == 0 {
            authorType    := strings.ToLower(util.SafeGetString(cauthor.Type))
            login         := util.SafeGetString(cauthor.Login)
            profileURL    := util.SafeGetString(cauthor.HTMLURL)
            avatarURL     := util.SafeGetString(cauthor.AvatarURL)

            contribAuthor := model.Author{
                Service:        "github",
                Type:           authorType,
                AuthorId:       contribID,
                Login:          login,
                Name:           "",
                ProfileURL:     profileURL,
                AvatarURL:      avatarURL,
                Deceased:       false,
            }
            repoDB.Save(&contribAuthor)
        }

        var repoContrib []model.RepoContributor
        repoDB.Where("repo_id = ? AND author_id = ?", repoModel.RepoId, contribID).Find(&repoContrib)
        if len(repoContrib) == 0 {
            contribInfo := model.RepoContributor{
                RepoId:         repoModel.RepoId,
                AuthorId:       contribID,
                Contribution:   cfactor,
            }
            repoDB.Save(&contribInfo)
        } else {
            repoContrib[0].Contribution = cfactor
            repoDB.Save(&repoContrib[0])
        }
    }

    if 0 < updatedDate.Sub(lastUpdate) || ctrl.Config.Update.ForceReadme {
        util.GithubReadmeScrap(repoModel.RepoPage, ctrl.Config.General.ReadmePath + repoModel.Slug + ".html")
    }
    return nil
}

func main() {
    // config
    cfgPath, ok := os.LookupEnv(config.EnvConfigFilePath)
    if !ok {
        cfgPath = "config.yaml"
    }
    cfg, err := config.NewConfig(cfgPath)

    // database
    repoDB, err := gorm.Open(cfg.Database.DatabaseType, cfg.Database.DatabasePath)
    if err != nil {
        log.Error(trace.Wrap(err, "Failed to connect database"))
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
    for _, repo := range repos {
        if err := accessGithubAPI(repoDB, ctrl, &repo); err != nil {
            log.Error(err.Error())
        }
    }

    log.Info("Update process ended at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
}