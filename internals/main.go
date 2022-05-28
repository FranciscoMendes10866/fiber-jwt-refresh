package internals

import (
	"go-refresh/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppInternalsInterface interface {
	Connect() error
}

type appInternals struct{}

func NewAppInternals() AppInternalsInterface {
	return &appInternals{}
}

var Database *gorm.DB

func (*appInternals) Connect() error {
	var err error
	Database, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	Database.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
	return nil
}

// This is the AppInternals instance
var AppInternals AppInternalsInterface = NewAppInternals()
