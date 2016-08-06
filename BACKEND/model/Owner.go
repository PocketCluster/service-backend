package model

import (
	"github.com/jinzhu/gorm"
)

type Owner struct{
	gorm.Model
	// primary author of this repository
	AuthorId	string
	// repository ID
	RepoId	    string
}
