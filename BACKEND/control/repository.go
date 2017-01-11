package control

import (
    "net/http"
    "strings"
    "io/ioutil"
    "time"
    "path"
    "fmt"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"

    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/storage"
)

func (ctrl *Controller) Repository(c web.C, r *http.Request) (string, int) {
    var (
        repositories []model.Repository
        repo model.Repository
        owner model.Author;
        contribs []model.Author;
        repoContribs []model.RepoContributor
        repoSupp model.RepoSupplement

        metaDB *gorm.DB         = ctrl.GetMetaDB(c)
        suppDB storage.Nosql    = ctrl.GetSuppleDB(c)
        slug string             = strings.ToLower(c.URLParams["repo"])

    )
    if len(slug) == 0 {
        return "", http.StatusNotFound
    }

    // Find the repo by slug
    metaDB.Where("slug = ?", slug).First(&repositories)
    if len(repositories) == 0 {
        log.Error(trace.Wrap(fmt.Errorf("Cannot find the target repository : %s",slug)))
        return "", http.StatusNotFound
    }
    repo = repositories[0]

    // Find Owner
    metaDB.Where("author_id = ?", repo.AuthorId).First(&owner)

    var content map[string]interface{} = map[string]interface{} {
        "DEFAULT_LANG":  "utf-8",
        "ISINDEX":       false,
        "SITENAME":      ctrl.Config.Site.SiteName,
        "SITEURL":       ctrl.Config.Site.SiteURL,
        "THEME_LINK":    ctrl.Site.ThemeLink,
        "CATEGORIES":    model.GetActivatedCategory(repo.Category),
        "title":         repo.Title,
        "repo":          &repo,
        "owner":         &owner,
    }

    // Find Contribution relation
    metaDB.Where("repo_id = ?", repo.RepoId).Not("author_id = ?", owner.AuthorId).Order("contribution desc").Limit(10).Find(&repoContribs)
    if len(repoContribs) != 0 {
        var contribId []string = make([]string, len(repoContribs))
        for i, r := range repoContribs {
            if owner.AuthorId != r.AuthorId {
                contribId[i] = r.AuthorId
            }
        }
        // Find Contributors
        metaDB.Where("author_id in (?)", contribId).Find(&contribs)
        content["contribs"] = &contribs
        content["hasContribs"] = true
    }

    // Find Language/ Releases/ and Tags
    suppDB.AcquireLock(repo.RepoId, time.Second)
    err := suppDB.GetObj([]string{model.RepoSuppBucket}, repo.RepoId, &repoSupp)
    suppDB.ReleaseLock(repo.RepoId)
    if err != nil {
        trace.Wrap(err)
    } else {
        list := repoSupp.RecentPublication()
        if len(list) != 0 {
            content["releases"] = list
        }
    }

    // Patch readme
    readme, err := ioutil.ReadFile(path.Join("readme/", slug + ".html"))
    if err != nil {
        log.Error(trace.Wrap(err, "Cnnot read readme html file"))
    }
    content["readme"] = string(readme)

    return util.RenderLayout(ctrl.Config.General.TemplatePath, "repo.html.mustache", "base.html.mustache", content), http.StatusOK
}
