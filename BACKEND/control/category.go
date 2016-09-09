package control

import (
	"net/http"
	"strings"
	"strconv"
	"log"

	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
)

// Category route
func (controller *Controller) Category(c web.C, r *http.Request) (string, int) {
	var repositories []model.Repository
	var category string = strings.ToLower(c.URLParams["cat"])
	if len(category) == 0 {
		return "", http.StatusNotFound
	}
	//FIXME : Titalize
	var title string = strings.TrimSpace(c.URLParams["cat"])
	if !model.IsCategoryPresent(category) {
		return "", http.StatusNotFound
	}

	var db *gorm.DB = controller.GetGORM(c)
	db.Where("category = ?", category).Order("updated desc").Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
	if len(repositories) == 0 {
		return "", http.StatusNotFound
	}

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : false,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
		"CATEGORIES"		   : model.GetActivatedCategory(category),
		"title"				   : title,
		"repositories"		   : &repositories,
	}

	if SingleColumnCount * TotalRowCount <= len(repositories) {
		content["nextpagelink"] = "/category/" + category + ".html/2"
	}

	return util.RenderLayout("index.html.mustache", "base.html.mustache", content), http.StatusOK
}

// Category route
func (controller *Controller) CategoryPaged(c web.C, r *http.Request) (string, int) {
	page, err := strconv.Atoi(c.URLParams["page"])
	if err != nil {
		log.Panic("Cannot convert page string to number : " + err.Error())
		return "", http.StatusNotFound
	}
	if page <= 0 {
		log.Panic("Page number cannot be smaller than 0.")
		return "", http.StatusNotFound
	}

	var repositories []model.Repository
	var category string = strings.ToLower(c.URLParams["cat"])
	if len(category) == 0 {
		return "", http.StatusNotFound
	}
	//FIXME : Titalize
	var title string = strings.TrimSpace(c.URLParams["cat"])
	if !model.IsCategoryPresent(category) {
		return "", http.StatusNotFound
	}

	var db *gorm.DB = controller.GetGORM(c)
	//FIXME : how to guard on querying for large page #?
	db.Where("category = ?", category).Order("updated desc").Offset(SingleColumnCount * TotalRowCount * page).Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
	if len(repositories) == 0 {
		return "", http.StatusNotFound
	}

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : false,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
		"CATEGORIES"		   : model.GetActivatedCategory(category),
		"title"				   : title,
		"repositories"		   : &repositories,
	}

	if SingleColumnCount * TotalRowCount <= len(repositories) {
		content["nextpagelink"] = "/category/" + category + ".html/" + strconv.Itoa(page + 1)
	}

	return util.RenderLayout("index.html.mustache", "base.html.mustache", content), http.StatusOK
}
