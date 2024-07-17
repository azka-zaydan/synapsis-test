package payment

import (
	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/service"
	"github.com/azka-zaydan/synapsis-test/transport/http/middleware"
	"github.com/azka-zaydan/synapsis-test/transport/http/response"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	PaymentSvc service.PaymentService
	auth       *middleware.Authentication
}

func (h *PaymentHandler) Router(r fiber.Router) {
	payment := r.Group("/payment", h.auth.JWTAuth())

	payment.Post("/pay", h.Pay)
}

func ProvidePaymentHandler(svc service.PaymentService, auth *middleware.Authentication) PaymentHandler {
	return PaymentHandler{
		PaymentSvc: svc,
		auth:       auth,
	}
}

// CreateProduct creates a new product
// @Summary creates a new product
// @Description This endpoint creates a new product
// @Tags v1/product
// @Param Authorization header string true "Bearer Token"
// @Param payRequest body dto.PayRequest true "orderID to pay"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/payment/pay/:id [post]
func (h *PaymentHandler) Pay(c *fiber.Ctx) error {
	return response.WithMessage(c, fiber.StatusOK, "OK")
}
