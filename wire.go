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
	cartRepo "github.com/azka-zaydan/synapsis-test/internal/domain/cart/repository"
	cartSvc "github.com/azka-zaydan/synapsis-test/internal/domain/cart/service"
	orderRepo "github.com/azka-zaydan/synapsis-test/internal/domain/order/repository"
	orderSvc "github.com/azka-zaydan/synapsis-test/internal/domain/order/service"
	paymentRepo "github.com/azka-zaydan/synapsis-test/internal/domain/payment/repository"
	paymentSvc "github.com/azka-zaydan/synapsis-test/internal/domain/payment/service"
	productRepo "github.com/azka-zaydan/synapsis-test/internal/domain/product/repository"
	productService "github.com/azka-zaydan/synapsis-test/internal/domain/product/service"
	userRepo "github.com/azka-zaydan/synapsis-test/internal/domain/user/repository"
	userSvc "github.com/azka-zaydan/synapsis-test/internal/domain/user/service"
	authHandler "github.com/azka-zaydan/synapsis-test/internal/handlers/auth"
	cartHandler "github.com/azka-zaydan/synapsis-test/internal/handlers/cart"
	paymentHandler "github.com/azka-zaydan/synapsis-test/internal/handlers/payment"
	productHandler "github.com/azka-zaydan/synapsis-test/internal/handlers/product"

	"github.com/azka-zaydan/synapsis-test/transport/http"
	"github.com/azka-zaydan/synapsis-test/transport/http/middleware"
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

var authMiddleware = wire.NewSet(
	middleware.ProvideAuthentication,
)

var domainUser = wire.NewSet(
	userRepo.ProvideUserRepositoryMySQL,
	wire.Bind(new(userRepo.UserRepository), new(*userRepo.UserRepositoryMySQL)),
	userSvc.ProvideUserServiceImpl,
	wire.Bind(new(userSvc.UserService), new(*userSvc.UserServiceImpl)),
)

var domainAuth = wire.NewSet(
	authService.ProvideAuthServiceImpl,
	wire.Bind(new(authService.AuthService), new(*authService.AuthServiceImpl)),
)

var domainProduct = wire.NewSet(
	productRepo.ProvideProductRepositoryMySQL,
	wire.Bind(new(productRepo.ProductRepository), new(*productRepo.ProductRepositoryMySQL)),
	productService.ProvideProductServiceImpl,
	wire.Bind(new(productService.ProductService), new(*productService.ProductServiceImpl)),
)

var domainCart = wire.NewSet(
	cartRepo.ProvideCartRepositoryMySQL,
	wire.Bind(new(cartRepo.CartRepository), new(*cartRepo.CartRepositoryMySQL)),
	cartSvc.ProvideCartServiceImpl,
	wire.Bind(new(cartSvc.CartService), new(*cartSvc.CartServiceImpl)),
)

var domainPayment = wire.NewSet(
	paymentRepo.ProvidePaymentRepositoryMySQL,
	wire.Bind(new(paymentRepo.PaymentRepository), new(*paymentRepo.PaymentRepositoryMySQL)),
	paymentSvc.ProvidePaymentServiceImpl,
	wire.Bind(new(paymentSvc.PaymentService), new(*paymentSvc.PaymentServiceImpl)),
)

var domainOrder = wire.NewSet(
	orderRepo.ProvideOrderRepositoryMySQL,
	wire.Bind(new(orderRepo.OrderRepository), new(*orderRepo.OrderRepositoryMySQL)),
	orderSvc.ProvideOrderServiceImpl,
	wire.Bind(new(orderSvc.OrderService), new(*orderSvc.OrderServiceImpl)),
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainAuth, domainUser, domainProduct, domainCart, domainPayment, domainOrder,
)

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "*"),
	router.ProvideRouter,
	authHandler.ProvideAuthHandler,
	productHandler.ProvideProductHandler,
	cartHandler.ProvideCartHandler,
	paymentHandler.ProvidePaymentHandler,
)

// Wiring for everything.
func InitializeService() *http.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// middleware
		authMiddleware,
		// domains
		domains,
		// routing
		routing,
		// selected transport layer
		http.ProvideHTTP)
	return &http.HTTP{}
}
