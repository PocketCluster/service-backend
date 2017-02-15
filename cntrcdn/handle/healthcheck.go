package handle

import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "os"
    "time"
)

func HealthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    hostname, _ := os.Hostname()
    now := time.Now().Format("Mon, 2 Jan 2006 15:04:05 MST")
    fmt.Fprintf(w, "%s, %s\n",hostname, now)
}
