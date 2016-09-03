package control

import (
	"net/http"
	"html/template"
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

	t := controller.GetTemplate(c)
	widgets := util.Parse(t, "home", nil)
	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsIndex"] = true
	c.Env["Title"] = "Default Project - free Go website project template"
	c.Env["Content"] = template.HTML(widgets)
	return util.Parse(t, "main", c.Env), http.StatusOK
}
