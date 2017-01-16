package search

import (
    "github.com/blevesearch/bleve"
    "github.com/blevesearch/bleve/analysis/analyzer/keyword"
    "github.com/blevesearch/bleve/analysis/lang/en"
    "github.com/blevesearch/bleve/mapping"
)

func BuildIndexMapping() (mapping.IndexMapping, error) {

    // a generic reusable mapping for english text
    englishText := bleve.NewTextFieldMapping()
    englishText.Analyzer = en.AnalyzerName

    // a generic reusable mapping for keyword text
    kword := bleve.NewTextFieldMapping()
    kword.Analyzer = keyword.Name

    repoMapping := bleve.NewDocumentMapping()
    repoMapping.AddFieldMappingsAt("Category", kword)
    repoMapping.AddFieldMappingsAt("Title", englishText)
    repoMapping.AddFieldMappingsAt("Readme", englishText)

    indexMapping := bleve.NewIndexMapping()
    indexMapping.AddDocumentMapping("repo", repoMapping)
    indexMapping.TypeField       = "Category"
    indexMapping.DefaultAnalyzer = "en"
    return indexMapping, nil
}
