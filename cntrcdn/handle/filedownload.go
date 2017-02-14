package handle

import (
//    "fmt"
    "net/http"
    "path"
    "strings"

//    log "github.com/Sirupsen/logrus"
//    "github.com/gravitational/trace"
    "github.com/julienschmidt/httprouter"

    "github.com/stkim1/cntrcdn/fileserve"
)

const (
    ImageUrlPrefix = "/image/v014/"
    cdnFsRoot = "/cdn-content/"
)

var (
    cdnRoot = http.Dir(cdnFsRoot)
)

func FileDownload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    // header validation
    // jwt check

    upath := r.URL.Path
    if !strings.HasPrefix(upath, "/") {
        upath = "/" + upath
    }
    upath = path.Clean(upath)
    fname := strings.Replace(upath, ImageUrlPrefix, "", -1)

    fileserve.ServeImageFile(w, r, cdnRoot, fname)
}

