package auth

import (
	"github.com/azka-zaydan/synapsis-test/internal/domain/auth/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/auth/service"
	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	AuthSvc service.AuthService
}

func (h *AuthHandler) Router(r fiber.Router) {
	auth := r.Group("/auth")

	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
}

func ProvideAuthHandler(svc service.AuthService) AuthHandler {
	return AuthHandler{
		AuthSvc: svc,
	}
}

// Register registers a new user.
// @Summary registers a new user.
// @Description This endpoint registers a new user.
// @Tags v1/auth
// @Param info body dto.RegisterDto true "user info."
// @Produce json
// @Success 201 {object} response.Base{data=dto.JWTResponse}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterDto
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[RegisterHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}
	token, err := h.AuthSvc.Register(c.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("[RegisterHandler] Failed From Auth Service")
		return response.WithError(c, err)
	}
	return response.WithJSON(c, fiber.StatusOK, token)
}

// Login logs in a user.
// @Summary logs in a user.
// @Description This endpoint logs in a user.
// @Tags v1/auth
// @Param info body dto.RegisterDto true "user info."
// @Produce json
// @Success 201 {object} response.Base{data=dto.JWTResponse}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.RegisterDto
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[LoginHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}
	token, err := h.AuthSvc.Login(c.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("[LoginHandler] Failed From Auth Service")
		return response.WithError(c, err)
	}
	return response.WithJSON(c, fiber.StatusOK, token)
}
