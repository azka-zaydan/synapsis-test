package jwt

import (
	"time"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	cfg *configs.Config
	key []byte
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJwtService(cfg *configs.Config) *JwtService {
	key := []byte(cfg.JWT.Key)
	return &JwtService{
		cfg: cfg,
		key: key,
	}
}

func (s *JwtService) GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(s.cfg.JWT.ExpiresIn)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.key)
}

func (s *JwtService) ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
