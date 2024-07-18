package cart

import (
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/service"
	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/shared/jwt"
	"github.com/azka-zaydan/synapsis-test/transport/http/middleware"
	"github.com/azka-zaydan/synapsis-test/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type CartHandler struct {
	auth    *middleware.Authentication
	CartSvc service.CartService
}

func (h *CartHandler) Router(r fiber.Router) {
	cart := r.Group("/cart", h.auth.JWTAuth())

	cart.Post("/add-items", h.AddItems)
	cart.Get("/list-items", h.ListItems)
	cart.Post("/remove-items", h.DeleteItems)

	cart.Post("/checkout", h.Checkout)
}

func ProvideCartHandler(svc service.CartService, auth *middleware.Authentication) CartHandler {
	return CartHandler{
		CartSvc: svc,
		auth:    auth,
	}
}

// ListItems gets all cart items
// @Summary gets all cart items
// @Description This endpoint gets all cart items
// @Tags v1/cart
// @Param Authorization header string true "Bearer Token"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/list-items [get]
func (h *CartHandler) ListItems(c *fiber.Ctx) error {
	userIDStr := jwt.GetClaims(c)["userID"].(string)

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		log.Error().Err(err).Msg("[ListItemsHandler] Failed Converting into UUID")
		return response.WithError(c, failure.BadRequest(err))
	}

	res, err := h.CartSvc.ListItems(c.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("[ListItemsHandler] Failed ListItems")
		return response.WithError(c, err)
	}

	return response.WithJSON(c, fiber.StatusOK, res)
}

// ListItems gets all cart items
// @Summary gets all cart items
// @Description This endpoint gets all cart items
// @Tags v1/cart
// @Param Authorization header string true "Bearer Token"
// @Param addItemsRequest body dto.AddItemsRequest true "items to be added"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/add-items [post]
func (h *CartHandler) AddItems(c *fiber.Ctx) error {
	userIDStr := jwt.GetClaims(c)["userID"].(string)

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		log.Error().Err(err).Msg("[AddItemsHandler] Failed Converting into UUID")
		return response.WithError(c, failure.BadRequest(err))
	}
	var req dto.AddItemsRequest
	err = c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[AddItemsHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}

	res, err := h.CartSvc.AddItems(c.Context(), req, userID)
	if err != nil {
		log.Error().Err(err).Msg("[AddItemsHandler] Failed ListItems")
		return response.WithError(c, err)
	}

	return response.WithJSON(c, fiber.StatusOK, res)
}

// DeleteItems remove cart items based on request
// @Summary remove cart items based on request
// @Description This endpoint removes cart items based on request
// @Tags v1/cart
// @Param Authorization header string true "Bearer Token"
// @Param addItemsRequest body dto.DeleteItemsRequest true "items to be deleted"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/remove-items [post]
func (h *CartHandler) DeleteItems(c *fiber.Ctx) error {
	userIDStr := jwt.GetClaims(c)["userID"].(string)

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItemsHandler] Failed Converting into UUID")
		return response.WithError(c, failure.BadRequest(err))
	}
	var req dto.DeleteItemsRequest
	err = c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItemsHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}

	res, err := h.CartSvc.DeleteItems(c.Context(), req, userID)
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItemsHandler] Failed ListItems")
		return response.WithError(c, err)
	}

	return response.WithJSON(c, fiber.StatusOK, res)
}

// Checkout checks out items based on request
// @Summary checks out items based on request
// @Description This endpoint checks out items based on request
// @Tags v1/cart
// @Param Authorization header string true "Bearer Token"
// @Param addItemsRequest body dto.CheckoutRequest true "items to be deleted"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/checkout [post]
func (h *CartHandler) Checkout(c *fiber.Ctx) error {
	userIDStr := jwt.GetClaims(c)["userID"].(string)

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		log.Error().Err(err).Msg("[CheckoutHandler] Failed Converting into UUID")
		return response.WithError(c, failure.BadRequest(err))
	}
	var req dto.CheckoutRequest
	err = c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[CheckoutHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}

	res, err := h.CartSvc.Checkout(c.Context(), req, userID)
	if err != nil {
		log.Error().Err(err).Msg("[CheckoutHandler] Failed ListItems")
		return response.WithError(c, err)
	}

	return response.WithJSON(c, fiber.StatusOK, res)
}
