package main

import (
    "net/http"
    "runtime"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/julienschmidt/httprouter"

    "github.com/stkim1/cdn/handle"
)

/*
import (
    "os"
    "os/signal"
    "syscall"
)

type Server struct {}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    body := "Hello World\n"
    // Try to keep the same amount of headers
    w.Header().Set("Server", "gophr")
    w.Header().Set("Connec  tion", "keep-alive")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Content-Length", fmt.Sprint(len(body)))
    fmt.Fprint(w, body)
}

func main_old() {
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)
    signal.Notify(sigchan, syscall.SIGTERM)

    server := Server{}
    log.Printf("Running...")

    go func() {
        http.Handle("/", server)
        err := http.ListenAndServe(":8080", http.FileServer(http.Dir("/cdn-content/")));
        if err != nil {
            log.Fatal(err)
        }
    }()

    <-sigchan
}
*/

func main() {
    log.Printf("Running...")

    runtime.GOMAXPROCS(runtime.NumCPU())

    router := httprouter.New()
    router.GET(handle.URLHealthCheck,     handle.HealthCheck)
    router.GET(handle.URLContainerFilter, handle.FileDownload)

    err := http.ListenAndServe(":8080", router);
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
}
