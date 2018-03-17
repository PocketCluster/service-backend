package model

import (
    "regexp"

    "github.com/jinzhu/gorm"
)

const (
    ColEmail         string = "email"
    ColInvitation    string = "invitation"
    ColInvHash       string = "invhash"
    ColDevHash       string = "devhash"

    RegHashChecker   string = "^[a-z0-9]{40}$"
)

type AuthIdentity struct {
    gorm.Model
    // user email
    Email         string    `gorm:"column:email;type:VARCHAR(256)"`
    // user invitation
    Invitation    string    `gorm:"column:invitation;type:VARCHAR(24)"`
    // hashed user invitation
    InvHash       string    `gorm:"column:invhash;type:VARCHAR(40)"`
    // hashed device ID
    DevHash       string    `gorm:"column:devhash;type:VARCHAR(40)"`
}

func (AuthIdentity) TableName() string {
    return "auth_identity"
}

func IsValidHash(hashstr string) bool {
    if len(hashstr) != 40 {
        return false
    }
    match, _ := regexp.MatchString(RegHashChecker, hashstr)
    return match
}