package model

import (
    "time"
)

func MakeTagEntryKey(repoID string) string {
    return "tag-" + repoID
}

type RepoTag struct {
    // release date
    Published       time.Time       `msgpack:"pubdate"`
    // release name
    Name            string          `msgpack:"name"`
    // release note
    Note            string          `msgpack:"note"`
    // SHA
    SHA             string          `msgpack:"sha"`
    // HTML URL
    WebLink         string          `msgpack:"weblink"`
}

type ListTag []RepoTag

func (slice ListTag) Len() int {
    return len(slice)
}

func (slice ListTag) Less(i, j int) bool {
    return 0 < slice[i].Published.Sub(slice[j].Published);
}

func (slice ListTag) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

