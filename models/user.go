package models

import (
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID               string `gorm:"primarykey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        sql.NullTime `gorm:"index"`
	Username         string       `gorm:"unique"`
	PasswordHash     string
	CanResetPassword bool `gorm:"not null;default:false"`
}

func (user *User) BeforeCreate(*gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.NewString()
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (user *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
