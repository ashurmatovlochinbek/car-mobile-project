package util

import (
	"car-mobile-project/config"
	"car-mobile-project/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	claims := &Claims{
		UserId:      user.UserId.String(),
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(user *models.User, config *config.Config) (string, error) {
	claims := &Claims{
		UserId:      user.UserId.String(),
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
