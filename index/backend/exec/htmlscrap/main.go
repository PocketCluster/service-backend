package main

import (
    "os"
    "path"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/blevesearch/bleve"

    "github.com/stkim1/backend/config"
    "github.com/stkim1/backend/model"
    "github.com/stkim1/backend/util"
    pocketsearch "github.com/stkim1/backend/search"
)

func main()  {
    // config
    cfgPath, ok := os.LookupEnv(config.EnvConfigFilePath)
    if !ok {
        cfgPath = "config.yaml"
    }
    cfg, err := config.NewConfig(cfgPath)
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    rsIndx, err := bleve.Open(cfg.Search.IndexStoragePath)
    if err != nil {
        m, err := pocketsearch.BuildIndexMapping()
        if err != nil {
            log.Fatal(err)
        }
        rsIndx, err = bleve.New(cfg.Search.IndexStoragePath, m)
        if err != nil {
            log.Fatal(err)
        }
    }
    defer rsIndx.Close()

    // database
    repoDB, err := gorm.Open(cfg.Database.DatabaseType, cfg.Database.DatabasePath)
    if err != nil {
        log.Fatal(errors.WithStack(err))
        return
    }
    defer repoDB.Close()

    var repos []model.Repository
    repoDB.Find(&repos)
    for i, repo := range repos {
        readme, err := util.GithubReadmeScrap(repo.RepoPage, path.Join(cfg.General.ReadmePath, repo.Slug + ".html"))
        if err != nil {
            log.Error(err)
        } else {
            sr := pocketsearch.NewSerachRepo(&(repos[i]), &readme)
            sr.Index(rsIndx)
        }
    }
}