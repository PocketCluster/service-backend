package control

import (
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
	"log"
)

// Home page route
func (controller *Controller) Index(c web.C, r *http.Request) (string, int) {
	var repositories []model.Repository
	var db *gorm.DB = controller.GetGORM(c)
	db.Order("updated desc").Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
	if len(repositories) == 0 {
		return "", http.StatusNotFound
	}

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : true,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
		"CATEGORIES"		   : model.GetDefaultCategory(),
		"repositories"		   : &repositories,
		"nextpagelink" 		   : "/index.html?page=2",
	}

	return util.Render("index.html.mustache", "base.html.mustache", content), http.StatusOK
}

func (controller *Controller) IndexPaged(c web.C, r *http.Request) (string, int) {

	log.Print("paged index\n")

	var repositories []model.Repository
	var db *gorm.DB = controller.GetGORM(c)
	db.Order("updated desc").Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
	if len(repositories) == 0 {
		return "", http.StatusNotFound
	}

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : true,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
		"CATEGORIES"		   : model.GetDefaultCategory(),
		"repositories"		   : &repositories,
	}

	if SingleColumnCount * TotalRowCount <= len(repositories) {
		content["nextpagelink"] = "/index.html?page=2"
	}

	return util.Render("index.html.mustache", "base.html.mustache", content), http.StatusOK
}