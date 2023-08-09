package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"toddler-player/server/database"
)

func GetAccessToken(authData SpotifyAuthPayload) (SpotifyAuthResponse, error) {
	apiUrl := "https://accounts.spotify.com"
	resource := "/api/token"
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	var responseData SpotifyAuthResponse

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
		return responseData, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &responseData)

	return responseData, nil
}

func GetUserData(authData SpotifyAuthResponse, id string, db *database.DatabaseConnection) (SpotifyUserData, error) {

	var responseData SpotifyUserData
	resource := "v1/me"

	body, err, authData := spotifyGetRequestWithRetry(authData, resource)
	if authData.AccessToken != "" {
		SaveAuthData(db, id, authData)
	}

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func GetDevices(authData SpotifyAuthResponse, id string, db *database.DatabaseConnection) (SpotifyDeviceList, error) {

	var responseData SpotifyDeviceList
	resource := "v1/me/player/devices"

	body, err, authData := spotifyGetRequestWithRetry(authData, resource)
	if authData.AccessToken != "" {
		SaveAuthData(db, id, authData)
	}

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

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		return []byte{}, InvalidToken
	}

	return body, nil
}

func spotifyGetRequestWithRetry(authData SpotifyAuthResponse, resource string) ([]byte, error, SpotifyAuthResponse) {
	respData, err := spotifyGetRequest(authData.AccessToken, resource)

	if err == nil {
		return respData, err, SpotifyAuthResponse{}
	}

	if !errors.Is(err, InvalidToken) {
		return respData, err, SpotifyAuthResponse{}
	}

	apiUrl := "https://accounts.spotify.com"
	authResource := "/api/token"
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	var authResponse SpotifyAuthResponse

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", authData.RefreshToken)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = authResource
	urlStr := u.String()

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	request.Header.Add("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		return []byte{}, err, SpotifyAuthResponse{}
	}
	if resp.StatusCode != 200 {
		return []byte{}, AuthError, SpotifyAuthResponse{}
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &authResponse)

	retryData, err := spotifyGetRequest(authData.AccessToken, resource)

	return retryData, err, authResponse
}
