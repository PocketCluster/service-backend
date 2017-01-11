package model

import (
    "time"
    "fmt"
)

func MakeTagEntryKey(repoID string) string {
    return fmt.Sprintf("tag-%s",repoID)
}

type RepoTag struct {
    // tag date
    Published       time.Time       `msgpack:"pubdate"`
    // tag version
    Version         string          `msgpack:"version"`
    // tag SHA
    SHA             string          `msgpack:"sha"`
    // HTML URL
    WebLink         string          `msgpack:"weblink"`
}

func (r *RepoTag) PublishedDate() string {
    return r.Published.Format(releaseDateFormat)
}

type ListTag []RepoTag

func (slice ListTag) Len() int {
    return len(slice)
}

func (slice ListTag) Less(i, j int) bool {
    return time.Duration(0) < slice[i].Published.Sub(slice[j].Published);
}

func (slice ListTag) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func (slice ListTag) FirstTenElements() []map[string]string {
    var cnt int = len(slice)
    if 10 < cnt {
        cnt = 10
    }
    // FIXME : this is ugly as mustache does not work with property function
    var list []map[string]string
    for _, tag := range slice[:cnt] {
        list = append(list, map[string]string {
            "PublishedDate":    tag.PublishedDate(),
            "Version":          tag.Version,
            "WebLink":          tag.WebLink,
        })
    }
    return list
}