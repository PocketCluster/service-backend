package model

import "strings"

type Category struct {
	Name			string
	Url			    string
	IsActive		bool
}
var CategoryList []string = []string{"Example", "Toolset", "Model", "Library", "Framework"}
var CategoryUrl []string = []string{"category/example.html", "category/toolset.html", "category/model.html", "category/library.html", "category/framework.html"}

func GetDefaultCategory() []Category {
	return []Category{
		{"Example",    "category/example.html",   false},
		{"Toolset",    "category/toolset.html",   false},
		{"Model",	   "category/model.html", 	  false},
		{"Library",    "category/library.html",   false},
		{"Framework",  "category/framework.html", false},
	}
}

func GetActivatedCategory(activeCategory string) []Category {
	categories := make([]Category, len(CategoryList))
	for i, _ := range CategoryList {
		cat := Category{CategoryList[i], CategoryUrl[i], false}
		if strings.ToLower(activeCategory) == strings.ToLower(CategoryList[i]) {
			cat.IsActive = true
		}
		categories[i] = cat
	}
	return categories
}