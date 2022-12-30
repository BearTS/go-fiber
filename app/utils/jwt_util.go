package utils

import (
	"context"
	"time"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Token struct {
	AccessToken struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"accessToken"`
	RefreshToken struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"refreshToken"`
}

var RefreshTokenCollection *mongo.Collection = database.GetCollection(database.DB, "RefreshToken")

func CreateJWTTokenUser(user models.User) (*Token, error) {
	PhoneAvailable := false
	if user.Phone != "" {
		PhoneAvailable = true
	}
	claims := jwt.MapClaims{
		"id":             user.Id,
		"email":          user.Email,
		"role":           "user",
		"PhoneAvailable": PhoneAvailable,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte("secret"))
	accessTokenExpiresAt := time.Now().Add(time.Hour * 24).Unix()
	if err != nil {
		return nil, err
	}
	refreshToken := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24 * 3).Unix()
	obj := models.RefreshToken{
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		User:      user.Id,
	}
	if _, err := RefreshTokenCollection.InsertOne(context.Background(), obj); err != nil {
		return nil, err
	}
	return &Token{
		AccessToken: struct {
			Token     string `json:"token"`
			ExpiresAt int64  `json:"expiresAt"`
		}{
			Token:     accessToken,
			ExpiresAt: accessTokenExpiresAt,
		},
		RefreshToken: struct {
			Token     string `json:"token"`
			ExpiresAt int64  `json:"expiresAt"`
		}{
			Token:     refreshToken,
			ExpiresAt: expiresAt,
		},
	}, nil

}
