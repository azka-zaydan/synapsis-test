package main

//go:generate go run github.com/swaggo/swag/cmd/swag init
//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/azka-zaydan/synapsis-test/shared/logger"
	"github.com/gofiber/fiber/v2"
)

var config *configs.Config

func main() {
	app := fiber.New()
	// Initialize logger
	logger.InitLogger()

	// Initialize config
	config = configs.Get()

	// Set desired log level
	logger.SetLogLevel(config)

	// Wire everything up
	http := InitializeService()

	// consumers := InitializeEvent()

	// // Start consumers
	// consumers.Start()

	// Run server
	http.SetupAndServe(app)
}
