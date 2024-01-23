package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateAllTokens(username string) (string, string, error) {

	secret := []byte(os.Getenv("SECRET"))
	accessClaims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	}

	accessjwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessjwtToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}
	refreshjwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshjwtToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
