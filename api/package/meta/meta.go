package meta

import (
    "fmt"
    "net/http"
    "path"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/julienschmidt/httprouter"
    "github.com/stkim1/sharedpkg/cforigin"
    "github.com/stkim1/sharedpkg/errmsg"
    "github.com/stkim1/api"
    "github.com/stkim1/api/abnormal"
)

func serveMeta(w http.ResponseWriter, r *http.Request, fsRoot, fileName string) {
    fs := http.Dir(fsRoot)

    f, err := fs.Open(fmt.Sprintf("%s.json",fileName))
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

    // redirect if the directory name doesn't end in a slash
    if d.IsDir() {
        log.Error(errors.Errorf("invalid meta request %s", r.URL.Path))
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonResourceNotFound, http.StatusNotFound)
        return
    }

    w.Header().Set("Server", "PocketCluster API Service")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Etag", d.Name()) // file name
    w.Header().Set("Cache-Control", "max-age=3600") // 24 hrs
    w.Header().Set("Expires", time.Now().UTC().Add(time.Hour).Format("Mon, 2 Jan 2006 15:04:05 MST"))

    // ServeContent will check modification time
    // we can pass empty string to name if we already set content-type
    http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func PackageMeta(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    err := cforigin.IsOriginAllowedCountry(r)
    if err != nil {
        log.Debugf(err.Error())
        abnormal.ResponseJsonError(w, errmsg.ErrMsgJsonUnallowedCountry, http.StatusForbidden)
        return
    }

    urlComp := strings.Split(path.Clean(r.URL.Path), "/")
    mName := urlComp[len(urlComp) - 1]
    serveMeta(w, r, api.FSPackageRootMeta, mName)
}
