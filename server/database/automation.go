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

func (d *DatabaseConnection) GetAutomation(nfcTag NfcTag, out Automation) error {
	err := d.DB.Model(&Automation{}).Preload("NfcTag").Where("nfc_tag = ?", nfcTag).First(&out).Error
	return err
}

func (d *DatabaseConnection) ListAutomations(nfcTag NfcTag) ([]Automation, error) {
	var automations []Automation
	err := d.DB.Model(&Automation{}).Preload("NfcTag").Find(&automations).Error
	return automations, err
}
