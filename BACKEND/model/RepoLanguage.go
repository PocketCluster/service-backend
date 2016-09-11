package model

import (
    "github.com/jinzhu/gorm"
)

type RepoLanguage struct {
    gorm.Model
    // repository ID
    RepoId            string
    // Programming Language
    Language          string
    // Percentage
    Percentage        float64
}