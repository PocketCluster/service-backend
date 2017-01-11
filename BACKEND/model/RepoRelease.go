package model

import (
    "time"
    "fmt"
)

func MakeReleaseEntryKey(repoID string) string {
    return fmt.Sprintf("release-%s",repoID)
}

type RepoRelease struct {
    // release date
    Published       time.Time       `msgpack:"pubdate"`
    // release version
    Version         string          `msgpack:"version"`
    // HTML URL
    WebLink         string          `msgpack:"weblink"`
}

func (r *RepoRelease) PublishedDate() string {
    return r.Published.Format(releaseDateFormat)
}

type ListRelease []RepoRelease

func (slice ListRelease) Len() int {
    return len(slice)
}

func (slice ListRelease) Less(i, j int) bool {
    return time.Duration(0) < slice[i].Published.Sub(slice[j].Published);
}

func (slice ListRelease) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func (slice ListRelease) FirstTenElements() []map[string]string {
    var cnt int = len(slice)
    if 10 < cnt {
        cnt = 10
    }
    // FIXME : this is ugly as mustache does not work with property function
    var list []map[string]string
    for _, rel := range slice[:cnt] {
        list = append(list, map[string]string {
            "PublishedDate":    rel.PublishedDate(),
            "Version":          rel.Version,
            "WebLink":          rel.WebLink,
        })
    }
    return list
}