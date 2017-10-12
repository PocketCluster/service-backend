package control

import (
    "encoding/base64"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "github.com/jinzhu/gorm"
    "github.com/zenazn/goji/web"
    humanize "github.com/dustin/go-humanize"
    "github.com/blevesearch/bleve"
    "github.com/blevesearch/bleve/search/query"
    psearch "github.com/stkim1/backend/search"
    "github.com/stkim1/backend/model"
    "github.com/stkim1/backend/util"
)

const (
    URLQueryPartTerm string = "term"
    URLQueryPartPage string = "page"
    URLQueryPartEtrm string = "etrm" // base64 encoded term
)

func (ctrl *Controller) ServeSearch(c web.C, req *http.Request) (string, int) {
    var (
        content map[string]interface{} = map[string]interface{}{
            "SITENAME":        ctrl.Config.SiteName,
            "DEFAULT_LANG":    "utf-8",
            "SITEURL":         ctrl.Config.SiteURL,
            "THEME_LINK":      ctrl.Site.ThemeLink,
            "TOTAL_COUNT":     humanize.FormatInteger("##,###.", int(ctrl.TotalRepoCount.Load().(int64))),
            "CATEGORIES":      model.GetDefaultCategory(),
        }
        db        *gorm.DB = ctrl.GetMetaDB(c)
        repoFound []model.Repository
        qfrom     int = 0
        qterm, encodedTerm string
    )

    // find the index to operate on
    index := ctrl.GetSearchIndex(c)
    if index == nil {
        content["ERROR_MESSAGE"] = "internal search index error"
        return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
    }

    // read which page this is in
    qpage := req.URL.Query().Get(URLQueryPartPage)
    if len(qpage) != 0 {
        if ipage, err := strconv.Atoi(qpage); err == nil {
            qfrom = ipage
        }
    }

    // read search queries when page # is 0
    if qfrom == 0 {
        qterm = req.URL.Query().Get(URLQueryPartTerm)
        if len(qterm) == 0 {
            content["ERROR_MESSAGE"] = "invalid search term"
            return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
        }

    // read encoded query for other pages
    } else {
        encodedTerm = req.URL.Query().Get(URLQueryPartEtrm)
        estr, err := base64.RawStdEncoding.DecodeString(encodedTerm)
        if len(encodedTerm) == 0 || err != nil {
            content["ERROR_MESSAGE"] = "invalid search term"
            return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
        }
        qterm = string(estr)
    }

    // TODO : need sanitize term

    var (
        term = strings.ToLower(qterm)
        size = SingleColumnCount * TotalRowCount
        from = size * qfrom
    )

    // build a query
/*
    // (2017/10/10) we can use exact field matching query + compound. That's for later
    srchqry := bleve.NewTermQuery(term)
    srchqry.SetField(psearch.SearchFieldReadme)
    srchqry.SetBoost(1.0)
*/
    srchqry := bleve.NewQueryStringQuery(term)
    srchqry.SetBoost(1.0)

    // build a request
    sreq := bleve.NewSearchRequestOptions(srchqry, size, from, false)
    sreq.Fields = []string{psearch.SearchFieldTitle, psearch.SearchFieldReadme}

    // validate the query
    if srqv, ok := sreq.Query.(query.ValidatableQuery); ok {
        if err := srqv.Validate(); err != nil {
            content["ERROR_MESSAGE"] = "invalid search query"
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
        content["ERROR_MESSAGE"] = "no more result found."
    } else {
        for _, hit := range srsp.Hits {
            var repoHit model.Repository
            db.First(&repoHit, "repo_id = ?", hit.ID)
            repoFound = append(repoFound, repoHit)
        }
        content["REPOSITORIES"] = &repoFound
    }

    if size <= len(srsp.Hits) {
        if qfrom == 0 && len(encodedTerm) == 0 {
            encodedTerm = base64.RawStdEncoding.EncodeToString([]byte(qterm))
        }
        nlink := fmt.Sprintf("/search?%v=%v", URLQueryPartEtrm, encodedTerm)
        content["nextpagelink"] = nlink
    }

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "navhead.html.mustache", "search.html.mustache", content), http.StatusOK
}
