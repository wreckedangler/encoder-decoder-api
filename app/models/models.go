package models

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type File struct {
	ID               uint `gorm:"primaryKey"`
	OriginalFilename string
	Filename         string
	Filesize         int64
	UploadDate       time.Time
	Password         string
	FileHash         string `gorm:"uniqueIndex"`
}
