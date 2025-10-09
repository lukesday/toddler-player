package spotify

import (
	"errors"
)

type SpotifyAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyDevice struct {
	Id               string
	IsActive         bool `json:"is_active"`
	IsPrivateSession bool `json:"is_private_session"`
	IsRestricted     bool `json:"is_restricted"`
	Name             string
	Type             string
	VolumePercent    int  `json:"volume_percent"`
	SupportsVolume   bool `json:"supports_volume"`
}

type SpotifyUserData struct {
	SessionId   string `json:"session_id"`
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Href        string `json:"href"`
	ID          string `json:"id"`
	Product     string `json:"product"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
}

type SpotifyDeviceList struct {
	SessionId string `json:"session_id"`
	Devices   []SpotifyDevice
}

type SpotifyAuthPayload struct {
	Code     string `json:"code"`
	Redirect string `json:"redirect"`
}

var ErrInvalidToken = errors.New("Invalid Token")
var ErrAuth = errors.New("Reauthentication Required")
