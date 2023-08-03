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

func (d *DatabaseConnection) GetTag(UID string, out NfcTag) error {
	err := d.DB.Where("nfc_uid = ?", UID).First(&out).Error
	return err
}
