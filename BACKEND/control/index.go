package control

import (
	"net/http"
	"strconv"
	"log"

	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
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
		"nextpagelink" 		   : "/index.html/2",
	}

	return util.Render("index.html.mustache", "base.html.mustache", content), http.StatusOK
}

func (controller *Controller) IndexPaged(c web.C, r *http.Request) (string, int) {
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
	var db *gorm.DB = controller.GetGORM(c)
	//FIXME : how to guard on querying for large page #?
	db.Order("updated desc").Offset(SingleColumnCount * TotalRowCount * page).Limit(SingleColumnCount * TotalRowCount).Find(&repositories)
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
		content["nextpagelink"] = "/index.html/" + strconv.Itoa(page + 1)
	}

	return util.Render("index.html.mustache", "base.html.mustache", content), http.StatusOK
}