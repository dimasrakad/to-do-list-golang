package utils

import (
	"log"
	"strconv"
	"time"
	"to-do-list-golang/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

var cfg = config.LoadConfig()
var (
	AccessSecret  = []byte(cfg.JWTAccessSecret)
	RefreshSecret = []byte(cfg.JWTRefreshSecret)
	AccessExpire  int
	RefreshExpire int
)

func generateToken(userID uint, expire int, secret []byte) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateAccessToken(userID uint) (string, error) {
	AccessExpire, err := strconv.Atoi(cfg.JWTAccessExpire)
	if err != nil {
		log.Fatalf("Invalid JWT_ACCESS_EXPIRE: %v", err)
	}

	return generateToken(userID, AccessExpire, AccessSecret)
}

func GenerateRefreshToken(userID uint) (string, error) {
	RefreshExpire, err := strconv.Atoi(cfg.JWTRefreshExpire)
	if err != nil {
		log.Fatalf("Invalid JWT_REFRESH_EXPIRE: %v", err)
	}

	return generateToken(userID, RefreshExpire, RefreshSecret)
}
