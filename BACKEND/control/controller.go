package control

import (
    "github.com/gorilla/sessions"
    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"
    "github.com/stkim1/BACKEND/model"
    "strings"
)

type Controller struct {
}

func (controller *Controller) GetSession(c web.C) *sessions.Session {
    return c.Env["Session"].(*sessions.Session)
}

func (controller *Controller) GetGORM(c web.C) *gorm.DB {
    return c.Env["GORM"].(*gorm.DB)
}

func (controller *Controller) IsXhr(c web.C) bool {
    return c.Env["IsXhr"].(bool)
}


/*
    repo1, repo2, repo3 := GetAssignedRepoColumn(len(repositories))
    for index, _ := range repositories {
        subindex := index % SingleColumnCount
        switch int(index / SingleColumnCount) {
            case 0: {
                repo1[subindex] = &repositories[index]
                break
            }
            case 1: {
                repo2[subindex] = &repositories[index]
                break
            }
            case 2: {
                repo3[subindex] = &repositories[index]
                break
            }
        }
    }
*/

// Assignment of repo to column within a page
const SingleColumnCount int = 10
const TotalRowCount int = 3

func GetAssignedRepoColumn(repoCount int) ([]*model.Repository, []*model.Repository, []*model.Repository) {

    remainCount := repoCount
    remainCheck := func() int {
        if remainCount <= 0 {
            return 0
        } else if 0 < remainCount && remainCount <= SingleColumnCount {
            count := remainCount
            remainCount = 0
            return count
        } else {
            remainCount -= SingleColumnCount
            return SingleColumnCount
        }
    }

    remain := remainCheck()
    var repo1, repo2, repo3 []*model.Repository
    if remain != 0 {
        repo1 = make([]*model.Repository, remain)
    }
    remain = remainCheck()
    if remain != 0 {
        repo2 = make([]*model.Repository, remain)
    }
    remain = remainCheck()
    if remain != 0 {
        repo3 = make([]*model.Repository, remain)
    }
    return repo1, repo2, repo3
}

const GithubClientIdentity string = "c74abcf03e61e209b3c3"
const GithubClientSecret string = "da0f7d33d02552282e72a7e594d39ba76f96d478"

func GetGithubAPILink(githubLink string) string {
    URL := strings.Replace(githubLink , "https://github.com/", "https://api.github.com/repos/", -1)
    URL += "?client_id=" + GithubClientIdentity + "&client_secret=" + GithubClientSecret
    return URL
}