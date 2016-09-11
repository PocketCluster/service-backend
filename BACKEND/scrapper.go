package main

import (
    "log"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
)

func main() {
    db, err := gorm.Open("sqlite3", "pc-index.db")
    if err != nil {
        panic("failed to connect database " + err.Error() )
    }

    var repos []model.Repository
    db.Find(&repos)
    for _, repo := range repos {
        log.Print("commencing " + repo.RepoPage + "...")
        util.GithubReadmeScrap(repo.RepoPage, repo.Slug + ".html")
    }
}