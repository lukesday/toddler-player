package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
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

	body, _ := io.ReadAll(resp.Body)

	json.Unmarshal(body, &responseData)

	return responseData, nil
}

func GetUserData(authData SpotifyAuthResponse, id string, db *database.DatabaseConnection) (SpotifyUserData, error) {

	var responseData SpotifyUserData
	resource := "v1/me"

	body, err := spotifyGetRequestWithRetry(authData, resource, id, db)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func GetUserDataWithoutRetry(authData SpotifyAuthResponse) (SpotifyUserData, error) {

	var responseData SpotifyUserData
	resource := "v1/me"

	body, err := spotifyGetRequest(authData.AccessToken, resource)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func GetDevices(authData SpotifyAuthResponse, id string, db *database.DatabaseConnection) (SpotifyDeviceList, error) {

	var responseData SpotifyDeviceList
	resource := "v1/me/player/devices"

	body, err := spotifyGetRequestWithRetry(authData, resource, id, db)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func spotifyGetRequest(spotifyAccessToken, resource string) ([]byte, error) {

	apiUrl := "https://api.spotify.com"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, urlStr, strings.NewReader(""))
	request.Header.Add("Authorization", "Bearer "+spotifyAccessToken)
	resp, _ := client.Do(request)

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		return []byte{}, ErrInvalidToken
	}

	return body, nil
}

func spotifyGetRequestWithRetry(authData SpotifyAuthResponse, resource, id string, db *database.DatabaseConnection) ([]byte, error) {
	respData, err := spotifyGetRequest(authData.AccessToken, resource)

	if err == nil {
		return respData, err
	}

	if !errors.Is(err, ErrInvalidToken) {
		return respData, err
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
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return []byte{}, ErrAuth
	}

	body, _ := io.ReadAll(resp.Body)

	json.Unmarshal(body, &authResponse)

	SaveAuthData(db, id, "", authResponse)

	retryData, err := spotifyGetRequest(authData.AccessToken, resource)

	return retryData, err
}
