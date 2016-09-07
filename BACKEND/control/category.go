package control

import (
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
)

// Home page route
func (controller *Controller) Category(c web.C, r *http.Request) (string, int) {
	var repositories []model.Repository
	var category string = strings.ToLower(c.URLParams["cat"])
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
		content["nextpagelink"] = "/category/" + category + ".html?page=2"
	}

	return util.Render("index.html.mustache", "base.html.mustache", content), http.StatusOK
}
