package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    log "github.com/Sirupsen/logrus"
    //"github.com/gravitational/trace"
)

var (
    abort bool
)

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    body := "Hello World\n"
    // Try to keep the same amount of headers
    w.Header().Set("Server", "gophr")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Content-Length", fmt.Sprint(len(body)))
    fmt.Fprint(w, body)
}

func main() {
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)
    signal.Notify(sigchan, syscall.SIGTERM)

    server := Server{}

    go func() {
        http.Handle("/", server)
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatal(err)
        }
    }()

    log.Print("Hello World")

    <-sigchan
}
