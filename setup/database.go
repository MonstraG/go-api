package setup

import (
	"github.com/glebarez/sqlite"
	"go-server/models"
	"go-server/setup/appConfig"
	"gorm.io/gorm"
	"log"
)

const foreignKeySwitch = "?_pragma=foreign_keys(1)"

func OpenDb(appConfig appConfig.AppConfig) *gorm.DB {
	if len(appConfig.DatabaseFile) == 0 {
		log.Fatalf("Database file not specified")
	}

	dsn := appConfig.DatabaseFile + foreignKeySwitch
	log.Printf("Opening database %s", dsn)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database:\n%v\n", err)
	}

	seedDb(db, appConfig)

	return db
}

func seedDb(db *gorm.DB, appConfig appConfig.AppConfig) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to migrate users:\n%v\n", err)
	}

	passwordHash, err := models.HashPassword(appConfig.DefaultUser.Password)
	if err != nil {
		log.Fatalf("Failed to hash default user password:\n%v\n", err)
	}

	db.FirstOrCreate(&models.User{
		Username:     appConfig.DefaultUser.Username,
		PasswordHash: passwordHash,
	})
}
