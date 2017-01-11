package model

import "time"

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