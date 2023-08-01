package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type NfcTag struct {
	gorm.Model
	NfcUID string `gorm:"unique; not null"`
}

type DatabaseConnection struct {
	DB *gorm.DB
}

func InitialiseDatabase() DatabaseConnection {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&NfcTag{})

	return DatabaseConnection{DB: db}
}

func (d *DatabaseConnection) Enroll(UID string) {
	d.DB.Create(&NfcTag{
		NfcUID: UID,
	})
}
