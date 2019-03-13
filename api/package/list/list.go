package list

import (
    "net/http"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/julienschmidt/httprouter"
    "github.com/stkim1/service-backend/shared/errmsg"
    "github.com/stkim1/service-backend/shared/cforigin"
    "github.com/stkim1/service-backend/api"
    "github.com/stkim1/service-backend/api/abnormal"
)

func serveList(w http.ResponseWriter, r *http.Request, fsRoot, fileName string) {

    f, err := http.Dir(fsRoot).Open(fileName)
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
    w.Header().Set("Expires", time.Now().UTC().Add(time.Hour).Format("Mon, 2 Jan 2006 15:04:05 MST"))

    // ServeContent will check modification time
    // we can pass empty string to name if we already set content-type
    http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func PackageList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    err := cforigin.IsOriginAllowedCountry(r)
    if err != nil {
        log.Error(err.Error())
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonUnallowedCountry, http.StatusForbidden)
        return
    }

    serveList(w, r, api.FSPackageRootList, api.FilePackageList)
}