package control

import (
	"time"
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/stkim1/BACKEND/model"
	// FIXME: when xmlcontent PR is merged, change this repo to the original
	"github.com/sungwoncho/go-sitemap-generator/stm"
)

func (controller *Controller) Sitemap(c web.C, r *http.Request) (string, int) {
	c.Env["Content-Type"] = "application/xml"

	sm := stm.NewSitemap()
	sm.Create()
	sm.SetDefaultHost("https://index.pocketcluster.io")
	sm.SetVerbose(false)

	var repos []model.Repository
	controller.GetGORM(c).Find(&repos)
	for _, repo := range repos {
		sm.Add(stm.URL{"loc": repo.Slug, "lastmod": time.Now(), "changefreq":"daily", "priority":"0.5"})
	}

	return string(sm.XMLContent()), http.StatusOK
}
