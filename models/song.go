package models

import (
	"gorm.io/gorm"
	"time"
)

type Song struct {
	gorm.Model
	YoutubeId string `gorm:"unique"`
	Duration  int
	Title     string
	File      string
}

type SongQueueItem struct {
	gorm.Model
	SongId   uint
	Song     Song
	StartsAt time.Time
	EndsAt   time.Time
}
