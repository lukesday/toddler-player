package database

import (
	"errors"

	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	Key   string `gorm:"unique; not null"`
	Value string
}

func (d *DatabaseConnection) SetStatus(key string, value string) {
	matchingStatus := Status{}
	d.DB.Where("key= ?", key).First(&matchingStatus)

	if matchingStatus.Key != "" {
		if matchingStatus.Value != value {
			d.DB.Model(&matchingStatus).Update("value", value)

		}
	} else {
		d.DB.Create(&Status{
			Key:   key,
			Value: value,
		})
	}
}

func (d *DatabaseConnection) CheckStatus(key string) (string, error) {
	matchingStatus := Status{}
	d.DB.Where("key= ?", key).First(&matchingStatus)

	if matchingStatus.Key == "" {
		return "", errors.New("Status doesn't exist")
	}

	return matchingStatus.Value, nil
}
