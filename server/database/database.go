package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	//"gorm.io/gorm/logger"
)

type DatabaseConnection struct {
	DB *gorm.DB
}

func InitialiseDatabase() DatabaseConnection {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&NfcTag{},
		&Status{},
		&RequestLog{},
		&Automation{},
		&User{},
	)

	return DatabaseConnection{DB: db}
}
