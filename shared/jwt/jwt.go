package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenData struct {
	Username string
}

func NewToken(username string) *jwt.Token {
	claims := jwt.MapClaims{
		"name": username,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}

func ParseToken(token *jwt.Token) *TokenData {
	claims := token.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return &TokenData{
		Username: name,
	}
}
