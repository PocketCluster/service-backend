package model

import (
    "github.com/jinzhu/gorm"
    "time"
)

type RepoCommit struct {
    gorm.Model
    // repository ID
    RepoId            string
    // commit SHA
    Commit            string
    // commit Date
    Date              time.Time
}