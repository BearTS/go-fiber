package services

import (
	"github.com/bearts/go-fiber/app/models"
	"github.com/golang-jwt/jwt/v4"
)

func CreateJWTTokenUser(user models.User) (string, error) {

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  "user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}
