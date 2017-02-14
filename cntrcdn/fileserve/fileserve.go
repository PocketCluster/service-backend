package fileserve

import (
    "fmt"
    "net/http"
    "os"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
)

// toHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that toHTTPError does not
// actually return err.Error(), since msg and httpStatus are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func toHTTPError(err error) (msg string, httpStatus int) {
    if os.IsNotExist(err) {
        return "404 page not found. Your ip address is also recorded.", http.StatusNotFound
    }
    if os.IsPermission(err) {
        return "403 Forbidden. Your ip address is also recorded.", http.StatusForbidden
    }
    // Default:
    return "500 Internal Server Error. Your ip address is also recorded.", http.StatusInternalServerError
}

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be plain text.
func Error(w http.ResponseWriter, error string, code int) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
    fmt.Fprintln(w, error)
}

func ServeImageFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string) {
    f, err := fs.Open(name)
    if err != nil {
        log.Error(trace.Wrap(err))
        msg, code := toHTTPError(err)
        Error(w, msg, code)
        return
    }
    defer f.Close()

    d, err := f.Stat()
    if err != nil {
        log.Error(trace.Wrap(err))
        msg, code := toHTTPError(err)
        Error(w, msg, code)
        return
    }

    // redirect if the directory name doesn't end in a slash
    if d.IsDir() {
        log.Error(trace.Wrap(fmt.Errorf("Requested file is a directory. This should not happen.")))
        Error(w, "404 page not found. Your ip address is also recorded.", http.StatusNotFound)
        return
    }

    w.Header().Set("Server", "PocketCluster Container CDN")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Etag", d.Name()) // file name
    w.Header().Set("Cache-Control", "max-age=86400") // 24 hrs
    w.Header().Set("Expires", time.Now().Add(time.Second * 86400).Format("Mon, 2 Jan 2006 15:04:05 MST"))

    // ServeContent will check modification time
    // we can pass empty string to name if we already set content-type
    http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
