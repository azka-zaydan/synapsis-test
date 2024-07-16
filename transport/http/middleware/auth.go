package middleware

import (
	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Authentication struct {
	cfg *configs.Config
}

func ProvideAuthentication(cfg *configs.Config) *Authentication {
	return &Authentication{
		cfg: cfg,
	}
}

func (m *Authentication) JWTAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(m.cfg.JWT.Key),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
			}
			return nil
		},
	})
}
