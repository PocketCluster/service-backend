package model

import (
    "github.com/jinzhu/gorm"
)

/*
  "owner": {
    "login": "rqlite",
    "id": 18683023,
    "avatar_url": "https://avatars.githubusercontent.com/u/18683023?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/rqlite",
    "html_url": "https://github.com/rqlite",
    "followers_url": "https://api.github.com/users/rqlite/followers",
    "following_url": "https://api.github.com/users/rqlite/following{/other_user}",
    "gists_url": "https://api.github.com/users/rqlite/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/rqlite/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/rqlite/subscriptions",
    "organizations_url": "https://api.github.com/users/rqlite/orgs",
    "repos_url": "https://api.github.com/users/rqlite/repos",
    "events_url": "https://api.github.com/users/rqlite/events{/privacy}",
    "received_events_url": "https://api.github.com/users/rqlite/received_even ts",
    "type": "Organization",
    "site_admin": false
  }
 */
type Author struct{
    gorm.Model
    // Is this from Github/Gitlab/Bitbucket
    Service         string
    // Profile Type : Organization/Personal/etc
    Type            string
    // two abbreviate chars + numbering : gh23247808
    AuthorId        string
    // Author Full Name
    Login           string
    // Author Full Name
    Name            string

    // home Page Link
    HomeURL         string
    // service profile Page
    ProfileURL      string
    // Avatar Link
    AvatarURL       string
    // If this author is deceased
    Deceased        bool
}
