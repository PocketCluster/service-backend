package main

import (
    "os"
    "path"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/config"
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
    }

    util.GithubReadmeScrap("https://github.com/spark-jobserver/spark-jobserver", path.Join(cfg.General.ReadmePath + "spark-jobserver.html"))
}