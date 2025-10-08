package setup

import (
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/models"
	"go-api/infrastructure/myLog"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const foreignKeySwitch = "?_pragma=foreign_keys(1)"

func OpenDb(appConfig appConfig.AppConfig) *gorm.DB {
	if len(appConfig.DatabaseFile) == 0 {
		myLog.Fatal.Logf("Database file not specified")
	}

	dsn := appConfig.DatabaseFile + foreignKeySwitch
	myLog.Info.Logf("Opening database %s", dsn)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		myLog.Fatal.Logf("failed to open database:\n\t%v", err)
	}

	seedDb(db, appConfig)

	return db
}

func seedDb(db *gorm.DB, appConfig appConfig.AppConfig) {
	err := db.AutoMigrate(&models.User{}, &models.QueuedSong{})
	if err != nil {
		myLog.Fatal.Logf("failed to migrate users:\n\t%v", err)
	}

	seedUser(db, appConfig.DefaultUser)
	seedUser(db, appConfig.GuestUser)
}

func seedUser(db *gorm.DB, user appConfig.DefaultUser) {
	passwordHash, err := models.HashPassword(user.Password)
	if err != nil {
		myLog.Fatal.Logf("Failed to hash default user %s password:\n\t%v", user.Username, err)
	}

	userModel := models.User{
		Username:     user.Username,
		PasswordHash: passwordHash,
	}

	result := db.Where(models.User{Username: user.Username}).FirstOrCreate(&userModel)

	if result.Error != nil {
		myLog.Fatal.Logf("Failed to insert default user %s:\n\t%v", user.Username, err)
	}
}
