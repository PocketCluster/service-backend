package handle

import (
    "net/http"
    "path"
    "strings"

//    log "github.com/Sirupsen/logrus"
//    "github.com/gravitational/trace"
    "github.com/julienschmidt/httprouter"

    "github.com/stkim1/cntrcdn/fileserve"
)

var (
    cdnRoot = http.Dir(FScdnRoot)
)

func FileDownload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    // header validation
    // jwt check

    upath := r.URL.Path
    if !strings.HasPrefix(upath, "/") {
        upath = "/" + upath
    }
    upath = path.Clean(upath)
    fname := strings.Replace(upath, PrefixContainer, "", -1)

    fileserve.ServeImageFile(w, r, cdnRoot, fname)
}

