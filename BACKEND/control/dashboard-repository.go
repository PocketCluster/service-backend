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
    "github.com/boltdb/bolt"
    "github.com/google/go-github/github"
    "gopkg.in/vmihailenco/msgpack.v2"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/config"

    "github.com/davecgh/go-spew/spew"
)

const (
    modeStrings string = "preview update submit"
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
    mode := strings.ToLower(strings.TrimSpace(c.URLParams["mode"]))
    if len(mode) == 0 || !strings.Contains(modeStrings, mode) {
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

// this is how we'll get lang/release/tag and save it to boltdb
{
    langs, _, err := ctrl.GetGithubRepoLanguages(repoURL)
    if err != nil {
        log.Error(trace.Wrap(err))
        return "", http.StatusNotFound
    }
    log.Info(spew.Sdump(langs))

    releases, _, err := ctrl.GetGithubAllReleases(repoURL)
    if err != nil {
        log.Info(err.Error())
        log.Error(trace.Wrap(err))
        return "", http.StatusNotFound
    }
    log.Info(spew.Sdump(releases))

    tags, _, err := ctrl.GetGithubAllTags(repoURL)
    if err != nil {
        log.Info(err.Error())
        log.Error(trace.Wrap(err))
        return "", http.StatusNotFound
    }
    //log.Info(spew.Sdump(tags))

    value, err := msgpack.Marshal(tags)
    if err != nil {
        log.Info(err.Error())
    } else {
        log.Info(spew.Sdump(value))
    }


    rid, err := util.SafeGetInt(repo.ID)
    if err != nil {
        log.Error("Cannot parse repository id")
    }
    repoID := "gh" + strconv.Itoa(rid)

    var tagBucket = []byte("github-repo-tags")
    var key = model.MakeTagEntryKey(repoID)

    // store some data
    err = ctrl.GetSuppleDB(c).Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists(tagBucket)
        if err != nil {
            return err
        }
        err = bucket.Put(key, value)
        if err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        log.Error(err)
    }

    // retrieve the data
    err = ctrl.GetSuppleDB(c).View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(tagBucket)
        if bucket == nil {
            return fmt.Errorf("Bucket %q not found!", tagBucket)
        }

        val := bucket.Get(key)
        log.Info(spew.Sdump(val))

        return nil
    })

    if err != nil {
        log.Error(err)
    }
}


    switch mode {
    case "preview": {
        response, err := getPreview(ctrl.GetMetaDB(c), requests, repo)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot preview repo info"))
            return "{}", http.StatusNotFound
        }
        json, err := json.Marshal(response);
        if err != nil {
            log.Error(trace.Wrap(err))
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    }
    case "update": {
        // Decode contributor API
        contribs, _, err := ctrl.GetGithubContributors(repoURL)
        if err != nil {
            log.Error(trace.Wrap(err, "Retrieving repository contribution data failed " + util.SafeGetString(repo.HTMLURL)))
            return "", http.StatusNotFound
        }
        responses, err := updateRepo(ctrl.GetMetaDB(c), ctrl.Config, requests, repo, contribs)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot update the repo info " + util.SafeGetString(repo.HTMLURL)))
            return "{}", http.StatusNotFound
        }
        json, err := json.Marshal(responses); if err != nil {
            log.Error(trace.Wrap(err, "Cannot marshal json"))
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    }
    case "submit": {
        // Decode contributor API
        contribs, _, err := ctrl.GetGithubContributors(repoURL)
        if err != nil {
            log.Error(trace.Wrap(err, "Retrieving repository contribution data failed " + util.SafeGetString(repo.HTMLURL)))
            return "", http.StatusNotFound
        }
        responses, err := submitRepo(ctrl.GetMetaDB(c), ctrl.Config, requests, repo, contribs)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot submit the repo info " + util.SafeGetString(repo.HTMLURL)))
            return "{}", http.StatusNotFound
        }
        json, err := json.Marshal(responses); if err != nil {
            log.Error(trace.Wrap(err, "Cannot marshal json"))
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    }
    }
    return "{}", http.StatusNotFound
}

func submitRepo(repoDB *gorm.DB, config *config.Config, reqs map[string]string, repoData *github.Repository, ctribs []*github.Contributor) (map[string]interface{}, error) {
    //TODO check validity of these variables
    var (
        // title
        title string            = strings.TrimSpace(reqs["add-repo-title"])
        // Description
        description string      = strings.TrimSpace(reqs["add-repo-desc"])
        // get Slug
        slug string             = strings.TrimSpace(reqs["add-repo-slug"])
        // Category
        category string         = strings.ToLower(strings.TrimSpace(reqs["add-repo-cat"]))
        // Project Page
        projectPage string      = strings.TrimSpace(reqs["add-repo-proj"])
        // logo image
        logoImage string        = strings.TrimSpace(reqs["add-repo-logo"])
        // repo Page
        repoPage  string        = strings.TrimSpace(reqs["add-repo-url"])
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
    for _, cauthor := range ctribs {
        // contribution
        if cauthor == nil {
            log.Error(trace.Wrap(errors.New("Null contribution data. WTF?")))
            continue
        }

        // user id
        cid, err := util.SafeGetInt(cauthor.ID)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot access contributor ID"))
            continue
        }
        contribID := "gh" + strconv.Itoa(cid)

        // how many times this contributor has worked
        cfactor, err := util.SafeGetInt(cauthor.Contributions)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot parse contribution count"))
            continue
        }

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
                RepoId:         repoID,
                AuthorId:       contribID,
                Contribution:   cfactor,
            }
            repoDB.Save(&contribInfo)
        }
    }

    return map[string]interface{}{
        "status" :"ok",
    }, nil
}

