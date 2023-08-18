package database

import "gorm.io/gorm"

type Automation struct {
	gorm.Model
	NfcTagID int
	NfcTag   NfcTag
	DeviceId string
	MediaId  string
	UserID   string
	User     User `gorm:"references:UserID"`
}

func (d *DatabaseConnection) CreateAutomation(nfcTag NfcTag, deviceId, mediaId string, user User) {
	d.DB.Create(&Automation{
		NfcTag:   nfcTag,
		DeviceId: deviceId,
		MediaId:  mediaId,
		User:     user,
	})
}

func (d *DatabaseConnection) GetAutomation(id uint, out *Automation) error {
	return d.DB.Model(&Automation{}).Preload("NfcTag").Where("id = ?", id).First(&out).Error
}

func (d *DatabaseConnection) GetAutomationByNfcTag(nfcTag NfcTag, out *Automation) error {
	return d.DB.Model(&Automation{}).Preload("NfcTag").Where("nfc_tag_ID = ?", nfcTag.ID).First(&out).Error
}

func (d *DatabaseConnection) ListAutomations(userId string, out *[]Automation) error {
	return d.DB.Model(&Automation{}).Preload("NfcTag").Where("user_id = ?", userId).Find(&out).Error
}
