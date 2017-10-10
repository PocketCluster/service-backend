//  Copyright (c) 2014 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package search

import (
    "encoding/json"
    "net/http"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	bhttp "github.com/blevesearch/bleve/http"
    "github.com/stkim1/backend/util"
)

type ListFieldsHandler struct {
    defaultIndexName string
    IndexNameLookup  varLookupFunc
}

func NewListFieldsHandler(defaultIndexName string) *ListFieldsHandler {
    return &ListFieldsHandler{
        defaultIndexName: defaultIndexName,
    }
}

func (h *ListFieldsHandler) ServeList(c web.C, req *http.Request) (string, int) {
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

    fields, err := index.Fields()
    if err != nil {
        log.Errorf("cannot access fields list: %v", err)
        return util.JsonErrorResponse(errors.WithMessage(err, "cannot access fields list"))
    }

    data, err := json.Marshal(struct {
        Fields []string `json:"fields"`
    }{
        Fields: fields,
    })
    if err != nil {
        log.Errorf("error parsing result query: %v", err)
        return util.JsonErrorResponse(errors.WithMessage(err, "error parsing result query"))
    }
    return string(data), http.StatusOK
}
