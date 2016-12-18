package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"

    "github.com/stkim1/BACKEND/model"
)

func main() {
    db, err := gorm.Open("sqlite3", "pc-index.db")
    if err != nil {
        panic("failed to connect database " + err.Error() )
    }
    // Migrate the schema
    db.AutoMigrate(&model.Author{}, &model.Repository{}, &model.RepoCommit{}, &model.RepoVersion{}, &model.RepoLanguage{}, &model.RepoContributor{});

    // set relation
    // db.Model(&model.Repository{}).Related(&model.RepoVersion{})
    // db.Model(&model.Repository{}).Related(&model.RepoCommit{})
    // db.Model(&model.Repository{}).Related(&model.RepoLanguage{})
}