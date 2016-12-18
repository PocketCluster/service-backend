package control

import (
    "net/http"
    "strings"
    "strconv"
    "errors"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"

    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/model"
)

// Category route
func (ctrl *Controller) Category(c web.C, r *http.Request) (string, int) {
    var repositories []model.Repository
    var category string = strings.ToLower(c.URLParams["cat"])
    if len(category) == 0 {
        return "", http.StatusNotFound
    }
    //FIXME : Titalize
    var title string = strings.TrimSpace(c.URLParams["cat"])
    if !model.IsCategoryPresent(category) {
        return "", http.StatusNotFound
    }

    var db *gorm.DB = ctrl.GetGORM(c)
    db.Where("category = ?", category).Order("updated desc").Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }

    var content map[string]interface{} = map[string]interface{} {
        "ISINDEX":         false,
        "SITENAME":        ctrl.Site.SiteName,
        "DEFAULT_LANG":    "utf-8",
        "SITEURL":         ctrl.Config.SiteURL,
        "THEME_LINK":      ctrl.Site.ThemeLink,
        "CATEGORIES":      model.GetActivatedCategory(category),
        "title":           title,
        "repositories":    &repositories,
    }

    if SingleColumnCount * TotalRowCount <= len(repositories) {
        content["nextpagelink"] = "/category/" + category + "2.html"
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "index.html.mustache", "base.html.mustache", content), http.StatusOK
}

// Category route
func (ctrl *Controller) CategoryPaged(c web.C, r *http.Request) (string, int) {
    page, err := strconv.Atoi(c.URLParams["page"])
    if err != nil {
        log.Error(trace.Wrap(err, "Cannot convert page string to number"))
        return "", http.StatusNotFound
    }
    if page <= 0 {
        log.Error(trace.Wrap(errors.New("Page number cannot be smaller than 0")))
        return "", http.StatusNotFound
    }

    var repositories []model.Repository
    var category string = strings.ToLower(c.URLParams["cat"])
    if len(category) == 0 {
        return "", http.StatusNotFound
    }
    //FIXME : Titalize
    var title string = strings.TrimSpace(c.URLParams["cat"])
    if !model.IsCategoryPresent(category) {
        return "", http.StatusNotFound
    }

    var db *gorm.DB = ctrl.GetGORM(c)
    //FIXME : how to guard on querying for large page #?
    db.Where("category = ?", category).Order("updated desc").Offset(SingleColumnCount * TotalRowCount * page).Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }

    var content map[string]interface{} = map[string]interface{} {
        "ISINDEX":         false,
        "SITENAME":        ctrl.Site.SiteName,
        "DEFAULT_LANG":    "utf-8",
        "SITEURL":         ctrl.Site.SiteURL,
        "THEME_LINK":      ctrl.Site.ThemeLink,
        "CATEGORIES":      model.GetActivatedCategory(category),
        "title":           title,
        "repositories":    &repositories,
    }

    if SingleColumnCount * TotalRowCount <= len(repositories) {
        content["nextpagelink"] = "/category/" + category + strconv.Itoa(page + 1) + ".html"
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "index.html.mustache", "base.html.mustache", content), http.StatusOK
}
