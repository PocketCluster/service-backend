package control

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/pkg/errors"
    "github.com/zenazn/goji/web"
    "github.com/blevesearch/bleve"
    "github.com/blevesearch/bleve/search/query"
    psearch "github.com/stkim1/backend/search"
    "github.com/stkim1/backend/util"
)

func (ctrl *Controller) ServeSearch(c web.C, req *http.Request) (string, int) {
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

    data, err := json.Marshal(srsp)
    if err != nil {
        return util.JsonErrorResponse(errors.WithMessage(err,"error parsing result query"))
    }
    return string(data), http.StatusOK
}
