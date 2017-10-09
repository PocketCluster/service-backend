package util

import (
    "encoding/json"
    "net/http"

    log "github.com/Sirupsen/logrus"
)

func JsonErrorResponse(err error) (string, int) {
    log.Error(err.Error())
    jerr, err := json.Marshal(map[string]interface{} {
        "status" : "duplicated",
        "reason" : err.Error(),
    })
    if err != nil {
        log.Error(err.Error())
        return "{}", http.StatusNotFound
    }
    return string(jerr), http.StatusNotFound
}
