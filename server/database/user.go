package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID    string `gorm:"unique"`
	Token     string
	Refresh   string
	SpotifyId string
}

func (d *DatabaseConnection) CreateUser(token, refresh, spotifyId string) string {
	id := uuid.New().String()

	d.DB.Create(&User{
		UserID:    id,
		Token:     token,
		Refresh:   refresh,
		SpotifyId: spotifyId,
	})

	return id
}

func (d *DatabaseConnection) GetUser(id string, out *User) error {
	return d.DB.Model(&User{}).Where("user_id = ?", id).First(&out).Error
}

func (d *DatabaseConnection) GetUserBySpotifyId(spotifyId string, out *User) error {
	return d.DB.Model(&User{}).Where("spotify_id = ?", spotifyId).First(&out).Error
}

func (d *DatabaseConnection) UpdateUser(id string, token, refresh string, out *User) error {
	if err := d.DB.Model(&User{}).Where("user_id = ?", id).First(&out).Error; err != nil {
		return err
	}
	out.Token = token
	out.Refresh = refresh
	return d.DB.Save(&out).Error
}
