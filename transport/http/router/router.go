package router

import (
	"github.com/azka-zaydan/synapsis-test/internal/handlers/auth"
	"github.com/gofiber/fiber/v2"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	AuthHanlder auth.AuthHandler
}

// Router is the router struct containing handlers.
type Router struct {
	DomainHandlers DomainHandlers
}

// ProvideRouter is the provider function for this router.
func ProvideRouter(domainHandlers DomainHandlers) Router {
	return Router{
		DomainHandlers: domainHandlers,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(app *fiber.App) {
	app.Route("/v1", func(router fiber.Router) {
		r.DomainHandlers.AuthHanlder.Router(router)
	})
}
