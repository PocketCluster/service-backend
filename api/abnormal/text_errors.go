package abnormal

import (
    "fmt"
    "net/http"
    "os"
)

// ToTextHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that ToTextHTTPError does not
// actually return err.Error(), since msg and httpStatus are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func ToTextHTTPError(err error) (msg string, httpStatus int) {
    if os.IsNotExist(err) {
        return "page not found", http.StatusNotFound
    }
    if os.IsPermission(err) {
        return "forbidden page", http.StatusForbidden
    }
    // Default:
    return "internal service error", http.StatusInternalServerError
}

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be plain text.
func ResponseTextError(w http.ResponseWriter, errMsg string, code int) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
    fmt.Fprintln(w, errMsg)
}

