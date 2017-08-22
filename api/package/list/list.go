package list

import (
    "net/http"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/julienschmidt/httprouter"
    "github.com/stkim1/api"
    "github.com/stkim1/api/abnormal"
)

func PackageList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    f, err := http.Dir(api.FSPackageRoot).Open(api.FilePackageList)
    if err != nil {
        log.Error(errors.WithStack(err))
        msg, code := abnormal.ToJsonHTTPError(err)
        abnormal.ResponseJsonError(w, msg, code)
        return
    }
    defer f.Close()

    d, err := f.Stat()
    if err != nil {
        log.Error(errors.WithStack(err))
        msg, code := abnormal.ToJsonHTTPError(err)
        abnormal.ResponseJsonError(w, msg, code)
        return
    }

    w.Header().Set("Server", "PocketCluster API Service")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Etag", d.Name()) // file name
    w.Header().Set("Cache-Control", "max-age=3600") // 1 hr
    w.Header().Set("Expires", time.Now().UTC().Add(time.Second * 3600).Format("Mon, 2 Jan 2006 15:04:05 MST"))

    // ServeContent will check modification time
    // we can pass empty string to name if we already set content-type
    http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
