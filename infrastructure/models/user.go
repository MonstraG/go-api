package models

import (
	"context"
	"database/sql"
	"go-api/infrastructure/crypto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               string    `gorm:"primarykey"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time
	DeletedAt        sql.NullTime `gorm:"index"`
	Username         string       `gorm:"unique"`
	PasswordHash     string       `gorm:"not null"`
	PasswordSalt     string
	CanResetPassword bool `gorm:"not null;default:false"`
}

func NewUser(username string, password string) User {
	salt := crypto.NewSalt()
	passwordHash := crypto.HashPassword(password, salt)

	return User{
		Username:     username,
		PasswordHash: passwordHash,
		PasswordSalt: salt,
	}
}

func (user *User) BeforeCreate(*gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.NewString()
	return
}

func (user *User) CheckPassword(password string) bool {
	hash := crypto.HashPassword(password, user.PasswordSalt)
	return hash == user.PasswordHash
}

func FindUser(db *gorm.DB, userId string) (User, error) {
	// todo: move this method somewhere better
	ctx := context.Background()
	return gorm.G[User](db).Where("id = ?", userId).First(ctx)
}

func (user *User) IsAdmin() bool {
	// I don't have roles yet)
	return user.Username == "MonstraG"
}
