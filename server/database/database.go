package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseConnection struct {
	DB *gorm.DB
}

func InitialiseDatabase() DatabaseConnection {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&NfcTag{}, &Status{}, &RequestLog{}, &Automation{})

	return DatabaseConnection{DB: db}
}
