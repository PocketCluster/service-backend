package main

import (
    "os"
    "path"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/blevesearch/bleve"
)

import (
    "github.com/stkim1/BACKEND/config"
    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
    pocketsearch "github.com/stkim1/BACKEND/search"
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
        log.Fatal(trace.Wrap(err))
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