package control

import (
    "net/http"
    "strconv"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
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
        db        *gorm.DB = ctrl.GetMetaDB(c)
        repoFound []model.Repository
    )

    // find the index to operate on
    index := ctrl.GetSearchIndex(c)
    if index == nil {
        return util.JsonErrorResponse(errors.Errorf("search index error"))
    }

    // read search queries
    qterm := req.URL.Query().Get("term")
    if len(qterm) == 0 {
        return util.JsonErrorResponse(errors.Errorf("invalid search query"))
    }
    // TODO : need sanitize term
    log.Infof("Query term %v", qterm)

    // read which page this is in
    var qfrom int = 0
    qpage := req.URL.Query().Get("page")
    if len(qpage) != 0 {
        if ipage, err := strconv.Atoi(qpage); err == nil {
            qfrom = ipage
        }
    }

    var (
        size = SingleColumnCount * TotalRowCount
        from = size * qfrom
    )

    // build a query
    srchqry := bleve.NewTermQuery(qterm)
    srchqry.SetField(psearch.SearchFieldReadme)
    srchqry.SetBoost(1.0)

    // build a request
    sreq := bleve.NewSearchRequestOptions(srchqry, size, from, false)

    // validate the query
    if srqv, ok := sreq.Query.(query.ValidatableQuery); ok {
        if err := srqv.Validate(); err != nil {
            return util.JsonErrorResponse(errors.WithMessage(err,"error validating query"))
        }
    }

    // execute the query
    srsp, err := index.Search(sreq)
    if err != nil {
        return util.JsonErrorResponse(errors.WithMessage(err,"error executing query"))
    }

    if len(srsp.Hits) == 0 {
        return util.JsonErrorResponse(errors.Errorf("end of search"))
    }

    for _, hit := range srsp.Hits {
        var repoHit model.Repository
        db.First(&repoHit, "repo_id = ?", hit.ID)
        repoFound = append(repoFound, repoHit)
    }

    var content map[string]interface{} = map[string]interface{} {
        "SITENAME":        ctrl.Config.SiteName,
        "DEFAULT_LANG":    "utf-8",
        "SITEURL":         ctrl.Config.SiteURL,
        "THEME_LINK":      ctrl.Site.ThemeLink,
        "CATEGORIES":      model.GetDefaultCategory(),
        "repositories":    &repoFound,
    }

    if size <= len(srsp.Hits) {
        content["nextpagelink"] = "/index" + strconv.Itoa(from + 1) + ".html"
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "index.html.mustache", content), http.StatusOK
}
