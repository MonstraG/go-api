package models

import (
	"context"
	"database/sql"
	"go-api/infrastructure/pass"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID `gorm:"primarykey"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time
	DeletedAt        sql.NullTime `gorm:"index"`
	Username         string       `gorm:"unique"`
	PasswordHash     string       `gorm:"not null"`
	PasswordSalt     string
	CanResetPassword bool `gorm:"not null;default:false"`
	IsAdmin          bool `gorm:"not null;default:false"`
}

func NewUser(username string, password string) User {
	salt := pass.NewSalt()
	passwordHash := pass.HashPassword(password, salt)

	return User{
		Username:     username,
		PasswordHash: passwordHash,
		PasswordSalt: salt,
	}
}

func (user *User) BeforeCreate(*gorm.DB) (err error) {
	user.ID = uuid.NewV7()
	return
}

func (user *User) CheckPassword(password string) bool {
	hash := pass.HashPassword(password, user.PasswordSalt)
	return hash == user.PasswordHash
}

func FindUser(db *gorm.DB, userId uuid.UUID) (User, error) {
	// todo: move this method somewhere better
	ctx := context.Background()
	return gorm.G[User](db).Where("id = ?", userId).First(ctx)
}

func (user *User) IsAdminOld() bool {
	// todo: make this a flag
	return user.Username == "MonstraG"
}
