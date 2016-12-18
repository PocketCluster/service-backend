package control

import (
    "bytes"
    "net/http"

    "github.com/zenazn/goji/web"
    "github.com/stkim1/BACKEND/model"
    // FIXME: when xmlcontent PR is merged, change this repo to the original
    "github.com/sungwoncho/go-sitemap-generator/stm"
)

func (ctrl *Controller) Sitemap(c web.C, r *http.Request) (string, int) {
    c.Env["Content-Type"] = "application/xml"

    sm := stm.NewSitemap()
    sm.Create()
    sm.SetDefaultHost(ctrl.Config.Site.SiteURL)
    sm.SetVerbose(false)

    var repos []model.Repository
    ctrl.GetGORM(c).Find(&repos)

    var buffer bytes.Buffer
    for _, repo := range repos {
        buffer.Reset()
        buffer.WriteString(repo.Slug)
        buffer.WriteString(".html")
        sm.Add(stm.URL{"loc": buffer.String(), "lastmod": repo.Updated, "changefreq":"daily", "priority":"0.5"})
    }

    return string(sm.XMLContent()), http.StatusOK
}
