package main

import (
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "runtime"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/thoas/stats"
    "github.com/julienschmidt/httprouter"

    _ "github.com/mattn/go-sqlite3"
    "github.com/jinzhu/gorm"

    "github.com/stkim1/api"
    "github.com/stkim1/api/auth"
    "github.com/stkim1/api/health"
    "github.com/stkim1/api/package/list"
    "github.com/stkim1/api/package/repo"
    "github.com/stkim1/api/package/meta"
    "github.com/stkim1/api/package/sync"
)

func main() {
    var (
        router = httprouter.New()
        s = stats.New()
    )
    orm, derr := gorm.Open("sqlite3", "/api-service/v014/authdata.sql")
    if derr != nil {
        log.Errorf("[DATABASE] auth db open error %v", errors.WithStack(derr).Error())
        os.Exit(2)
    }
    defer orm.Close()
    authsrvc, err := auth.NewAuthGateway(orm)
    if err != nil {
        log.Errorf("[DATABASE] initialization error %v", errors.WithStack(err).Error())
        os.Exit(2)
    }

    // setup logging
    log.SetFormatter(&log.TextFormatter{})
    logRecordPath, err := filepath.Abs("/var/log/api-service.log")
    if err != nil {
        log.Error(errors.WithStack(err).Error())
        os.Exit(1)
    }
    logRecord, err := os.OpenFile(logRecordPath, os.O_WRONLY | os.O_CREATE, 0640)
    if err != nil {
        log.Error(errors.WithStack(err).Error())
        os.Exit(1)
    }
    log.SetLevel(log.DebugLevel)
    log.SetOutput(logRecord)

    // set runetime
    runtime.GOMAXPROCS(runtime.NumCPU())

    // setup route path
    router.GET(api.URLPackageList, list.PackageList)
    router.GET(api.URLPackageRepo, repo.RepoList)
    router.GET(api.URLPackageSync, sync.PackageSync)
    router.GET(api.URLPackageMeta, meta.PackageMeta)

    // setup auth path
    router.POST(api.URLAuthCheck,  authsrvc.IsUserAuthValid)

    // misc
    router.GET(api.URLHealthCheck, health.HealthCheck)
    router.GET(api.URLAppStats, func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        s, err := json.Marshal(s.Data())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        w.Write(s)
    })

    // start serving
    log.Printf("API Service Running...")
    err = http.ListenAndServe(":8080", router);
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
}
