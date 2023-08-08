package router

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
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
	VolumePercent    int `json:"volume_percent"`
}

type SpotifyUserData struct {
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
	Devices []SpotifyDevice
}

type SpotifyAuthPayload struct {
	Code     string `json:"code"`
	Redirect string `json:"redirect"`
}

func (r *Router) UseSpotify() {
	//Replace this with user specific cookies and access tokens - Use a session
	var spotifyAccessToken string

	r.App.Get("/api/spotify/devices", func(c *fiber.Ctx) error {
		if devices, err := getDevices(spotifyAccessToken); err == nil {
			return c.JSON(devices)
		}

		//Add atttempt to refresh in here, possibly as a wrapper to the getDevices function
		return c.SendStatus(401)
	})

	r.App.Post("/api/spotify/auth", func(c *fiber.Ctx) error {
		payload := SpotifyAuthPayload{}

		if err := c.BodyParser(&payload); err != nil {
			log.Println("error = ", err)
			return err
		}

		if token, err := getAccessToken(payload); err != nil {
			return c.SendStatus(401)
		} else {
			spotifyAccessToken = token
		}

		if userData, err := getUserData(spotifyAccessToken); err != nil {
			return err
		} else {
			return c.JSON(userData)
		}
	})
}

func getAccessToken(authData SpotifyAuthPayload) (string, error) {
	apiUrl := "https://accounts.spotify.com"
	resource := "/api/token"
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authData.Code)
	data.Set("redirect_uri", authData.Redirect)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	request.Header.Add("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var responseData SpotifyAuthResponse

	json.Unmarshal(body, &responseData)

	return responseData.AccessToken, nil
}

func getUserData(spotifyAccessToken string) (SpotifyUserData, error) {

	var responseData SpotifyUserData
	resource := "v1/me"

	body, err := spotifyGetRequest(spotifyAccessToken, resource)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func getDevices(spotifyAccessToken string) (SpotifyDeviceList, error) {

	var responseData SpotifyDeviceList
	resource := "v1/me/player/devices"

	body, err := spotifyGetRequest(spotifyAccessToken, resource)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func spotifyGetRequest(spotifyAccessToken string, resource string) ([]byte, error) {

	apiUrl := "https://api.spotify.com"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, urlStr, strings.NewReader(""))
	request.Header.Add("Authorization", "Bearer "+spotifyAccessToken)
	resp, _ := client.Do(request)
	if resp.StatusCode == 401 {
		return []byte{}, RefreshError
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

var RefreshError = errors.New("Invalid Token")
