package main

import (
    "os"
    "path"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

import (
    "github.com/stkim1/BACKEND/config"
    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
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

    // database
    repoDB, err := gorm.Open(cfg.Database.DatabaseType, cfg.Database.DatabasePath)
    if err != nil {
        log.Fatal(trace.Wrap(err))
        return
    }
    defer repoDB.Close()

    var repos []model.Repository
    repoDB.Find(&repos)
    for _, repo := range repos {
        util.GithubReadmeScrap(repo.RepoPage, path.Join(cfg.General.ReadmePath, repo.Slug + ".html"))
    }
}