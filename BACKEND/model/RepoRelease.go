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
    // Tagged version
    TagVersion      string          `msgpack:"tagver"`
    // HTML URL
    WebLink         string          `msgpack:"weblink"`
}

type ListRelease []*RepoRelease

func (slice ListRelease) Len() int {
    return len(slice)
}

func (slice ListRelease) Less(i, j int) bool {
    return time.Duration(0) < slice[i].Published.Sub(slice[j].Published);
}

func (slice ListRelease) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}