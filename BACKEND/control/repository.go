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
    var slug string = strings.ToLower(c.URLParams["repo"])

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

    var content map[string]interface{} = map[string]interface{} {
        "ISINDEX"			   : false,
        "SITENAME"			   : "PocketCluster Index",
        "SITEURL"              : "https://index.pocketcluster.io",
        "THEME_STATIC_DIR"     : "theme",
        "CATEGORIES"		   : model.GetActivatedCategory(repo.Category),
        "title"                : repo.Title,
        "repo"                 : &repo,
        "owner"                : &owner,
    }

    // Find Contribution relation
    db.Where("repo_id = ?", repo.RepoId).Not("author_id = ?", owner.AuthorId).Order("contribution desc").Limit(10).Find(&repoContribs)
    if len(repoContribs) != 0 {
        var contribId []string = make([]string, len(repoContribs))
        for i, r := range repoContribs {
            if owner.AuthorId != r.AuthorId {
                contribId[i] = r.AuthorId
            }
        }
        // Find Contributors
        db.Where("author_id in (?)", contribId).Find(&contribs)
        content["contribs"] = &contribs
        content["hasContribs"] = true
    }

    // Patch readme
    readme, err := ioutil.ReadFile(path.Join("readme/", slug + ".html"))
    if err != nil {
        log.Panic("Cnnot read readme")
    }
    content["readme"] = string(readme)

    return util.RenderLayout("repo.html.mustache", "base.html.mustache", content), http.StatusOK
}
