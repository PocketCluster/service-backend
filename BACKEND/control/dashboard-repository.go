package control

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "strconv"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"
    "github.com/google/go-github/github"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/config"
)

func (ctrl *Controller) DashboardRepository(c web.C, r *http.Request) (string, int) {
    if !ctrl.IsSafeConnection(r) {
        return "", http.StatusNotFound
    }

    requests := map[string]string{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&requests); if err != nil {
        log.Error(trace.Wrap(err, "Cannot decode request body "))
        return "{}", http.StatusNotFound
    }

    // Check what mode this is in
    mode := strings.ToLower(c.URLParams["mode"])
    if len(mode) == 0 || !(mode == "preview" || mode == "submit") {
        log.Error(trace.Wrap(errors.New("Cannot response without a proper mode : " + mode)))
        return "", http.StatusNotFound
    }

    // GITHUB API REQUEST
    repoURL := requests["add-repo-url"]
    if len(repoURL) == 0 {
        log.Error(trace.Wrap(errors.New("Repository URL [add-repo-url] cannot be null")))
        return "{}", http.StatusNotFound
    }
    repo, _, err := ctrl.GetGithubRepoMeta(repoURL)
    if err != nil {
        log.Error(trace.Wrap(err, "Retrieving repository failed"))
        return "", http.StatusNotFound
    }

    if mode == "preview" {
        response, err := getPreview(ctrl.GetGORM(c), requests, repo)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot preview repo info"))
            return "{}", http.StatusNotFound
        }
        json, err:= json.Marshal(response);
        if err != nil {
            log.Error(trace.Wrap(err))
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    } else {
        // Decode contributor API
        contribs, _, err := ctrl.GetGithubContributorsStat(repoURL)
        if err != nil {
            log.Error(trace.Wrap(err, "Retrieving repository failed"))
            return "", http.StatusNotFound
        }
        responses, err := submitRepo(ctrl.GetGORM(c), ctrl.Config, requests, repo, contribs)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot submit the repo info"))
            return "{}", http.StatusNotFound
        }
        json, err:= json.Marshal(responses); if err != nil {
            log.Error(trace.Wrap(err, "Cannot marshal json"))
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    }
}

