package control

import (
	"bytes"
	"html/template"

	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
)

type Controller struct {
}

func (controller *Controller) GetSession(c web.C) *sessions.Session {
	return c.Env["Session"].(*sessions.Session)
}

func (controller *Controller) GetTemplate(c web.C) *template.Template {
	return c.Env["Template"].(*template.Template)
}

func (controller *Controller) GetGORM(c web.C) *gorm.DB {
	return c.Env["GORM"].(*gorm.DB)
}

func (controller *Controller) IsXhr(c web.C) bool {
	return c.Env["IsXhr"].(bool)
}

func (controller *Controller) Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}