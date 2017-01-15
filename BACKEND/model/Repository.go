package model

import (
    "time"

    "github.com/jinzhu/gorm"
)

const repoDateFormat string = "Jan. 2, 2006"

type Repository struct{
    gorm.Model
    // two abbreviate chars + numbering : gh23247808
    RepoId             string       `gorm:"index;size:255"`
    // Primary author of this repository. use Authors.author_id
    AuthorId           string       `gorm:"index;size:255"`
    // If this repo is deceased
    Deceased           bool
    // Is this from Github/Gitlab/Bitbucket?
    Service            string
    // Repository Name
    Title              string
    // Full name (owner nick + reponame)
    RepoName           string

    // Logo Image link
    LogoImage          string
    // default branch
    Branch             string
    // check if this is original
    Forked             bool

    // Star count
    StarCount          int64
    // Fork count
    ForkCount          int64
    // Watcher count
    WatchCount         int64

    // Supplmentary Page Link
    ProjectPage        string
    // Wiki page Link
    WikiPage           string
    // Repository Page Link (Github/GitLab/BitBuket
    RepoPage           string

    // Slug for index pocketcluster.io
    Slug               string
    // Dependencies : Spark, Hadoop, etc...
    Tags               string
    // Framework/Library/Example/etc
    Category           string
    // Short Description
    Summary            string       `sql:"type:text"`

    // Created Date
    Created            time.Time
    // Updated Date
    Updated            time.Time
}

func (r *Repository) CreatedDate() string {
    //return repo.Created.Format("Jan. 2 2006 3:04 PM")
    return r.Created.Format(repoDateFormat)
}

func (r *Repository) UpdatedDate() string {
    return r.Updated.Format(repoDateFormat)
}