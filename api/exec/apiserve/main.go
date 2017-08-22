package main

import (
    "encoding/json"
    "net/http"
    "runtime"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/thoas/stats"
    "github.com/julienschmidt/httprouter"

    "github.com/stkim1/api"
    "github.com/stkim1/api/health"
)

func main() {
    var (
        router = httprouter.New()
        s = stats.New()
    )

    // set log level
    log.SetLevel(log.DebugLevel)

    // set runetime
    runtime.GOMAXPROCS(runtime.NumCPU())

    // setup route path
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
    err := http.ListenAndServe(":8080", router);
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
}
