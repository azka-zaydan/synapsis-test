package jwt

import (
	"time"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	jwtV5 "github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	cfg *configs.Config
	key []byte
}

type Claims struct {
	UserID   string `json:"UserID"`
	Username string `json:"username"`
	jwtV5.RegisteredClaims
}

func NewJwtService(cfg *configs.Config) *JwtService {
	key := []byte(cfg.JWT.Key)
	return &JwtService{
		cfg: cfg,
		key: key,
	}
}

func (s *JwtService) GenerateJWT(username, id string) (string, error) {
	expirationTime := time.Now().Add(s.cfg.JWT.ExpiresIn)
	claims := &Claims{
		Username: username,
		UserID:   id,
		RegisteredClaims: jwtV5.RegisteredClaims{
			ExpiresAt: jwtV5.NewNumericDate(expirationTime),
		},
	}
	token := jwtV5.NewWithClaims(jwtV5.SigningMethodHS256, claims)
	return token.SignedString(s.key)
}

func (s *JwtService) ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwtV5.ParseWithClaims(tokenStr, claims, func(token *jwtV5.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		if err == jwtV5.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}

func GetClaims(c *fiber.Ctx) jwt.MapClaims {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}
