package auth

import (
	"github.com/azka-zaydan/synapsis-test/internal/domain/auth/service"
	"github.com/azka-zaydan/synapsis-test/transport/http/response"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthSvc service.AuthService
}

func (h *AuthHandler) Router(r fiber.Router) {
	auth := r.Group("/auth")

	auth.Post("/register", h.Register)
	auth.Post("/login", h.Register)
}

func ProvideAuthHandler(svc service.AuthService) AuthHandler {
	return AuthHandler{
		AuthSvc: svc,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	return response.WithJSON(c, fiber.StatusOK, "OK")
}
