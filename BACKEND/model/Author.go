package model

import (
	"github.com/jinzhu/gorm"
)

type Author struct{
	gorm.Model
	// two abbreviate chars + numbering : gh23247808
	AuthorId 		string
	// Is this from Github/Gitlab/Bitbucket
	Service     	string
	// Profile name
	Name			string
	// Profile Type : Organization/Personal/etc
	Type            string

	// home Page Link
	HomeURL			string
	// service profile Page
	ProfileURL		string
	// Avatar Link
	AvatarURL		string
}
