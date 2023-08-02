package database

import (
	"gorm.io/gorm"
)

type RequestLog struct {
	gorm.Model
	Url    string
	Status string
}

func (d *DatabaseConnection) LogRequest(url, status string) {
	d.DB.Create(&RequestLog{
		Url:    url,
		Status: status,
	})
}

func (d *DatabaseConnection) GetLogs() []RequestLog {
	matchingLogs := []RequestLog{}
	d.DB.Find(&matchingLogs)

	return matchingLogs
}
