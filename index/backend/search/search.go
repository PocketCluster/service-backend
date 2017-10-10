package search

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
    "github.com/zenazn/goji/web"
	"github.com/blevesearch/bleve"
	bhttp "github.com/blevesearch/bleve/http"
    "github.com/blevesearch/bleve/search/query"
    "github.com/stkim1/backend/util"
)

type varLookupFunc func(req *http.Request) string

// SearchHandler can handle search requests sent over HTTP
type SearchHandler struct {
    defaultIndexName string
    IndexNameLookup  varLookupFunc
}

func NewSearchHandler(defaultIndexName string) *SearchHandler {
    return &SearchHandler{
        defaultIndexName: defaultIndexName,
    }
}

func (h *SearchHandler) ServeSearch(c web.C, req *http.Request) (string, int) {
    // find the index to operate on
    var indexName string
    if h.IndexNameLookup != nil {
        indexName = h.IndexNameLookup(req)
    }
    if indexName == "" {
        indexName = h.defaultIndexName
    }
    index := bhttp.IndexByName(indexName)
    if index == nil {
        log.Errorf("no such index '%s'", indexName)
        return util.JsonErrorResponse(errors.Errorf("no such index '%s'", indexName))
    }

    // read the request body
    requestBody, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Errorf("error reading request body: %v", err)
        return util.JsonErrorResponse(errors.WithMessage(err,"error reading request body"))
    }

    log.Infof("request body: %s", requestBody)

    // parse the request
    var searchRequest bleve.SearchRequest
    err = json.Unmarshal(requestBody, &searchRequest)
    if err != nil {
        log.Errorf("error parsing query: %v", err)
        return util.JsonErrorResponse(errors.WithMessage(err,"error parsing query"))
    }

    log.Infof("parsed request %#v", searchRequest)

    // validate the query
    if srqv, ok := searchRequest.Query.(query.ValidatableQuery); ok {
        err = srqv.Validate()
        if err != nil {
            log.Errorf("error validating query: %v", err)
            return util.JsonErrorResponse(errors.WithMessage(err,"error validating query"))
        }
    }

    // execute the query
    searchResponse, err := index.Search(&searchRequest)
    if err != nil {
        log.Errorf("error executing query: %v", err)
        return util.JsonErrorResponse(errors.WithMessage(err,"error executing query"))
    }

    data, err := json.Marshal(searchResponse)
    if err != nil {
        log.Errorf("error parsing result query: %v", err)
        return util.JsonErrorResponse(errors.WithMessage(err,"error parsing result query"))
    }
    return string(data), http.StatusOK
}
