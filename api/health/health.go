package health

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/julienschmidt/httprouter"
)

func HealthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    hostname, _ := os.Hostname()
    now := time.Now().Local().Format("Mon, 2 Jan 2006 15:04:05 MST")
    fmt.Fprintf(w, "%s, %s\n",hostname, now)
}
