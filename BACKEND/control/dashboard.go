package control

import (
    "net/http"
    "log"
    "strings"

    "github.com/zenazn/goji/web"
    "github.com/stkim1/BACKEND/util"
)

// Category route
func (controller *Controller) DashboardFront(c web.C, r *http.Request) (string, int) {

    // access control based on IP
    ipAddress := getIPAdress(r)
    if ipAddress != "198.199.115.209" {
        log.Print("Cannot display page without proper access from VPN")
        return "", http.StatusNotFound
    }

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
