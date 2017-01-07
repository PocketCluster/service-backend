package control

import (
    "fmt"
    "strings"
    "sort"

    //log "github.com/Sirupsen/logrus"
    //"github.com/davecgh/go-spew/spew"
    "github.com/google/go-github/github"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
)

func (ctrl *Controller) GetGithubRepoMeta(repoURL string) (*github.Repository, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }
    return ctrl.githubClient.Repositories.Get(owner, repo)
}

func (ctrl *Controller) GetGithubContributors(repoURL string) ([]*github.Contributor, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }
    // We'll execlude anonymous users as it doesn't provide much information
    // https://developer.github.com/v3/repos/#list-contributors
    //opts := &github.ListContributorsOptions{Anon: "true"}
    opts := &github.ListContributorsOptions{}
    return ctrl.githubClient.Repositories.ListContributors(owner, repo, opts)
}

func (ctrl *Controller) GetGithubContributorsStat(repoURL string) ([]*github.ContributorStats, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }
    return ctrl.githubClient.Repositories.ListContributorsStats(owner, repo)
}

func (ctrl *Controller) GetGithubRepoLanguages(repoURL string) (model.ListLanguage, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }
    // (map[string]int, *Response, error)
    languages, resp, err := ctrl.githubClient.Repositories.ListLanguages(owner, repo)
    if err != nil {
        return nil, nil, err
    }

    // count
    var total int64 = 0
    for _, c := range languages {
        total += int64(c)
    }
    // change
    var listLang model.ListLanguage
    for lang, count := range languages {
        listLang = append(listLang, model.RepoLanguage{
            Language: lang,
            Percentage: float32(float32(count)/float32(total)),
        })
    }
    sort.Sort(listLang)
    return listLang, resp, err
}

func (ctrl *Controller) GetGithubAllReleases(repoURL string) (model.ListRelease, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }

    // ([]*RepositoryRelease, *Response, error)
    // TODO : iterate to get all the releases
    releases, resp, err := ctrl.githubClient.Repositories.ListReleases(owner, repo, &github.ListOptions{Page:1, PerPage:100})
    if err != nil {
        return nil, nil, err
    }

    var listRelease model.ListRelease
    for _, rel := range releases {
        listRelease = append(listRelease, model.RepoRelease{
            Published:      util.SafeGetTimestamp(rel.PublishedAt),
            Name:           util.SafeGetString(rel.Name),
            Note:           util.SafeGetString(rel.Body),
            WebLink:        util.SafeGetString(rel.HTMLURL),
        })
    }
    sort.Sort(listRelease)
    return listRelease, resp, err
}

func (ctrl *Controller) GetGithubAllTags(repoURL string) (model.ListTag, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }

    // ([]*RepositoryRelease, *Response, error)
    // TODO : iterate to get all the releases
    tags, resp, err := ctrl.githubClient.Repositories.ListTags(owner, repo, &github.ListOptions{Page:1, PerPage:100})
    if err != nil {
        return nil, nil, err
    }

    //log.Info(spew.Sdump(tags))

    var listTag model.ListTag
    for _, tag := range tags {
        listTag = append(listTag, model.RepoTag{
//            Published:      util.SafeGetTime(tag.Commit.Author.Date),
            Name:           util.SafeGetString(tag.Name),
//            Note:           util.SafeGetString(tag.Commit.Message),
            SHA:            util.SafeGetString(tag.Commit.SHA),
            WebLink:        util.SafeGetString(tag.Commit.URL),
        })
    }
    sort.Sort(listTag)
    return listTag, resp, err
}
