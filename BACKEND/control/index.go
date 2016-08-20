package control

import (
	"net/http"

	"html/template"

	"github.com/stkim1/BACKEND/util"
	"github.com/zenazn/goji/web"
)


// Home page route
func (controller *Controller) Index(c web.C, r *http.Request) (string, int) {
	t := controller.GetTemplate(c)

	widgets := util.Parse(t, "home", nil)

	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsIndex"] = true

	c.Env["Title"] = "Default Project - free Go website project template"
	c.Env["Content"] = template.HTML(widgets)

	return util.Parse(t, "main", c.Env), http.StatusOK
}
