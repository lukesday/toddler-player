package database

import "gorm.io/gorm"

type Automation struct {
	gorm.Model
	NfcTagID int
	NfcTag   NfcTag
	MediaId  string
	Name     string
	UserID   string
	User     User `gorm:"references:UserID"`
}

func (d *DatabaseConnection) CreateAutomation(nfcTag NfcTag, mediaId, name string, user User) {
	d.DB.Create(&Automation{
		NfcTag:  nfcTag,
		MediaId: mediaId,
		Name:    name,
		User:    user,
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

func (d *DatabaseConnection) DeleteAutomation(id uint) error {
	return d.DB.Delete(&Automation{}, id).Error
}
