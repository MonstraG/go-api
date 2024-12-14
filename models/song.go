package models

import (
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	YoutubeId string `gorm:"unique"`
}
