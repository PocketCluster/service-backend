package search

import (
    "github.com/blevesearch/bleve"
    "github.com/blevesearch/bleve/analysis/analyzer/keyword"
    "github.com/blevesearch/bleve/analysis/lang/en"
    "github.com/blevesearch/bleve/mapping"
)

const (
    IndexNameRepoMeta  string       = "indexmeta"
)

const (
    SearchTypeField    string       = "Type"
    SearchTypeRepoMeta string       = "repometa"
)

func BuildIndexMapping() (mapping.IndexMapping, error) {

    // a generic reusable mapping for english text
    englishText := bleve.NewTextFieldMapping()
    englishText.Analyzer = en.AnalyzerName

    // a generic reusable mapping for keyword text
    kword := bleve.NewTextFieldMapping()
    kword.Analyzer = keyword.Name

    repoMapping := bleve.NewDocumentMapping()
    repoMapping.AddFieldMappingsAt("Slug", kword)
    repoMapping.AddFieldMappingsAt("Category", kword)
    repoMapping.AddFieldMappingsAt("Title", englishText)
    repoMapping.AddFieldMappingsAt("Readme", englishText)

    indexMapping := bleve.NewIndexMapping()
    indexMapping.AddDocumentMapping(SearchTypeRepoMeta, repoMapping)
    indexMapping.TypeField       = SearchTypeField
    indexMapping.DefaultAnalyzer = "en"
    return indexMapping, nil
}
