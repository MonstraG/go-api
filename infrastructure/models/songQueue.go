package models

import (
	"time"
)

type QueuedSong struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Path      string
	Duration  time.Duration
	EndsAt    time.Time
}
