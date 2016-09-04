package control

import (
    "net/http"
    "strings"
    "io/ioutil"
    "path"
    "log"

    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"
    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/model"
)

func (controller *Controller) Repository(c web.C, r *http.Request) (string, int) {
    var repositories []model.Repository
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
    db.Where("Slug = ?", repoSlug).First(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }

    var repo model.Repository  = repositories[0]
    var content map[string]interface{} = map[string]interface{} {
        "DEFAULT_LANG"         : "utf-8",
        "SITEURL"              : "https://index.pocketcluster.io",
        "THEME_STATIC_DIR"     : "theme",
        "title"                : repo.Title,
        "repo"                 : &repo,
    }

    readme, err := ioutil.ReadFile(path.Join("readme/", "readme.html"))
    if err != nil {
        log.Panic("Cnnot read readme")
    }
    content["readme"]             = string(readme)

    return util.Render("repo.html.mustache", "base.html.mustache", content), http.StatusOK
}
