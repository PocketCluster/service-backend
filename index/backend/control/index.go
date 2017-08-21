package control

import (
    "net/http"
    "strconv"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"
    humanize "github.com/dustin/go-humanize"

    "github.com/stkim1/backend/util"
    "github.com/stkim1/backend/model"
)

// Home page route
func (ctrl *Controller) Index(c web.C, r *http.Request) (string, int) {
    var (
        db *gorm.DB = ctrl.GetMetaDB(c)
        repositories []model.Repository
    )
    db.Order("updated desc").Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }

    var content map[string]interface{} = map[string]interface{} {
        "SITENAME":        ctrl.Config.Site.SiteName,
        "DEFAULT_LANG":    "utf-8",
        "SITEURL":         ctrl.Config.Site.SiteURL,
        "THEME_LINK":      ctrl.Site.ThemeLink,
        "TOTAL_COUNT":     humanize.FormatInteger("##,###.", int(ctrl.TotalRepoCount.Load().(int64))),
        "CATEGORIES":      model.GetDefaultCategory(),
        "repositories":    &repositories,
        "nextpagelink":    "/index2.html",
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "index.html.mustache", content), http.StatusOK
}

func (ctrl *Controller) IndexPaged(c web.C, r *http.Request) (string, int) {
    var (
        db *gorm.DB = ctrl.GetMetaDB(c)
        repositories []model.Repository
        page int
        err error
    )
    page, err = strconv.Atoi(c.URLParams["page"])
    if err != nil {
        log.Error(errors.WithMessage(err,"Cannot convert index page query-string to proper database index"))
        return "", http.StatusNotFound
    }
    if page <= 0 {
        log.Error(trace.Wrap(errors.New("Page number cannot be smaller than 0")))
        return "", http.StatusNotFound
    }

    //FIXME : how to guard on querying for large page #?
    db.Order("updated desc").Offset(SingleColumnCount * TotalRowCount * page).Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }

    var content map[string]interface{} = map[string]interface{} {
        "SITENAME":        ctrl.Config.SiteName,
        "DEFAULT_LANG":    "utf-8",
        "SITEURL":         ctrl.Config.SiteURL,
        "THEME_LINK":      ctrl.Site.ThemeLink,
        "TOTAL_COUNT":     humanize.FormatInteger("##,###.", int(ctrl.TotalRepoCount.Load().(int64))),
        "CATEGORIES":      model.GetDefaultCategory(),
        "repositories":    &repositories,
    }

    if SingleColumnCount * TotalRowCount <= len(repositories) {
        content["nextpagelink"] = "/index" + strconv.Itoa(page + 1) + ".html"
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "index.html.mustache", content), http.StatusOK
}