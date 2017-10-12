package search

import (
    "github.com/blevesearch/bleve"
    "github.com/blevesearch/bleve/analysis/analyzer/keyword"
    "github.com/blevesearch/bleve/analysis/lang/en"
    "github.com/blevesearch/bleve/mapping"
)

const (
    IndexNameRepoMeta   string = "indexmeta"

    SearchTypeField     string = "Type"
    SearchTypeRepoMeta  string = "repometa"

    SearchFieldSlug     string = "Slug"
    SearchFieldCategory string = "Category"
    SearchFieldTitle    string = "Title"
    SearchFieldReadme   string = "Readme"
)

func BuildIndexMapping() (mapping.IndexMapping, error) {

    // a generic reusable mapping for english text
    englishText := bleve.NewTextFieldMapping()
    englishText.Analyzer = en.AnalyzerName

    // a generic reusable mapping for keyword text
    kword := bleve.NewTextFieldMapping()
    kword.Analyzer = keyword.Name

    repoMapping := bleve.NewDocumentMapping()
    repoMapping.AddFieldMappingsAt(SearchFieldSlug, kword)
    repoMapping.AddFieldMappingsAt(SearchFieldCategory, kword)
    repoMapping.AddFieldMappingsAt(SearchFieldTitle, englishText)
    repoMapping.AddFieldMappingsAt(SearchFieldReadme, englishText)

    indexMapping := bleve.NewIndexMapping()
    indexMapping.AddDocumentMapping(SearchTypeRepoMeta, repoMapping)
    indexMapping.TypeField       = SearchTypeField
    indexMapping.DefaultAnalyzer = "en"
    return indexMapping, nil
}
