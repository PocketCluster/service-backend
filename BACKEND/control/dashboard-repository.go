package control

import (
    "net/http"
    "log"
    "encoding/json"
    "strings"
    "strconv"
    "errors"
    "time"

    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"
    "github.com/stkim1/BACKEND/model"
    //github.com/google/go-github
)

func (Controller *Controller) DashboardRepository(c web.C, r *http.Request) (string, int) {
    requests := map[string]string{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&requests); if err != nil {
        log.Panic("Cannot decode request body " + err.Error())
        return "{}", http.StatusNotFound
    }

    // Check what mode this is in
    mode := strings.ToLower(c.URLParams["mode"]); if len(mode) == 0 || !(mode == "preview" || mode == "submit") {
        log.Panic("Cannot response without a proper mode")
        return "", http.StatusNotFound
    }

    // GITHUB API REQUEST
    if len(requests["add-repo-url"]) == 0 {
        log.Panic("Repo URL cannot be null")
        return "{}", http.StatusNotFound
    }
    apiResp, err := http.Get(GetGithubAPILink(requests["add-repo-url"])); if err != nil {
        log.Panic("Cannot Access API " + err.Error())
        return "{}", http.StatusNotFound
    }
    defer apiResp.Body.Close()

    // Decode Github API
    var githubData map[string]interface{}
    if err = json.NewDecoder(apiResp.Body).Decode(&githubData); err != nil {
        log.Panic("Cannot decode Github API body to JSON : " + err.Error())
        return "{}", http.StatusNotFound
    }

    if mode == "preview" {
        responses, err := Preview(requests, githubData); if err != nil {
            log.Panic("Cannot generate preview " + err.Error())
            return "{}", http.StatusNotFound
        }
        json, err:= json.Marshal(responses); if err != nil {
            log.Panic("Cannot marshal json " + err.Error())
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    } else {
        // Contributor Data
        contribResp, err := http.Get(GetGithubAPILink(requests["add-repo-url"] + "/contributors")); if err != nil {
            log.Panic("Cannot Access API " + err.Error())
            return "{}", http.StatusNotFound
        }
        defer contribResp.Body.Close()

        // Decode contributor API
        var contribData []map[string]interface{}
        if err = json.NewDecoder(contribResp.Body).Decode(&contribData); err != nil {
            log.Panic("Cannot decode Github API body to JSON : " + err.Error())
            return "{}", http.StatusNotFound
        }
        responses, err := Submit(Controller.GetGORM(c), requests, githubData, contribData); if err != nil {
            log.Panic("Cannot submit the repo info : " + err.Error())
            return "{}", http.StatusNotFound
        }
        json, err:= json.Marshal(responses); if err != nil {
            log.Panic("Cannot marshal json " + err.Error())
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    }
}

func Submit(db *gorm.DB, requests map[string]string, githubData map[string]interface{}, contribData []map[string]interface{}) (map[string]interface{}, error) {

    // title
    title       := requests["name"]
    // Description
    description := requests["description"]
    // get Slug
    slug        := requests["add-repo-slug"]
    // Category
    category    := strings.ToLower(requests["add-repo-category"])
    // Project Page
    projectPage := requests["add-project-page"]
    // logo image
    logoImage   := requests["add-logo-image"]
    // repo Page
    repoPage    := requests["add-repo-url"]

    // get repo id
    gid, ok := githubData["id"].(float64); if !ok {
        return nil, errors.New("Cannot parse Github ID")
    }
    repoID := "gh" + strconv.FormatInt(int64(gid), 10)

    // let's quickly Check database if this repo exists
    var repo []model.Repository
    db.Find(&repo, "repo_id = ? AND slug = ?", repoID, slug); if len(repo) != 0 {
        return map[string]interface{}{
            "status"             :"duplicated",
            "reason"             :"The repository already exists",
        }, nil
    }

    /* ------------------------------------------- Handle Owner information ----------------------------------------- */
    ownerData, ok := githubData["owner"].(map[string]interface{}); if !ok {
        return nil, errors.New("Cannot parse Owner info of the repo")
    }
    // owner id
    aid, ok := ownerData["id"].(float64); if !ok {
        return nil, errors.New("Cannot parse Owner ID")
    }
    authorID := "gh" + strconv.FormatInt(int64(aid), 10)
    // find owner
    var users []model.Author
    if db.Find(&users, "author_id = ?", authorID); len(users) == 0 {
        userType, ok    := ownerData["type"].(string); if !ok {
            return nil, errors.New("Cannot parse Owner type")
        }
        userLogin, ok   := ownerData["login"].(string); if !ok {
            return nil, errors.New("Cannot parse Owner login name")
        }
        profileUrl, ok  := ownerData["html_url"].(string); if !ok {
            return nil, errors.New("Cannot parse Owner profile page")
        }
        avatarUrl, ok   := ownerData["avatar_url"].(string); if !ok {
            return nil, errors.New("Cannot parse Owner avatar URL")
        }

        _ = model.Author{
            Service     :"github",
            Type        :userType,
            AuthorId    :authorID,
            Login       :userLogin,
            Name        :"",
            ProfileURL  :profileUrl,
            AvatarURL   :avatarUrl,
            Deceased    :false,
        }
    }

    /* ------------------------------------------- Handle Repository information ------------------------------------ */
    repoName, ok    := githubData["full_name"].(string); if !ok {
        return nil, errors.New("Cannot parse repo's full name")
    }
    branch, ok      := githubData["default_branch"].(string); if !ok {
        return nil, errors.New("Cannot parse repo's default branch")
    }
    forked, ok      := githubData["fork"].(bool); if !ok {
        return nil, errors.New("Cannot parse if the repo is forked")
    }
    starCount, ok   := githubData["stargazers_count"].(float64); if !ok {
        return nil, errors.New("Cannot parse the repo's star count")
    }
    forkCount, ok   := githubData["forks_count"].(float64); if !ok {
        return nil, errors.New("Cannot parse the repo's fork count")
    }
    watchCount, ok  := githubData["subscribers_count"].(float64); if !ok {
        return nil, errors.New("Cannot parse the repo's watch count")
    }
    created, ok      := githubData["created_at"].(string); if !ok {
        return nil, errors.New("Cannot parse when the repo's created")
    }
    createdDate, err := time.Parse(time.RFC3339, created); if err != nil {
        return nil, errors.New("Cannot parse when the repo's created " + err.Error())
    }
    updated, ok      := githubData["updated_at"].(string); if !ok {
        return nil, errors.New("Cannot parse when the repo's updated")
    }
    updatedDate, err := time.Parse(time.RFC3339, updated); if err != nil {
        return nil, errors.New("Cannot parse when the repo's created " + err.Error())
    }

    _ = model.Repository{
        RepoId          :repoID,
        AuthorId        :authorID,
        Deceased        :false,
        Service         :"github",
        Title           :title,
        RepoName        :repoName,
        LogoImage       :logoImage,
        Branch          :branch,
        Forked          :forked,
        StarCount       :int64(starCount),
        ForkCount       :int64(forkCount),
        WatchCount      :int64(watchCount),
        ProjectPage     :projectPage,
        WikiPage        :"",
        RepoPage        :repoPage,
        Slug            :slug,
        Tags            :"",
        Category        :category,
        Summary         :description,
        Created         :createdDate,
        Updated         :updatedDate,
    }

    /* ------------------------------------------- Handle Contributor information ----------------------------------- */

    for _, contrib := range contribData {
        // user id
        cid, ok     := contrib["id"].(float64); if !ok {
            return nil, errors.New("Cannot parse User ID")
        }
        contribID := "gh" + strconv.FormatInt(int64(cid), 10)
        cfactor, ok := contrib["contributions"].(float64); if !ok {
            return nil, errors.New("Cannot parse Contribution Factor")
        }

        // find this user
        var users []model.Author
        if db.Find(&users, "author_id = ?", contribID); len(users) == 0 {
            userType, ok    := contrib["type"].(string); if !ok {
                return nil, errors.New("Cannot parse Owner type")
            }
            userLogin, ok   := contrib["login"].(string); if !ok {
                return nil, errors.New("Cannot parse Owner login name")
            }
            profileUrl, ok  := contrib["html_url"].(string); if !ok {
                return nil, errors.New("Cannot parse Owner profile page")
            }
            avatarUrl, ok   := contrib["avatar_url"].(string); if !ok {
                return nil, errors.New("Cannot parse Owner avatar URL")
            }
            _ = model.Author{
                Service     :"github",
                Type        :userType,
                AuthorId    :contribID,
                Login       :userLogin,
                Name        :"",
                ProfileURL  :profileUrl,
                AvatarURL   :avatarUrl,
                Deceased    :false,
            }
        }

        var repoContrib []model.RepoContributor
        if db.Find(&repoContrib, "repo_id = ? AND author_id = ?", repoID, contribID); len(repoContrib) == 0 {
            _ = model.RepoContributor{
                RepoId      :repoID,
                AuthorId    :contribID,
                Contribution:int(cfactor),
            }
        }
    }
    return map[string]interface{}{
        "status"             :"ok",
        "reason"             :"",
    }, nil
}

func Preview(requests map[string]string, githubData map[string]interface{}) (map[string]interface{}, error) {

    // get Slug
    slug := strings.Replace(requests["add-repo-url"], "https://github.com/", "", -1)
    slug  = strings.ToLower(slug)
    slug  = strings.Replace(slug, "/", "-", -1)
    slug  = strings.Replace(slug, "_", "-", -1)
    slug  = strings.Replace(slug, ".", "-", -1)

    // get repo id
    gid, ok := githubData["id"].(float64); if !ok {
        return nil, errors.New("Cannot parse Github ID")
    }
    repoId := "gh" + strconv.FormatInt(int64(gid), 10)

    // title
    title, ok := githubData["name"].(string); if !ok {
        return nil, errors.New("Cannot parse Github Title")
    }

    // Description
    description, ok := githubData["description"].(string); if !ok {
        return nil, errors.New("Cannot parse Github Description")
    }

    if len(description) == 0 {
        description = requests["add-repo-desc"]
    }

    return map[string]interface{}{
        "add-repo-id"        :repoId,
        "add-repo-title"     :title,
        "add-repo-slug"      :slug,
        "add-repo-desc"      :description,
    }, nil
}