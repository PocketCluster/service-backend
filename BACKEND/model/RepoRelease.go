package model

import (
    "time"
)

func MakeReleaseEntryKey(repoID string) []byte {
    return []byte("release-" + repoID)
}

type RepoRelease struct {
    // release date
    Published       time.Time       `msgpack:"pubdate"`
    // release name
    Name            string          `msgpack:"name"`
    // release note
    Note            string          `msgpack:"note"`
    // HTML URL
    WebLink         string          `msgpack:"weblink"`
}

type ListRelease []RepoRelease

func (slice ListRelease) Len() int {
    return len(slice)
}

func (slice ListRelease) Less(i, j int) bool {
    return time.Duration(0) < slice[j].Published.Sub(slice[i].Published);
}

func (slice ListRelease) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}
