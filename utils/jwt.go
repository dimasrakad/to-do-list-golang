package utils

import (
	"time"
	"to-do-list-golang/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

var cfg = config.LoadConfig()
var JwtKey = []byte(cfg.JWTKey)

func GenerateJWT(userID uint) (string, error) {

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
