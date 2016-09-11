package control

import (
    "net/http"
    "log"
    "net"
    "strings"

    "github.com/zenazn/goji/web"
    "github.com/stkim1/BACKEND/util"
)

// Category route
func (controller *Controller) DashboardFront(c web.C, r *http.Request) (string, int) {

    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        log.Printf("userip: %q is not IP:port", r.RemoteAddr)
        return "", http.StatusNotFound
    }

    clientIP := net.ParseIP(ip)
    if clientIP == nil {
        log.Printf("userip: %q is not IP:port", r.RemoteAddr)
        return "", http.StatusNotFound
    }
    //forward := r.Header.Get("X-Forwarded-For")

    var content map[string]interface{} = map[string]interface{} {
        "ISINDEX"              : false,
        "SITENAME"             : "PocketCluster Index",
        "DEFAULT_LANG"         : "utf-8",
        "SITEURL"              : "https://index.pocketcluster.io",
        "THEME_STATIC_DIR"     : "theme",
    }

    mode := strings.ToLower(c.URLParams["mode"]); if len(mode) == 0 || !(mode == "overview" || mode == "repository") {
        log.Panic("Cannot display page without a proper mode")
        return "", http.StatusNotFound
    }

    return util.RenderLayout("dashboard/" + mode + ".html.mustache", "dashboard/base.html.mustache", content), http.StatusOK
}