func submitRepo(repoDB *gorm.DB, config *config.Config, reqs map[string]string, repoData *github.Repository, ctribs []*github.ContributorStats) (map[string]interface{}, error) {

    var (
        // title
        title string            = reqs["add-repo-title"]
        // Description
        description string      = reqs["add-repo-desc"]
        // get Slug
        slug string             = reqs["add-repo-slug"]
        // Category
        category string         = strings.ToLower(reqs["add-repo-category"])
        // Project Page
        projectPage string      = reqs["add-project-page"]
        // logo image
        logoImage string        = reqs["add-logo-image"]
        // repo Page
        repoPage  string        = reqs["add-repo-url"]
    )

    /* -------------------------------------------- Submit Error Checking ------------------------------------------- */
    /*                      These are the checks that prevents errors in submit process                               */
    /* -------------------------------------------------------------------------------------------------------------- */

    // Build repo id
    rid, err := util.SafeGetInt(repoData.ID)
    if err != nil {
        return nil, trace.Wrap(err, "Cannot parse repository id")
    }
    repoID := "gh" + strconv.Itoa(rid)

    // owner info
    var owner *github.User = repoData.Owner
    if owner == nil {
        return nil, errors.New("Cannot parse Owner info of the repo")
    }

    // owner id
    aid, err := util.SafeGetInt(owner.ID)
    if err != nil {
        return nil, trace.Wrap(err, fmt.Sprintf("Cannot parse Owner[%s] id from repo.Owner.ID", util.SafeGetString(owner.Login)))
    }
    authorID := "gh" + strconv.Itoa(aid)

    // let's quickly Check database if this repo exists
    var repoFound []model.Repository
    repoDB.Where("repo_id = ? AND slug = ?", repoID, slug).Find(&repoFound);
    if len(repoFound) != 0 {
        return map[string]interface{}{
            "status"    :"duplicated",
            "reason"    :"The repository already exists",
        }, nil
    }

    /* ------------------------------------------- Handle Owner information ----------------------------------------- */
    // find and match owner
    var foundAuthor []model.Author
    repoDB.Where("author_id = ?", authorID).Find(&foundAuthor);
    if len(foundAuthor) == 0 {
        authorType    := strings.ToLower(util.SafeGetString(owner.Type))
        login         := util.SafeGetString(owner.Login)
        name          := util.SafeGetString(owner.Name)
        profileURL    := util.SafeGetString(owner.HTMLURL)
        avatarURL     := util.SafeGetString(owner.AvatarURL)

        repoAuthor := model.Author{
            Service:    "github",
            Type:       authorType,
            AuthorId:   authorID,
            Login:      login,
            Name:       name,
            ProfileURL: profileURL,
            AvatarURL:  avatarURL,
            Deceased:   false,
        }
        repoDB.Save(&repoAuthor)
    }

    /* ------------------------------------------- Handle Repository information ------------------------------------ */
    repoName        := util.SafeGetString(repoData.FullName)
    branch          := util.SafeGetString(repoData.DefaultBranch)
    forked          := util.SafeGetBool(repoData.Fork)
    starCount       := int64(*repoData.StargazersCount)
    forkCount       := int64(*repoData.ForksCount)
    watchCount      := int64(*repoData.SubscribersCount)
    createdDate     := repoData.CreatedAt.Time
    updatedDate     := repoData.UpdatedAt.Time
    wikiPage        := ""
    if *repoData.HasWiki {
        wikiPage    = repoPage + "/wiki"
    }

    repoAdded := model.Repository{
        RepoId:         repoID,
        AuthorId:       authorID,
        Deceased:       false,
        Service:        "github",
        Title:          title,
        RepoName:       repoName,
        LogoImage:      logoImage,
        Branch:         branch,
        Forked:         forked,
        StarCount:      starCount,
        ForkCount:      forkCount,
        WatchCount:     watchCount,
        ProjectPage:    projectPage,
        WikiPage:       wikiPage,
        RepoPage:       repoPage,
        Slug:           slug,
        Tags:           "",
        Category:       category,
        Summary:        description,
        Created:        createdDate,
        Updated:        updatedDate,
    }
    repoDB.Save(&repoAdded)

    // upon successful repo save, save readme to file
    util.GithubReadmeScrap(repoPage, config.General.ReadmePath + slug + ".html")

    /* ------------------------------------------- Handle Contributor information ----------------------------------- */
    for _, ctrb := range ctribs {
        // contributor
        cauthor := ctrb.Author

        // user id
        cid, err := util.SafeGetInt(cauthor.ID)
        if err != nil {
            continue
        }
        contribID := "gh" + strconv.Itoa(cid)

        // how many times this contributor has worked
        cid, err = util.SafeGetInt(ctrb.Total)
        if err != nil {
            continue
        }
        cfactor := cid

        // find this user
        var users []model.Author
        repoDB.Where("author_id = ?", contribID).Find(&users)
        if len(users) == 0 {
            authorType      := strings.ToLower(util.SafeGetString(cauthor.Type))
            login           := util.SafeGetString(cauthor.Login)
            profileUrl      := util.SafeGetString(cauthor.HTMLURL)
            avatarUrl       := util.SafeGetString(cauthor.AvatarURL)

            contribAuthor := model.Author{
                Service     :"github",
                Type        :authorType,
                AuthorId    :contribID,
                Login       :login,
                Name        :"",
                ProfileURL  :profileUrl,
                AvatarURL   :avatarUrl,
                Deceased    :false,
            }
            repoDB.Save(&contribAuthor)
        }

        var repoContrib []model.RepoContributor
        repoDB.Where("repo_id = ? AND author_id = ?", repoID, contribID).Find(&repoContrib)
        if len(repoContrib) == 0 {
            contribInfo := model.RepoContributor{
                RepoId      :repoID,
                AuthorId    :contribID,
                Contribution:int(cfactor),
            }
            repoDB.Save(&contribInfo)
        }
    }

    return map[string]interface{}{
        "status" :"ok",
    }, nil
}

// TODO check if this already exists
func getPreview(repodb *gorm.DB, requests map[string]string, repo *github.Repository) (map[string]interface{}, error) {
    var (
        slug, repoID, description string
        response map[string]interface{} = make(map[string]interface{})
    )

    // Make Slug
    slug = strings.Replace(requests["add-repo-url"], "https://github.com/", "", -1)
    slug = strings.ToLower(slug)
    slug = strings.Replace(slug, "/", "-", -1)
    slug = strings.Replace(slug, "_", "-", -1)
    slug = strings.Replace(slug, ".", "-", -1)

    // Build repo id
    rid, err := util.SafeGetInt(repo.ID)
    if err != nil {
        return nil, errors.New("Cannot parse repository id")
    }
    repoID = "gh" + strconv.Itoa(rid)

    // let's quickly Check database if this repo exists
    var repoFound []model.Repository
    repodb.Where("repo_id = ? AND slug = ?", repoID, slug).Find(&repoFound);
    if len(repoFound) != 0 {
        response["status"] = "duplicated"
        response["reason"] = "The repository already exists"
    }

    // Description
    if repo.Description == nil || len(*repo.Description) == 0 {
        description = requests["add-repo-desc"]
    } else {
        description = *repo.Description
    }

    response["add-repo-id"]    = repoID
    response["add-repo-title"] = repo.Name
    response["add-repo-slug"]  = slug
    response["add-repo-desc"]  = description
    return response, nil
}
