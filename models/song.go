package models

import (
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	YoutubeId string `gorm:"unique"`
	Duration  int
}

type SongQueueItem struct {
	gorm.Model
	SongId uint
	Song   Song
}
