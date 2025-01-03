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

	err = db.AutoMigrate(&models.Song{})
	if err != nil {
		log.Fatalf("failed to migrate songs:\n%v\n", err)
	}
	err = db.AutoMigrate(&models.SongQueueItem{})
	if err != nil {
		log.Fatalf("failed to migrate songQueues:\n%v\n", err)
	}

	seedUser(db, appConfig.DefaultUser)
	seedUser(db, appConfig.GuestUser)
}

func seedUser(db *gorm.DB, user appConfig.DefaultUser) {
	passwordHash, err := models.HashPassword(user.Password)
	if err != nil {
		log.Fatalf("Failed to hash default user %s password:\n%v\n", user.Username, err)
	}

	userModel := models.User{
		Username:     user.Username,
		PasswordHash: passwordHash,
	}

	result := db.Where(models.User{Username: user.Username}).FirstOrCreate(&userModel)

	if result.Error != nil {
		log.Fatalf("Failed to insert default user %s:\n%v\n", user.Username, err)
	}
}
