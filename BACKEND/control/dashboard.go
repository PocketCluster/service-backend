package control

import (
    "errors"
    "net/http"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/zenazn/goji/web"

    "github.com/stkim1/BACKEND/util"
)

// Category route
func (ctrl *Controller) DashboardFront(c web.C, r *http.Request) (string, int) {
    if ctrl.IsSafeConnection(r) {
        return "", http.StatusNotFound
    }

    var content map[string]interface{} = map[string]interface{} {
        "ISINDEX":        false,
        "SITENAME":       ctrl.Config.Site.SiteName,
        "DEFAULT_LANG":   "utf-8",
        "SITEURL":        ctrl.Config.Site.SiteURL,
        "THEME_LINK":     ctrl.Site.ThemeLink,
    }

    mode := strings.ToLower(c.URLParams["mode"]); if len(mode) == 0 || !(mode == "overview" || mode == "repository") {
        log.Error(trace.Wrap(errors.New("Cannot display page without a proper mode")))
        return "", http.StatusNotFound
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "dashboard/" + mode + ".html.mustache", "dashboard/base.html.mustache", content), http.StatusOK
}
