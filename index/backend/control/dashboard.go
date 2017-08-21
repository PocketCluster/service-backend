package control

import (
    "errors"
    "net/http"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/zenazn/goji/web"

    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/control/overview"
    "github.com/stkim1/BACKEND/control/weekly"
    "github.com/stkim1/BACKEND/control/scrapping"
)

// Category route
func (ctrl *Controller) DashboardFront(c web.C, r *http.Request) (string, int) {
    if !ctrl.IsSafeConnection(r) {
        return "", http.StatusNotFound
    }
    mode := strings.ToLower(c.URLParams["mode"]);
    switch mode {
    case "overview":
        return overview.FrontEnd(ctrl.Config, c, r)
    case "weekly":
        return weekly.FrontEnd(ctrl.Config, ctrl.GetMetaDB(c))
    case "repository":
        var content map[string]interface{} = map[string]interface{} {
            "ISINDEX":        false,
            "SITENAME":       ctrl.Config.Site.SiteName,
            "DEFAULT_LANG":   "utf-8",
            "SITEURL":        ctrl.Config.Site.SiteURL,
            "THEME_LINK":     ctrl.Site.ThemeLink,
        }
        return util.RenderLayout(ctrl.Config.General.TemplatePath, "dashboard/repository.html.mustache", "dashboard/base.html.mustache", content), http.StatusOK
    case "scrapping":
        return scrapping.FrontEnd(ctrl.Config, c, r)
    }

    log.Error(trace.Wrap(errors.New("Cannot display page without a proper mode")))
    return "", http.StatusNotFound
}