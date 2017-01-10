package model

const RepoSuppBucket string = "repo-supp-bucket"

type RepoSupplement struct {
    RepoID          string              `msgpack:"repoid"`
    Languages       ListLanguage        `msgpack:"languages, inline, omitempty"`
    Releases        ListRelease         `msgpack:"releases, inline, omitempty"`
    Tags            ListTag             `msgpack:"tags, inline, omitempty"`
}