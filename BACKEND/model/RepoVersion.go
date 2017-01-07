package model

import (
    "time"
)

type RepoVersion struct {
     // repository ID
    RepoId          string
    // version string
    Version         string
    // tag/release/snapshot
    Type            string
    // release date
    Date            time.Time
}