package overview

import (
    "net/http"

    "github.com/zenazn/goji/web"
    "github.com/stkim1/backend/util"
    "github.com/stkim1/backend/config"
)

func FrontEnd(cfg *config.Config, c web.C, r *http.Request) (string, int) {
    var content map[string]interface{} = map[string]interface{} {
        "ISINDEX":        false,
        "SITENAME":       cfg.Site.SiteName,
        "DEFAULT_LANG":   "utf-8",
        "SITEURL":        cfg.Site.SiteURL,
        "THEME_LINK":     cfg.ThemeLink,
    }
    return util.RenderLayout(cfg.General.TemplatePath, "dashboard/overview.html.mustache", "dashboard/base.html.mustache", content), http.StatusOK
}