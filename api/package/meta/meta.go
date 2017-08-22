package meta

import (
    "net/http"
    "path"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/julienschmidt/httprouter"
    "github.com/stkim1/api"
    "github.com/stkim1/api/abnormal"
)

func serveMeta(w http.ResponseWriter, r *http.Request, fsRoot, fileName string) {
    fs := http.Dir(fsRoot)

    f, err := fs.Open(fileName)
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
        log.Error(errors.Errorf("Requested file is a directory. This should not happen!"))
        abnormal.ResponseJsonError(w, "{\"error\":\"resource not found\"}", http.StatusNotFound)
        return
    }

    w.Header().Set("Server", "PocketCluster API Service")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Etag", d.Name()) // file name
    w.Header().Set("Cache-Control", "max-age=3600") // 24 hrs
    w.Header().Set("Expires", time.Now().UTC().Add(time.Hour).Format("Mon, 2 Jan 2006 15:04:05 MST"))

    // ServeContent will check modification time
    // we can pass empty string to name if we already set content-type
    http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func PackageMeta(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    urlComp := strings.Split(path.Clean(r.URL.Path), "/")
    mName := urlComp[len(urlComp) - 1]
    serveMeta(w, r, api.FSPackageMetaRoot, mName)
}
