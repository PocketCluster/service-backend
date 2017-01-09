package control

import (
    "fmt"
    "strings"
    "sort"

    log "github.com/Sirupsen/logrus"
    //"github.com/davecgh/go-spew/spew"
    "github.com/google/go-github/github"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/util"
    "github.com/gravitational/trace"
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
    releases, resp, err := ctrl.githubClient.Repositories.ListReleases(owner, repo, &github.ListOptions{Page:1, PerPage:10})
    if err != nil {
        return nil, nil, err
    }

    var listRelease model.ListRelease
    for _, rel := range releases {
        listRelease = append(listRelease, model.RepoRelease{
            Published:      util.SafeGetTimestamp(rel.PublishedAt),
            Version:        util.SafeGetString(rel.Name),
            WebLink:        util.SafeGetString(rel.HTMLURL),
        })
    }
    sort.Sort(listRelease)
    return listRelease, resp, err
}

func (ctrl *Controller) GetGithubAllTags(repoURL string, oldTagList model.ListTag) (model.ListTag, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }

    // ([]*RepositoryRelease, *Response, error) : read 26 tags due to backport of apache repositories
    ghTags, resp, err := ctrl.githubClient.Repositories.ListTags(owner, repo, &github.ListOptions{Page:0, PerPage:26})
    if err != nil {
        return nil, nil, err
    }

    getOldTag := func(prevList model.ListTag, sha string) *model.RepoTag {
        if len(prevList) == 0 {
            return nil
        }
        for _, rel := range prevList {
            if rel.SHA == sha {
                return &rel
            }
        }
        return nil
    }

    var tagList model.ListTag
    if len(ghTags) == 1 {
        tag := ghTags[0]
        SHA := util.SafeGetString(tag.Commit.SHA)
        old := getOldTag(oldTagList, SHA)
        if old == nil {
            commit, resp, err := ctrl.githubClient.Git.GetCommit(owner, repo, SHA)
            if err != nil {
                trace.Wrap(err)
                return tagList, resp, err
            }
            tagList = append(tagList, model.RepoTag{
                Published:      util.SafeGetTime(commit.Committer.Date),
                Version:        util.SafeGetString(tag.Name),
                SHA:            SHA,
                WebLink:        fmt.Sprintf("https://github.com/%s/%s/commit/%s",owner, repo, SHA),
            })
        } else {
            tagList = append(tagList, *old)
        }
        return tagList, resp, nil

    } else {
        for _, tag := range ghTags {
            SHA := util.SafeGetString(tag.Commit.SHA)
            old := getOldTag(oldTagList, SHA)

            // this tag DNE in old list
            if old == nil {
                commit, _, err := ctrl.githubClient.Git.GetCommit(owner, repo, SHA)
                if err != nil {
                    trace.Wrap(err)
                    continue
                }
                tagList = append(tagList, model.RepoTag{
                    Published:      util.SafeGetTime(commit.Committer.Date),
                    Version:        util.SafeGetString(tag.Name),
                    SHA:            SHA,
                })
            } else {
                tagList = append(tagList, *old)
            }

            if len(ghTags) <= (len(tagList) + 1) {
                break
            }
        }
    }

    if len(tagList) != 0 {
        // sort for date
        sort.Sort(tagList)
        log.Info("Let's add web link to sorted list")

        for i, _ := range tagList {
            if len(tagList[i].WebLink) == 0 {
                tagList[i].WebLink = fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s",owner, repo, tagList[i + 1].SHA, tagList[i].SHA)
            }
            if len(tagList) <= (i + 2) {
                break
            }
        }
    }

    return tagList, resp, err
}