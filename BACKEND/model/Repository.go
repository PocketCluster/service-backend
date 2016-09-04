package model

import (
	"github.com/jinzhu/gorm"
)

type Repository struct{
	gorm.Model
	// two abbreviate chars + numbering : gh23247808
	RepoId 	        string	`gorm:"index;size:255"`
	// Is this from Github/Gitlab/Bitbucket?
	Service     	string
	// Repository Name
	Title 			string
	// Full name (owner nick + reponame)
	RepoName    	string

	// Logo Image link
	LogoImage		string
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
	// Repository Page Link (Github/GitLab/BitBuket
	RepoPage     	string

	// Slug for index pocketcluster.io
	Slug 			string
	// Dependencies : Spark, Hadoop, etc...
	Tags 			string
	// Framework/Library/Example/etc
	Category    	string
	// Short Description
	Summary     	string				`sql:"type:text"`
	// Full Readme htmlfile
	Readme      	string

	// Primary author of this repository
	Owner			Author
	// all the contributors
	Contributors 	[]RepoContributor
	// latest official release/tag/snapshot
	Versions     	[]RepoVersion
	// commit to the main repo only
	Commits 		[]RepoCommit
	// Programming Languages used
	Languages       []RepoLanguage
}
