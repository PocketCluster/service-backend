package model

import (
	"github.com/jinzhu/gorm"
)

type Repository struct{
	gorm.Model
	// two abbreviate chars + numbering : gh23247808
	RepoId 	        string
	// Is this from Github/Gitlab/Bitbucket
	Service     	string
	// Repository Name
	Title 			string
	// Full name (owner nick + reponame)
	RepoName    	string

	// Logo Image link
	LogoImage		string
	// Programming Language
	Languages   	string
	// default branch
	Branch 			string
	// check if this is original
	Forked			bool

	// Star count
	StarCount		int64
	// Fork count
	ForkCount   	int64
	// Watcher count
	WatchCount  	int64

	// Supplmentary Page Link
	ProjectPage 	string
	// Wiki page Link
	WikiPage    	string
	// Web Page Link
	WebPage     	string

	// Slug for index pocketcluster.io
	Slug 			string
	// Dependencies : Spark, Hadoop, etc...
	Tags 			string
	// Framework/Library/Example/etc
	Category    	string
	// Short Description
	Summary     	string
	// Full Readme htmlfile
	Readme      	string

	// Primary author of this repository
	Owned       	Owner
	// all the contributors
	Contributors 	[]Author
	// latest official release/tag/snapshot
	Version     	[]RepoVersion
	// commit to the main repo only
	Commit 			[]RepoCommit
}
