package main

import (
    "net/http"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/julienschmidt/httprouter"

    "github.com/stkim1/cntrcdn/handle"
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
    w.Header().Set("Connection", "keep-alive")
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

    router := httprouter.New()
    router.GET("/healthcheck",    handle.HealthCheck)
    router.GET(handle.ImageUrlPrefix + ":filename", handle.FileDownload)

    err := http.ListenAndServe(":8080", router);
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }
}
