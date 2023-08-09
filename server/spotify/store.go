package spotify

import (
	"errors"
	"toddler-player/server/database"
)

func GetAuthData(db *database.DatabaseConnection, id string) (SpotifyAuthResponse, error) {
	var user database.User
	err := db.GetUser(id, &user)
	if err != nil {
		return SpotifyAuthResponse{}, errors.New("No session data")
	}
	return SpotifyAuthResponse{
		AccessToken:  user.Token,
		RefreshToken: user.Refresh,
	}, nil
}

func SaveAuthData(db *database.DatabaseConnection, id, spotifyId string, authData SpotifyAuthResponse) (string, error) {
	var user database.User
	if spotifyId != "" {
		db.GetUserBySpotifyId(spotifyId, &user)
	} else {
		db.GetUser(id, &user)
	}

	if user.UserID != "" {
		err := db.UpdateUser(user.UserID, authData.AccessToken, authData.RefreshToken, &user)
		return user.UserID, err
	} else {
		userId := db.CreateUser(authData.AccessToken, authData.RefreshToken, spotifyId)
		return userId, nil
	}
}
