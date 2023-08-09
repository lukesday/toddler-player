package database

import "gorm.io/gorm"

type Automation struct {
	gorm.Model
	NfcTag   NfcTag `gorm:"foreignKey:NfcUID"`
	DeviceId string
	MediaId  string
}

func (d *DatabaseConnection) CreateAutomation(nfcTag NfcTag, deviceId, mediaId string) {
	d.DB.Create(&Automation{
		NfcTag:   nfcTag,
		DeviceId: deviceId,
		MediaId:  mediaId,
	})
}

func (d *DatabaseConnection) GetAutomation(id uint, out *Automation) error {
	return d.DB.Model(&Automation{}).Preload("NfcTag").Where("id = ?", id).First(&out).Error
}

func (d *DatabaseConnection) GetAutomationByNfcTag(nfcTag NfcTag, out *Automation) error {
	return d.DB.Model(&Automation{}).Preload("NfcTag").Where("nfc_tag = ?", nfcTag).First(&out).Error
}

func (d *DatabaseConnection) ListAutomations(out *[]Automation) error {
	return d.DB.Model(&Automation{}).Preload("NfcTag").Find(&out).Error
}
