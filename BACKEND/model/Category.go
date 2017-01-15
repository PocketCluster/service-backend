package model

import "strings"

type Category struct {
    Name            string
    Url             string
    IsActive        bool
}
var CategoryList []string = []string{"Example", "Toolset", "Model", "Library", "Framework"}
var CategoryUrl  []string = []string{"category/example.html", "category/toolset.html", "category/model.html", "category/library.html", "category/framework.html"}

func GetDefaultCategory() []Category {
    return []Category{
        {"Example",    "category/example.html",   false},
        {"Toolset",    "category/toolset.html",   false},
        {"Model",      "category/model.html",     false},
        {"Library",    "category/library.html",   false},
        {"Framework",  "category/framework.html", false},
    }
}

func GetActivatedCategory(activeCategory string) []Category {
    var targetCategory string = strings.ToLower(activeCategory)
    var categories []Category = make([]Category, len(CategoryList))
    for i, categoryName := range CategoryList {
        var cat Category = Category{CategoryList[i], CategoryUrl[i], false}
        if targetCategory == strings.ToLower(categoryName) {
            cat.IsActive = true
        }
        categories[i] = cat
    }
    return categories
}

func IsCategoryPresent(category string) bool {
    for _, entity := range CategoryList {
        if category == strings.ToLower(entity) {
            return true
        }
    }
    return false
}