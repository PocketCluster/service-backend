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
    var owner model.Author;
    var contribs []model.Author;
    var repoContribs []model.RepoContributor
    var db *gorm.DB = controller.GetGORM(c)
    var param string = strings.ToLower(c.URLParams["repo"])

    // FIXME : check with REGEX
    // when param does not ends with .html
    if !strings.HasSuffix(param, ".html") {
        return "", http.StatusNotFound
    }

    // Split params into string array
    var slug string = strings.Split(param, ".html")[0]
    if len(slug) == 0 {
        return "", http.StatusNotFound
    }

    // Find the repo by slug
    db.Where("slug = ?", slug).First(&repositories)
    if len(repositories) == 0 {
        return "", http.StatusNotFound
    }
    var repo model.Repository  = repositories[0]

    // Find Owner
    db.Where("author_id = ?", repo.AuthorId).First(&owner)

    // Find Contribution relation
    db.Where("repo_id = ?", repo.RepoId).Order("contribution desc").Limit(10).Find(&repoContribs)
    var contribId []string = make([]string, len(repoContribs))
    for i, r := range repoContribs {
        contribId[i] = r.AuthorId
    }

    // Find Contributors
    db.Where("author_id in (?)", contribId).Find(&contribs)

    var content map[string]interface{} = map[string]interface{} {
        "DEFAULT_LANG"         : "utf-8",
        "SITEURL"              : "https://index.pocketcluster.io",
        "THEME_STATIC_DIR"     : "theme",
        "title"                : repo.Title,
        "repo"                 : &repo,
        "owner"                : &owner,
        "contribs"             : &contribs,
    }

    readme, err := ioutil.ReadFile(path.Join("readme/", slug + ".html"))
    if err != nil {
        log.Panic("Cnnot read readme")
    }
    content["readme"]             = string(readme)

    return util.Render("repo.html.mustache", "base.html.mustache", content), http.StatusOK
}
