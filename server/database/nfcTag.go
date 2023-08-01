package database

import "gorm.io/gorm"

type NfcTag struct {
	gorm.Model
	NfcUID string `gorm:"unique; not null"`
}

func (d *DatabaseConnection) Enroll(UID string) {
	d.DB.Create(&NfcTag{
		NfcUID: UID,
	})
}

func (d *DatabaseConnection) Exists(UID string) bool {
	matchingTags := []NfcTag{}
	d.DB.Where("nfc_uid = ?", UID).Find(&matchingTags)
	return len(matchingTags) > 0
}
