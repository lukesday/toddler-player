package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	body, err := spotifyRequestWithRetry(authData, resource, id, db, http.MethodGet, nil)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func GetUserDataWithoutRetry(authData SpotifyAuthResponse) (SpotifyUserData, error) {

	var responseData SpotifyUserData
	resource := "v1/me"

	body, err := spotifyRequest(authData.AccessToken, resource, http.MethodGet, nil)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func GetDevices(authData SpotifyAuthResponse, id string, db *database.DatabaseConnection) (SpotifyDeviceList, error) {

	var responseData SpotifyDeviceList
	resource := "v1/me/player/devices"

	body, err := spotifyRequestWithRetry(authData, resource, id, db, http.MethodGet, nil)

	json.Unmarshal(body, &responseData)

	return responseData, err
}

func StartPlayback(authData SpotifyAuthResponse, id string, db *database.DatabaseConnection, contextUri string) error {
	resource := "v1/me/player/play"

	_, err := spotifyRequestWithRetry(authData, resource, id, db, http.MethodPut, strings.NewReader(fmt.Sprintf(`{"context_uri": "%s"}`, contextUri)))

	return err
}

func spotifyRequest(spotifyAccessToken, resource, method string, body io.Reader) ([]byte, error) {

	apiUrl := "https://api.spotify.com"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	request, _ := http.NewRequest(method, urlStr, body)
	request.Header.Add("Authorization", "Bearer "+spotifyAccessToken)
	resp, _ := client.Do(request)

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		return []byte{}, ErrInvalidToken
	}

	return responseBody, nil
}

func spotifyRequestWithRetry(authData SpotifyAuthResponse, resource, id string, db *database.DatabaseConnection, method string, body io.Reader) ([]byte, error) {
	respData, err := spotifyRequest(authData.AccessToken, resource, method, body)

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

	authBody, _ := io.ReadAll(resp.Body)

	json.Unmarshal(authBody, &authResponse)

	SaveAuthData(db, id, "", authResponse)

	retryData, err := spotifyRequest(authData.AccessToken, resource, http.MethodGet, body)

	return retryData, err
}
