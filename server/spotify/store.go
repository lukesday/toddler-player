package spotify

import (
	"errors"
	"toddler-player/server/database"

	"github.com/google/uuid"
)

func GetAuthData(db *database.DatabaseConnection, id string) (SpotifyAuthResponse, error) {
	userId := uuid.UUID([]byte(id))
	var user database.User
	err := db.GetUser(userId, user)
	if err != nil {
		return SpotifyAuthResponse{}, errors.New("No session data")
	}
	return SpotifyAuthResponse{
		AccessToken:  user.Token,
		RefreshToken: user.Refresh,
	}, nil
}

func SaveAuthData(db *database.DatabaseConnection, id string, authData SpotifyAuthResponse) (string, error) {
	userId := uuid.UUID([]byte(id))
	var user database.User
	err := db.GetUser(userId, user)
	if err == nil {
		err := db.UpdateUser(user.UserID, authData.AccessToken, authData.RefreshToken, user)
		return user.UserID.String(), err
	} else {
		userId = db.CreateUser(authData.AccessToken, authData.RefreshToken)
		return userId.String(), nil
	}
}
