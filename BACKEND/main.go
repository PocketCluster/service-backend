package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stkim1/BACKEND/model"
)

func main() {
	db, err := gorm.Open("sqlite3", "pc-index.db")
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.Author{}, &model.Owner{}, &model.Repository{}, &model.RepoCommit{}, &model.RepoVersion{});

}