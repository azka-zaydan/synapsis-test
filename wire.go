//go:build wireinject
// +build wireinject

package main

import (
	"github.com/azka-zaydan/synapsis-test/configs"
	// "github.com/azka-zaydan/synapsis-test/event"
	// fooBarBazEvent "github.com/azka-zaydan/synapsis-test/event/domain/foobarbaz"
	// "github.com/azka-zaydan/synapsis-test/event/producer"
	"github.com/azka-zaydan/synapsis-test/infras"
	// "github.com/azka-zaydan/synapsis-test/internal/domain/foobarbaz"
	authService "github.com/azka-zaydan/synapsis-test/internal/domain/auth/service"
	authHandler "github.com/azka-zaydan/synapsis-test/internal/handlers/auth"
	"github.com/azka-zaydan/synapsis-test/transport/http"

	// "github.com/azka-zaydan/synapsis-test/transport/http/middleware"
	"github.com/azka-zaydan/synapsis-test/transport/http/router"
	"github.com/google/wire"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistences = wire.NewSet(
	infras.ProvideMySQLConn,
	infras.RedisNewClient,
)

var domainAuth = wire.NewSet(
	authService.ProvideAuthServiceImpl,
	wire.Bind(new(authService.AuthService), new(*authService.AuthServiceImpl)),
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainAuth,
)

// var authMiddleware = wire.NewSet(
// 	middleware.ProvideAuthentication,
// )

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	router.ProvideRouter,
	authHandler.ProvideAuthHandler,
)

// Wiring for everything.
func InitializeService() *http.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// middleware
		// authMiddleware,
		// domains
		domains,
		// routing
		routing,
		// selected transport layer
		http.ProvideHTTP)
	return &http.HTTP{}
}
