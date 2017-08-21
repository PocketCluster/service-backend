package weekly

import (
    "net/http"
    "time"
    "fmt"

    "github.com/jinzhu/gorm"

    "github.com/stkim1/backend/util"
    "github.com/stkim1/backend/config"
    "github.com/stkim1/backend/model"
)

const (
    dateInterval int = -10
)

func findRepoInCategory(db *gorm.DB, dateInterval int, cat string) []map[string]string {
    var (
        then time.Time = time.Now().AddDate(0, 0, -5)
        repos []model.Repository
        desc []map[string]string
    )
    db.Where("? < created_at AND category = ?", then, cat).Find(&repos)
    if len(repos) == 0 {
        return nil
    }
    for _, r := range repos {
        desc = append(desc, map[string]string {
            "Title":        r.Title,
            "Link":         fmt.Sprintf("https://index.pocketcluster.io/%s.html", r.Slug),
            "Summary":      r.Summary,
        })
    }
    return desc
}

func FrontEnd(cfg *config.Config, db *gorm.DB) (string, int) {
    var (
        content map[string]interface{} = map[string]interface{} {
            "ISINDEX":        false,
            "SITENAME":       cfg.Site.SiteName,
            "DEFAULT_LANG":   "utf-8",
            "SITEURL":        cfg.Site.SiteURL,
            "THEME_LINK":     cfg.ThemeLink,
        }
    )

    content["EXAMPLE"]   = findRepoInCategory(db, dateInterval, "example")
    content["TOOLSET"]   = findRepoInCategory(db, dateInterval, "toolset")
    content["MODEL"]     = findRepoInCategory(db, dateInterval, "model")
    content["LIBRARY"]   = findRepoInCategory(db, dateInterval, "library")
    content["FRAMEWORK"] = findRepoInCategory(db, dateInterval, "framework")

    return util.RenderLayout(cfg.General.TemplatePath, "dashboard/weekly.html.mustache", "dashboard/base.html.mustache", content), http.StatusOK
}