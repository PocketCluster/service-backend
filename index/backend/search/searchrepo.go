package search

import (
    "github.com/blevesearch/bleve"

    "github.com/stkim1/service-backend/index/backend/model"
)

func NewSerachRepo(repo *model.Repository, readme *string) *SearchRepo {
    return &SearchRepo{
        RepoID:         repo.RepoId,
        Type:           SearchTypeRepoMeta,
        Slug:           repo.Slug,
        Category:       repo.Category,
        Title:          repo.Title,
        Readme:         *readme,
    }
}

type SearchRepo struct {
    RepoID     string       `json:"RepoID"`
    Type       string       `json:"Type"`
    // for 404 search & url generation
    Slug       string       `json:"Slug"`
    Category   string       `json:"Category"`
    Title      string       `json:"Title"`
    Readme     string       `json:"Readme"`
}

// Index is used to add the event in the bleve index.
func (s *SearchRepo) Index(index bleve.Index) error {
    return index.Index(s.RepoID, s)
}
