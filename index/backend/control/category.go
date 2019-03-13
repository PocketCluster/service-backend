package control

import (
    "net/http"
    "strconv"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/dustin/go-humanize"
    "github.com/jinzhu/gorm"
    "github.com/pkg/errors"
    "github.com/zenazn/goji/web"

    "github.com/stkim1/service-backend/index/backend/model"
    "github.com/stkim1/service-backend/index/backend/util"
)

// Category route
func (ctrl *Controller) Category(c web.C, r *http.Request) (string, int) {
    var (
        repositories []model.Repository
        db *gorm.DB        = ctrl.GetMetaDB(c)
        category string    = strings.ToLower(strings.TrimSpace(c.URLParams["cat"]))
        title string       = strings.Title(strings.TrimSpace(c.URLParams["cat"]))
    )
    if len(category) == 0 {
        return "", http.StatusNotFound
    }
    if !model.IsCategoryPresent(category) {
        return "", http.StatusNotFound
    }

    db.Where("category = ?", category).Order("updated desc").Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }

    var content map[string]interface{} = map[string]interface{} {
        "SITENAME":        ctrl.Site.SiteName,
        "DEFAULT_LANG":    "utf-8",
        "SITEURL":         ctrl.Config.SiteURL,
        "TOTAL_COUNT":     humanize.FormatInteger("##,###.", int(ctrl.TotalRepoCount.Load().(int64))),
        "THEME_LINK":      ctrl.Site.ThemeLink,
        "CATEGORIES":      model.GetActivatedCategory(category),
        "title":           title,
        "repositories":    &repositories,
    }

    if SingleColumnCount * TotalRowCount <= len(repositories) {
        content["nextpagelink"] = "/category/" + category + "2.html"
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "index.html.mustache", content), http.StatusOK
}

// Category route
func (ctrl *Controller) CategoryPaged(c web.C, r *http.Request) (string, int) {
    var (
        repositories []model.Repository
        category string    = strings.ToLower(strings.TrimSpace(c.URLParams["cat"]))
        title string       = strings.Title(strings.TrimSpace(c.URLParams["cat"]))
    )
    page, err := strconv.Atoi(strings.TrimSpace(c.URLParams["page"]))
    if err != nil {
        log.Error(errors.WithMessage(err, "Cannot convert page string to number"))
        return "", http.StatusNotFound
    }
    if page <= 0 {
        log.Error(errors.Errorf("Page number cannot be smaller than 0"))
        return "", http.StatusNotFound
    }
    if len(category) == 0 {
        return "", http.StatusNotFound
    }
    if !model.IsCategoryPresent(category) {
        return "", http.StatusNotFound
    }

    var db *gorm.DB = ctrl.GetMetaDB(c)
    //FIXME : how to guard on querying for large page #?
    db.Where("category = ?", category).Order("updated desc").Offset(SingleColumnCount * TotalRowCount * page).Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }

    var content map[string]interface{} = map[string]interface{} {
        "SITENAME":        ctrl.Site.SiteName,
        "DEFAULT_LANG":    "utf-8",
        "SITEURL":         ctrl.Site.SiteURL,
        "TOTAL_COUNT":     humanize.FormatInteger("##,###.", int(ctrl.TotalRepoCount.Load().(int64))),
        "THEME_LINK":      ctrl.Site.ThemeLink,
        "CATEGORIES":      model.GetActivatedCategory(category),
        "title":           title,
        "repositories":    &repositories,
    }

    if SingleColumnCount * TotalRowCount <= len(repositories) {
        content["nextpagelink"] = "/category/" + category + strconv.Itoa(page + 1) + ".html"
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "index.html.mustache", content), http.StatusOK
}