func updateRepo(repoDB *gorm.DB, config *config.Config, reqs map[string]string, repoData *github.Repository, ctribs []*github.Contributor) (map[string]interface{}, error) {
    //TODO check validity of these variables
    var (
        // title
        title string            = strings.TrimSpace(reqs["add-repo-title"])
        // Description
        description string      = strings.TrimSpace(reqs["add-repo-desc"])
        // get Slug
        slug string             = strings.TrimSpace(reqs["add-repo-slug"])
        // Category
        category string         = strings.ToLower(strings.TrimSpace(reqs["add-repo-cat"]))
        // Project Page
        projectPage string      = strings.TrimSpace(reqs["add-repo-proj"])
        // logo image
        logoImage string        = strings.TrimSpace(reqs["add-repo-logo"])
        // repo Page
        repoPage  string        = strings.TrimSpace(reqs["add-repo-url"])
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
    /* ------------------------------------------- Handle Repository information ------------------------------------ */
    {
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

        if len(repoFound) == 0 {
            log.Error(trace.Wrap(errors.New("Absence of repository from database in update should never happen : " + util.SafeGetString(repoData.HTMLURL))))
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
        } else {
            repoModel := repoFound[0]
            repoModel.StarCount       = starCount
            repoModel.ForkCount       = forkCount
            repoModel.WatchCount      = watchCount
            if repoModel.Title != title {
                repoModel.Title       = title
            }
            if repoModel.RepoName != repoName {
                repoModel.RepoName    = repoName
            }
            if repoModel.LogoImage != logoImage {
                repoModel.LogoImage   = logoImage
            }
            if repoModel.Branch != branch {
                repoModel.Branch      = branch
            }
            if repoModel.ProjectPage != projectPage {
                repoModel.ProjectPage = projectPage
            }
            if repoModel.WikiPage != wikiPage {
                repoModel.WikiPage    = wikiPage
            }
            if repoModel.RepoPage != repoPage {
                repoModel.RepoPage    = repoPage
            }
            if repoModel.Category != category {
                repoModel.Category    = category
            }
            if repoModel.Summary != description {
                repoModel.Summary     = description
            }
            repoDB.Save(&repoModel)

            if 0 < updatedDate.Sub(repoModel.Updated) {
                repoModel.Updated     = updatedDate
                util.GithubReadmeScrap(repoModel.RepoPage, config.General.ReadmePath + repoModel.Slug + ".html")
            }
        }
    }

    /* ------------------------------------------- Handle Contributor information ----------------------------------- */
    for _, cauthor := range ctribs {
        // contribution
        if cauthor == nil {
            log.Error(trace.Wrap(errors.New("Null contribution data. WTF?")))
            continue
        }

        // user id
        cid, err := util.SafeGetInt(cauthor.ID)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot access contributor ID"))
            continue
        }
        contribID := "gh" + strconv.Itoa(cid)

        // how many times this contributor has worked
        cfactor, err := util.SafeGetInt(cauthor.Contributions)
        if err != nil {
            log.Error(trace.Wrap(err, "Cannot parse contribution count"))
            continue
        }

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
                RepoId:         repoID,
                AuthorId:       contribID,
                Contribution:   cfactor,
            }
            repoDB.Save(&contribInfo)
        } else {
            repoContrib[0].Contribution = cfactor
            repoDB.Save(&repoContrib[0])
        }
    }

    return map[string]interface{}{
        "status" :"ok",
    }, nil
}

func getPreview(repodb *gorm.DB, requests map[string]string, repoData *github.Repository) (map[string]interface{}, error) {
    var (
        slug, repoID, description string
        //response map[string]interface{} = make(map[string]interface{})
    )

    // Make Slug
    slug = strings.Replace(strings.TrimSpace(requests["add-repo-url"]), "https://github.com/", "", -1)
    slug = strings.ToLower(slug)
    slug = strings.Replace(slug, "/", "-", -1)
    slug = strings.Replace(slug, "_", "-", -1)
    slug = strings.Replace(slug, ".", "-", -1)

    // Build repo id
    rid, err := util.SafeGetInt(repoData.ID)
    if err != nil {
        return nil, errors.New("Cannot parse repository id")
    }
    repoID = "gh" + strconv.Itoa(rid)

    // let's quickly Check database if this repo exists
    var repoFound []model.Repository
    repodb.Where("repo_id = ? AND slug = ?", repoID, slug).Find(&repoFound);
    if len(repoFound) != 0 {
        repoModel := repoFound[0]
        return map[string]interface{}{
            "status":           "duplicated",
            "reason":           "The repository already exists",
            "add-repo-id":      repoModel.RepoId,
            "add-repo-title":   repoModel.Title,
            "add-repo-slug":    repoModel.Slug,
            "add-repo-desc":    repoModel.Summary,
            "add-repo-cat":     repoModel.Category,
            "add-repo-proj":    repoModel.ProjectPage,
            "add-repo-logo":    repoModel.LogoImage,
        }, nil
    }

    // Description
    if repoData.Description == nil || len(*repoData.Description) == 0 {
        description = strings.TrimSpace(requests["add-repo-desc"])
    } else {
        description = *repoData.Description
    }

    return map[string]interface{}{
        "add-repo-id":      repoID,
        "add-repo-title":   repoData.Name,
        "add-repo-slug":    slug,
        "add-repo-desc":    description,
    }, nil
}
