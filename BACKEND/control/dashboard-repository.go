package control

import (
    "encoding/json"
    "errors"
    "net/http"
    "strings"
    "strconv"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"
    "github.com/google/go-github/github"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
)

func (ctl *Controller) DashboardRepository(c web.C, r *http.Request) (string, int) {
/*
    // access control based on IP
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        log.Printf("userip: %q is not IP:port", r.RemoteAddr)
        return "", http.StatusNotFound
    }

    clientIP := net.ParseIP(ip)
    if clientIP == nil {
        log.Printf("userip: %q is not IP:port", r.RemoteAddr)
        return "", http.StatusNotFound
    }
    forwarded := r.Header.Get("X-Forwarded-For")
    log.Print("Client IP " + string(clientIP) + " forwarded " + forwarded)
 */
    ipAddress := getIPAdress(r)
    if ipAddress != "198.199.115.209" {
        log.Print("Cannot display page without proper access from VPN")
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
    if len(requests["add-repo-url"]) == 0 {
        log.Error(trace.Wrap(errors.New("Repository URL [add-repo-url] cannot be null")))
        return "{}", http.StatusNotFound
    }
    repo, _, err := ctl.GetRepoMeta(requests["add-repo-url"])
    if err != nil {
        log.Info(trace.Wrap(err, "Retrieving repository failed"))
        return "", http.StatusNotFound
    }

    if mode == "preview" {
        json, err:= json.Marshal(getPreview(requests, repo));
        if err != nil {
            log.Error(trace.Wrap(err))
            return "{}", http.StatusNotFound
        }
        return string(json), http.StatusOK
    } else {
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
        responses, err := Submit(ctl.GetGORM(c), requests, githubData, contribData); if err != nil {
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
    title       := requests["add-repo-title"]
    // Description
    description := requests["add-repo-desc"]
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
    db.Where("repo_id = ? AND slug = ?", repoID, slug).Find(&repo); if len(repo) != 0 {
        return map[string]interface{}{
            "status"    :"duplicated",
            "reason"    :"The repository already exists",
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
    if db.Where("author_id = ?", authorID).Find(&users); len(users) == 0 {
        userType, ok    := ownerData["type"].(string); if !ok {
            return nil, errors.New("Cannot parse Owner type")
        } else {
            userType = strings.ToLower(userType)
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

        repoAuthor := model.Author{
            Service     :"github",
            Type        :userType,
            AuthorId    :authorID,
            Login       :userLogin,
            Name        :"",
            ProfileURL  :profileUrl,
            AvatarURL   :avatarUrl,
            Deceased    :false,
        }
        db.Save(&repoAuthor)
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

    repoAdded := model.Repository{
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
    db.Save(&repoAdded)

    // upon successful repo save, save readme to file
    util.GithubReadmeScrap(repoPage, "/www-server/readme/" + slug + ".html")

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
        if db.Where("author_id = ?", contribID).Find(&users); len(users) == 0 {
            userType, ok    := contrib["type"].(string); if !ok {
                return nil, errors.New("Cannot parse Owner type")
            } else {
                userType = strings.ToLower(userType)
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
            contribAuthor := model.Author{
                Service     :"github",
                Type        :userType,
                AuthorId    :contribID,
                Login       :userLogin,
                Name        :"",
                ProfileURL  :profileUrl,
                AvatarURL   :avatarUrl,
                Deceased    :false,
            }
            db.Save(&contribAuthor)
        }

        var repoContrib []model.RepoContributor
        if db.Where("repo_id = ? AND author_id = ?", repoID, contribID).Find(&repoContrib); len(repoContrib) == 0 {
            contribInfo := model.RepoContributor{
                RepoId      :repoID,
                AuthorId    :contribID,
                Contribution:int(cfactor),
            }
            db.Save(&contribInfo)
        }
    }

    return map[string]interface{}{
        "status"             :"ok",
    }, nil
}

// TODO check if this already exists
func getPreview(requests map[string]string, repo *github.Repository) map[string]interface{} {
    var (
        slug, repoID, description string
    )

    // Make Slug
    slug = strings.Replace(requests["add-repo-url"], "https://github.com/", "", -1)
    slug = strings.ToLower(slug)
    slug = strings.Replace(slug, "/", "-", -1)
    slug = strings.Replace(slug, "_", "-", -1)
    slug = strings.Replace(slug, ".", "-", -1)

    // Build repo id
    repoID = "gh" + strconv.FormatInt(int64(*repo.ID), 10)

    // Description
    if repo.Description == nil || len(description) == 0 {
        description = requests["add-repo-desc"]
    } else {
        description = *repo.Description
    }

    return map[string]interface{}{
        "add-repo-id"        :repoID,
        "add-repo-title"     :repo.Name,
        "add-repo-slug"      :slug,
        "add-repo-desc"      :description,
    }
}
