package main

import (
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/davecgh/go-spew/spew"

    "github.com/stkim1/backend/config"
)

func main() {
    cfgPath, ok := os.LookupEnv(config.EnvConfigFilePath)
    if !ok {
        cfgPath = "config.yaml"
    }

    cfg, err := config.NewConfig(cfgPath)
    if err == nil {
        log.Print(spew.Sdump(cfg))
    }

    if false {
        cfg := &config.Config{}
        err := config.SaveConfig(cfg, "config.yaml")
        if err == nil {
            log.Printf("%v", cfg)
        }
    }
}
