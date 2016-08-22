package control

import (
	"time"
	"bytes"
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

	var buffer bytes.Buffer
	var repos []model.Repository
	controller.GetGORM(c).Find(&repos)

	for _, repo := range repos {
		buffer.Reset()
		buffer.WriteString(repo.Slug)
		buffer.WriteString(".html")
		// FIXME : fix time to final updated timestamp
		sm.Add(stm.URL{"loc": buffer.String(), "lastmod": time.Now(), "changefreq":"daily", "priority":"0.5"})
	}

	return string(sm.XMLContent()), http.StatusOK
}
