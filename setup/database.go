package setup

import (
	"github.com/glebarez/sqlite"
	"go-api/models"
	"go-api/setup/appConfig"
	"go-api/setup/myLog"
	"gorm.io/gorm"
)

const foreignKeySwitch = "?_pragma=foreign_keys(1)"

func OpenDb(appConfig appConfig.AppConfig) *gorm.DB {
	if len(appConfig.DatabaseFile) == 0 {
		myLog.Fatal.Log("Database file not specified")
	}

	dsn := appConfig.DatabaseFile + foreignKeySwitch
	myLog.Info.Logf("Opening database %s", dsn)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		myLog.Fatal.Logf("failed to open database:\n%v\n", err)
	}

	seedDb(db, appConfig)

	return db
}

func seedDb(db *gorm.DB, appConfig appConfig.AppConfig) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		myLog.Fatal.Logf("failed to migrate users:\n%v\n", err)
	}

	seedUser(db, appConfig.DefaultUser)
	seedUser(db, appConfig.GuestUser)
}

func seedUser(db *gorm.DB, user appConfig.DefaultUser) {
	passwordHash, err := models.HashPassword(user.Password)
	if err != nil {
		myLog.Fatal.Logf("Failed to hash default user %s password:\n%v\n", user.Username, err)
	}

	userModel := models.User{
		Username:     user.Username,
		PasswordHash: passwordHash,
	}

	result := db.Where(models.User{Username: user.Username}).FirstOrCreate(&userModel)

	if result.Error != nil {
		myLog.Fatal.Logf("Failed to insert default user %s:\n%v\n", user.Username, err)
	}
}
