package payment

import (
	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/service"
	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/transport/http/middleware"
	"github.com/azka-zaydan/synapsis-test/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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
// @Tags v1/payment
// @Param Authorization header string true "Bearer Token"
// @Param payRequest body dto.PayRequest true "orderID to pay"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/payment/pay/ [post]
func (h *PaymentHandler) Pay(c *fiber.Ctx) error {
	var req dto.PayRequest
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[PayHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}

	res, err := h.PaymentSvc.Pay(c.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("[PayHandler] Failed ListItems")
		return response.WithError(c, err)
	}

	return response.WithJSON(c, fiber.StatusOK, res)
}
