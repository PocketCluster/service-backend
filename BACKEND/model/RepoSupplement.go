package model

import (
    "time"
    "sort"
    "strings"
)

const (
    RepoSuppBucket string       = "repo-supp-bucket"
    releaseDateFormat string    = "Jan. 2 2006"
)

// --- Representation format for users ---
type RecentPublish struct {
    Published       time.Time           `msgpack:"published"`
    Version         string              `msgpack:"version"`
    WebLink         string              `msgpack:"weblink"`
}

func converRelease(r *RepoRelease) *RecentPublish {
    var (
        published time.Time     = r.Published
        version string          = r.Version
        weblink string          = r.WebLink
    )

    if len(version) == 0 {
        var linkSplit []string = strings.Split(weblink, "/releases/tag/")
        version = linkSplit[len(linkSplit) - 1]
    }
    return &RecentPublish{
        Published:      published,
        Version:        version,
        WebLink:        weblink,
    }
}

func compareRelAndTag(r *RepoRelease, t *RepoTag) bool {
    diff := r.Published.Sub(t.Published)
    if diff < time.Duration(0) {
        diff *= -1
    }
    if len(r.TagVersion) != 0 && r.TagVersion == t.Version && diff < (time.Hour * 48) {
        return true
    }
    return false
}

func convertTag(r *RepoTag) *RecentPublish {
    return &RecentPublish{
        Published:      r.Published,
        Version:        r.Version,
        WebLink:        r.WebLink,
    }
}

func comparePubAndTag(rp *RecentPublish, t *RepoTag) bool {
    diff := rp.Published.Sub(t.Published)
    if diff < time.Duration(0) {
        diff *= -1
    }
    if strings.ToLower(rp.Version) == strings.ToLower(t.Version) && diff < (time.Hour * 48) {
        return true
    }
    return false
}

type ListPublished []*RecentPublish

func (l ListPublished) Len() int {
    return len(l)
}

func (l ListPublished) Less(i, j int) bool {
    return time.Duration(0) < l[i].Published.Sub(l[j].Published);
}

func (l ListPublished) Swap(i, j int) {
    l[i], l[j] = l[j], l[i]
}

// --- Actual Storage ---
type RepoSupplement struct {
    RepoID        string              `msgpack:"repoid"`
    Updated       time.Time           `msgpack:"updated"`
    Languages     ListLanguage        `msgpack:"languages, inline, omitempty"`
    Releases      ListRelease         `msgpack:"releases, inline, omitempty"`
    Tags          ListTag             `msgpack:"tags, inline, omitempty"`
    RecentPublish ListPublished       `msgpack:"published, inline, omitempty"`
}

func (r *RepoSupplement) BuildRecentPublication(maxEntity int) {
    var (
        pubList ListPublished
    )
    for i, _ := range r.Releases {
        pubList = append(pubList, converRelease(&(r.Releases[i])))
    }
    isRelieaseExist := func(rl *ListRelease, pl *ListPublished, t *RepoTag) bool {
        if len(*rl) == 0 {
            return false
        }
        for i, _ := range *rl {
            if compareRelAndTag(&((*rl)[i]), t) {
                return true
            }
        }
        if len(*pl) == 0 {
            return false
        }
        for _, p := range *pl {
            if comparePubAndTag(p, t) {
                return true
            }
        }
        return false
    }
    for i, _ := range r.Tags {
        if !isRelieaseExist(&(r.Releases), &pubList, &(r.Tags[i])) {
            pubList = append(pubList, convertTag(&(r.Tags[i])))
        }
    }
    var cnt int = len(pubList)
    if maxEntity < cnt {
        cnt = maxEntity
    }
    if cnt == 0 {
        return
    }
    sort.Sort(pubList)
    r.RecentPublish = pubList[:cnt]
}

func (r *RepoSupplement) RecentPublication() []map[string]string {
    if len(r.RecentPublish) == 0 {
        return nil
    }
    var list []map[string]string
    for _, pub := range r.RecentPublish {
        list = append(list, map[string]string {
            "PublishedDate":    pub.Published.Format(releaseDateFormat),
            "Version":          pub.Version,
            "WebLink":          pub.WebLink,
        })
    }
    return list
}