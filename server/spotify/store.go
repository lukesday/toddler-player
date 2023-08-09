package spotify

import (
	"errors"
	"log"
	"toddler-player/server/database"

	"github.com/google/uuid"
)

func GetAuthData(db *database.DatabaseConnection, id string) (SpotifyAuthResponse, error) {
	userId := parseUUID(id)
	var user database.User
	err := db.GetUser(userId, &user)
	if err != nil {
		return SpotifyAuthResponse{}, errors.New("No session data")
	}
	return SpotifyAuthResponse{
		AccessToken:  user.Token,
		RefreshToken: user.Refresh,
	}, nil
}

func SaveAuthData(db *database.DatabaseConnection, id, spotifyId string, authData SpotifyAuthResponse) (string, error) {
	userId := parseUUID(id)
	var user database.User
	if spotifyId != "" {
		db.GetUserBySpotifyId(spotifyId, &user)
	} else {
		db.GetUser(userId, &user)
	}

	log.Print(user)

	if user.UserID != uuid.Nil {
		err := db.UpdateUser(user.UserID, authData.AccessToken, authData.RefreshToken, &user)
		return user.UserID.String(), err
	} else {
		userId = db.CreateUser(authData.AccessToken, authData.RefreshToken, spotifyId)
		return userId.String(), nil
	}
}

func parseUUID(id string) uuid.UUID {
	if id == "" {
		return uuid.Nil
	}

	return uuid.UUID([]byte(id))
}
