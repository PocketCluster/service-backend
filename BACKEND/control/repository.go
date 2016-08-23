package control

import (
	"net/http"
	"html/template"

	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
	"github.com/zenazn/goji/web"
)


func (controller *Controller) Repository(c web.C, r *http.Request) (string, int) {
	t := controller.GetTemplate(c)
	widgets := util.Parse(t, "home", nil)

	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsIndex"] = true

	c.Env["Title"] = "Default Project - free Go website project template"
	c.Env["Content"] = template.HTML(widgets)


	var repo *model.Repository
	db := controller.GetGORM(c)
	db.Where("stub = ?", "dekhtiarjonathan-neural-nets-are-weird").First(repo)
	if repo == nil {
		return "", http.StatusNotFound
	}

	return util.Parse(t, "main", c.Env), http.StatusOK
}
