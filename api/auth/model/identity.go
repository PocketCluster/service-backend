package model

import (
    "github.com/jinzhu/gorm"
)

type AuthIdentity struct {
    gorm.Model
    // hashed user invitation
    Invitation    string    `gorm:"column:invitation;type:VARCHAR(40) UNIQUE" sql:"index"`
    // hashed device ID
    Device        string    `gorm:"column:device;type:VARCHAR(40)"`
}

