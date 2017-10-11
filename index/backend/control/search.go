package control

import (
    "fmt"
    "net/http"
//    "net/url"
    "strconv"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/jinzhu/gorm"
    "github.com/zenazn/goji/web"
    "github.com/blevesearch/bleve"
    "github.com/blevesearch/bleve/search/query"
    psearch "github.com/stkim1/backend/search"
    "github.com/stkim1/backend/model"
    "github.com/stkim1/backend/util"
)

func (ctrl *Controller) ServeSearch(c web.C, req *http.Request) (string, int) {
    var (
        content map[string]interface{} = map[string]interface{}{
            "SITENAME":     ctrl.Config.SiteName,
            "DEFAULT_LANG": "utf-8",
            "SITEURL":      ctrl.Config.SiteURL,
            "THEME_LINK":   ctrl.Site.ThemeLink,
            "CATEGORIES":   model.GetDefaultCategory(),
        }
        db        *gorm.DB = ctrl.GetMetaDB(c)
        repoFound []model.Repository
        qfrom          int = 0
    )

    // find the index to operate on
    index := ctrl.GetSearchIndex(c)
    if index == nil {
        content["ERROR_MESSAGE"] = "internal search index error"
        return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
    }

    // read search queries
    qterm := req.URL.Query().Get("term")
    if len(qterm) == 0 {
        content["ERROR_MESSAGE"] = "invalid search term"
        return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
    }
    // TODO : need sanitize term
    log.Infof("Query term %v", qterm)

    // read which page this is in
    qpage := req.URL.Query().Get("page")
    if len(qpage) != 0 {
        if ipage, err := strconv.Atoi(qpage); err == nil {
            qfrom = ipage
        }
    }

    var (
        term = strings.ToLower(qterm)
        size = SingleColumnCount * TotalRowCount
        from = size * qfrom
    )

    // build a query
    srchqry := bleve.NewTermQuery(term)
    srchqry.SetField(psearch.SearchFieldReadme)
    srchqry.SetBoost(1.0)

    // build a request
    sreq := bleve.NewSearchRequestOptions(srchqry, size, from, false)

    // validate the query
    if srqv, ok := sreq.Query.(query.ValidatableQuery); ok {
        if err := srqv.Validate(); err != nil {
            content["ERROR_MESSAGE"] = "error validating query"
            return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
        }
    }

    // execute the query
    srsp, err := index.Search(sreq)
    if err != nil {
        content["ERROR_MESSAGE"] = "internal search query error"
        return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
    }

    if len(srsp.Hits) == 0 {
        content["ERROR_MESSAGE"] = "no more search result :("
    } else {
        for _, hit := range srsp.Hits {
            var repoHit model.Repository
            db.First(&repoHit, "repo_id = ?", hit.ID)
            repoFound = append(repoFound, repoHit)
        }
        content["REPOSITORIES"] = &repoFound
    }

    if size <= len(srsp.Hits) {
        nlink := fmt.Sprintf("/search?term=%v&page=%v", qterm, (from + 1))
        //nlink := url.QueryEscape(nlink)
        content["nextpagelink"] = nlink
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
}
