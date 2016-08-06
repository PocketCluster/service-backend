package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type RepoVersion struct {
	gorm.Model
 	// repository ID
	RepoId			string
	// version string
	Version			string
	// tag/release/snapshot
	Type  			string
	// release date
	Date            time.Time
}