package product

import (
	"github.com/azka-zaydan/synapsis-test/internal/domain/product/model"
	"github.com/azka-zaydan/synapsis-test/internal/domain/product/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/product/service"
	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/transport/http/middleware"
	"github.com/azka-zaydan/synapsis-test/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type ProductHandler struct {
	auth    *middleware.Authentication
	service service.ProductService
}

func (h *ProductHandler) Router(r fiber.Router) {
	product := r.Group("/product", h.auth.JWTAuth())

	product.Post("/", h.CreateProduct)
	product.Post("/filter", h.GetProductsByFilter)
}

func ProvideProductHandler(auth *middleware.Authentication, svc service.ProductService) ProductHandler {
	return ProductHandler{
		auth:    auth,
		service: svc,
	}
}

// GetProducts gets all products by filter
// @Summary gets all products by filter
// @Description This endpoint gets all products by filter
// @Tags v1/product
// @Param Authorization header string true "Bearer Token"
// @Param Filter body model.Filter true "filter"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/product/filter [post]
func (h *ProductHandler) GetProductsByFilter(c *fiber.Ctx) error {
	var req model.Filter
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[GetProductsByFilterHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}
	err = dto.ValidateAndSetDefaultFilter(&req)
	if err != nil {
		log.Error().Err(err).Msg("[GetProductsByFilterHandler] Failed Validating Body")
		return response.WithError(c, failure.BadRequest(err))
	}
	res, err := h.service.GetProductByFilter(c.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("[GetProductsByFilterHandler] Failed GetProductByFilter")
		return response.WithError(c, err)
	}
	return response.WithMetadata(c, fiber.StatusOK, res.Data, res.Metadata)
}

// CreateProduct creates a new product
// @Summary creates a new product
// @Description This endpoint creates a new product
// @Tags v1/product
// @Param Authorization header string true "Bearer Token"
// @Param createProduct body dto.ProductCreateRequest true "create product body"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/product/ [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.ProductCreateRequest
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Err(err).Msg("[CreateProductHandler] Failed Parsing Body")
		return response.WithError(c, failure.BadRequest(err))
	}
	res, err := h.service.CreateProduct(c.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("[CreateProductHandler] Failed CreateProduct")
		return response.WithError(c, err)
	}
	return response.WithJSON(c, fiber.StatusOK, res)
}
