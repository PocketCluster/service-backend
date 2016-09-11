package main

import (
    "log"
    "encoding/json"
    "net/http"
    "time"
    "errors"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
    "github.com/stkim1/BACKEND/control"
)


func AccessGithubAPI(db *gorm.DB, repo *model.Repository) error {

    // URL CHECK
    if len(repo.RepoPage) == 0 {
        return errors.New("Cannot begin update a repo with empty URL")
    }

    // GITHUB API REQUEST
    apiResp, err := http.Get(control.GetGithubAPILink(repo.RepoPage)); if err != nil {
        return errors.New("Cannot Access Repo API " + err.Error())
    }
    defer apiResp.Body.Close()

    // DECODE GITHUB API REQUEST
    var githubData map[string]interface{}
    if err = json.NewDecoder(apiResp.Body).Decode(&githubData); err != nil {
        return errors.New("Cannot decode Github API body to JSON : " + err.Error())
    }

    // CONTRIBUTOR DATA FROM THE REPO
    contribResp, err := http.Get(control.GetGithubAPILink(repo.RepoPage + "/contributors")); if err != nil {
        return errors.New("Cannot Access Contributors API " + err.Error())
    }
    defer contribResp.Body.Close()

    // DECODE CONTRIBUTOR DATA
    var contribData []map[string]interface{}
    if err = json.NewDecoder(contribResp.Body).Decode(&contribData); err != nil {
        return errors.New("Cannot decode Github API body to JSON : " + err.Error())
    }

    // DESCRIPTION
    description, ok := githubData["description"].(string); if !ok {
        log.Print("Cannot parse Github Description")
    }
    // DEFAULT BRANCH
    branch, ok      := githubData["default_branch"].(string); if !ok {
        return errors.New("Cannot parse repo's default branch")
    }
    // STAR COUNT
    starCount, ok   := githubData["stargazers_count"].(float64); if !ok {
        return errors.New("Cannot parse the repo's star count")
    }
    // FORKED COUNT
    forkCount, ok   := githubData["forks_count"].(float64); if !ok {
        return errors.New("Cannot parse the repo's fork count")
    }
    // WATCH COUNT
    watchCount, ok  := githubData["subscribers_count"].(float64); if !ok {
        return errors.New("Cannot parse the repo's watch count")
    }
    // UPDATED DATE
    updated, ok      := githubData["updated_at"].(string); if !ok {
        return errors.New("Cannot parse when the repo's updated")
    }
    updatedDate, err := time.Parse(time.RFC3339, updated); if err != nil {
        return errors.New("Cannot parse when the repo's created " + err.Error())
    }

    if repo.Branch != branch {
        repo.Branch      = branch
    }
    if len(description) != 0 && repo.Summary != description {
        repo.Summary = description
    }
    repo.StarCount       = int64(starCount)
    repo.ForkCount       = int64(forkCount)
    repo.WatchCount      = int64(watchCount)
    repo.Updated         = updatedDate
    db.Save(repo)

    util.GithubReadmeScrap(repo.RepoPage, repo.Slug + ".html")
    return nil
}

func main() {
    db, err := gorm.Open("sqlite3", "pc-index.db")
    if err != nil {
        log.Panic("failed to connect database " + err.Error() )
        return
    }

    log.Print("Update process started at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
/*
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
            Dial: (&net.Dialer{
                Timeout:   60 * time.Second,
                KeepAlive: 60 * time.Second,
            }).Dial,
            // We use ABSURDLY large keys, and should probably not.
            TLSHandshakeTimeout:   60 * time.Second,
            ResponseHeaderTimeout: 60 * time.Second,
            ExpectContinueTimeout: 1 * time.Second,
        },
    }
*/
    var repos []model.Repository
    db.Find(&repos)
    for _, repo := range repos {
        //log.Print("Updating " + repo.RepoPage)
        if err := AccessGithubAPI(db, &repo); err != nil {
            log.Print("Updating " + repo.RepoPage + " failed! Reason : " + err.Error())
        }
    }

    log.Print("Update process ended at " + time.Now().Format("Jan. 2 2006 3:04 PM"))
}