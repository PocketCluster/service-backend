package model

import (
    "time"
)

type RepoCommit struct {
    // repository ID
    RepoId            string
    // commit SHA
    Commit            string
    // commit Date
    Date              time.Time
}