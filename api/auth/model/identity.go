package model

import (
    "github.com/jinzhu/gorm"
)

type AuthIdentity struct {
    gorm.Model
    // hashed user invitation
    UserID    string    `gorm:"column:user_id;type:VARCHAR(256) UNIQUE" sql:"index"`
    // hashed device ID
    DeviceID  string    `gorm:"column:device_id;type:VARCHAR(256)"`
}

