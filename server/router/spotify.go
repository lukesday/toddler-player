package router

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SpotifyDevice struct {
	Id               string
	IsActive         bool `json:"is_active"`
	IsPrivateSession bool `json:"is_private_session"`
	IsRestricted     bool `json:"is_restricted"`
	Name             string
	Type             string
	VolumePercent    int `json:"volume_percent"`
}

type SpotifyDeviceList struct {
	Devices []SpotifyDevice
}

func (r *Router) UseSpotify() {
	SpotifyAccessToken := getAccessToken()

	r.App.Get("/api/spotify/devices", func(c *fiber.Ctx) error {
		return c.JSON(getDevices(SpotifyAccessToken))
	})
}

func getAccessToken() string {
	apiUrl := "https://accounts.spotify.com"
	resource := "/api/token"
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	request.Header.Add("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)

	var responseData SpotifyAuthResponse

	json.Unmarshal(body, &responseData)

	return responseData.AccessToken
}

func getDevices(spotifyAccessToken string) SpotifyDeviceList {
	apiUrl := "https://api.spotify.com"
	resource := "v1/me/player/devices"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, urlStr, strings.NewReader(""))
	request.Header.Add("Authorization", "Bearer "+spotifyAccessToken)
	resp, _ := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)

	var responseData SpotifyDeviceList

	json.Unmarshal(body, &responseData)

	return responseData
}
