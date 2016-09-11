package model

import (
    "github.com/jinzhu/gorm"
)

type RepoContributor struct{
    gorm.Model
    // repository ID
    RepoId            string
    // primary author of this repository
    AuthorId          string
    // contribution count
    Contribution      int
}
