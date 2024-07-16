package product

import (
	"github.com/azka-zaydan/synapsis-test/transport/http/middleware"
	"github.com/azka-zaydan/synapsis-test/transport/http/response"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	auth *middleware.Authentication
}

func (h *ProductHandler) Router(r fiber.Router) {
	product := r.Group("/product", h.auth.JWTAuth())

	product.Post("/", h.GetProducts)
	product.Get("/", h.GetProducts)
	product.Get("/", h.GetProducts)
}

func ProvideProductHandler(auth *middleware.Authentication) ProductHandler {
	return ProductHandler{
		auth: auth,
	}
}

// GetProducts gets all products
// @Summary gets all products
// @Description This endpoint gets all products
// @Tags v1/product
// @Param Authorization header string true "Bearer Token"
// @Produce json
// @Success 201 {object} response.Base{}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/product/ [post]
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	return response.WithMessage(c, fiber.StatusOK, "OK")
}
