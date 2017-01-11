package model

import (
    "time"
    "sort"
)

const (
    RepoSuppBucket string       = "repo-supp-bucket"
    releaseDateFormat string    = "Jan. 2 2006"
)

type RepoSupplement struct {
    RepoID          string              `msgpack:"repoid"`
    Updated         time.Time           `msgpack:"updated"`
    Languages       ListLanguage        `msgpack:"languages, inline, omitempty"`
    Releases        ListRelease         `msgpack:"releases, inline, omitempty"`
    Tags            ListTag             `msgpack:"tags, inline, omitempty"`
}

type pubRelease interface {
    published() time.Time
    version() string
    weblink() string
}

type listPublished []pubRelease

func (l listPublished) Len() int {
    return len(l)
}

func (l listPublished) Less(i, j int) bool {
    return time.Duration(0) < l[i].published().Sub(l[j].published());
}

func (l listPublished) Swap(i, j int) {
    l[i], l[j] = l[j], l[i]
}

func (r *RepoSupplement) RecentPublication() []map[string]string {

    var pubList listPublished
    for i, _ := range r.Releases {
        pubList = append(pubList, &(r.Releases[i]))
    }
    for i, _ := range r.Tags {
        pubList = append(pubList, &(r.Tags[i]))
    }
    var cnt int = len(pubList)
    if 10 < cnt {
        cnt = 10
    }
    if cnt == 0 {
        return nil
    }

    sort.Sort(pubList)

    // FIXME : this is ugly as mustache does not work with property function
    var list []map[string]string
    for _, pub := range pubList[:cnt] {
        list = append(list, map[string]string {
            "PublishedDate":    pub.published().Format(releaseDateFormat),
            "Version":          pub.version(),
            "WebLink":          pub.weblink(),
        })
    }
    return list
}