package control

import (
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
)

func (controller *Controller) Repository(c web.C, r *http.Request) (string, int) {
	var repo []model.Repository
	var db *gorm.DB = controller.GetGORM(c)
	var param string = strings.ToLower(c.URLParams["repo"])

	// FIXME : check with REGEX
	// when param does not ends with .html
	if !strings.HasSuffix(param, ".html") {
		return "", http.StatusNotFound
	}

	// Split params into string array
	var repoSlug string = strings.Split(param, ".html")[0]
	if len(repoSlug) == 0 {
		return "", http.StatusNotFound
	}

	// Find the repo by slug
	db.Where("Slug = ?", repoSlug).First(&repo)
	if len(repo) == 0 {
		return "", http.StatusNotFound
	}

	var content map[string]string = map[string]string{}
	content["DEFAULT_LANG"] 	= "utf-8"
	content["SITEURL"] 			= "https://index.pocketcluster.io"
	content["THEME_STATIC_DIR"] = "/theme"
	content["title"] 			= "test title"

	return util.Render("repo.html.mustache", "base.html.mustache", content), http.StatusOK
}
