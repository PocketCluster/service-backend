package control

import (
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
)

type Controller struct {
}

func (controller *Controller) GetSession(c web.C) *sessions.Session {
	return c.Env["Session"].(*sessions.Session)
}

func (controller *Controller) GetGORM(c web.C) *gorm.DB {
	return c.Env["GORM"].(*gorm.DB)
}

func (controller *Controller) IsXhr(c web.C) bool {
	return c.Env["IsXhr"].(bool)
}