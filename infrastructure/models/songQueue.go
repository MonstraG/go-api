package models

import (
	"time"
)

type QueuedSong struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Path      string
}
