package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter(app *fiber.App, coinHandler *CoinHandler) {
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Static files (Images)
	// Assuming storage is at ./storage relative to execution
	app.Static("/storage", "./storage")

	// API Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/coins", coinHandler.AddCoin)
	v1.Get("/coins", coinHandler.ListCoins)
	v1.Get("/coins/:id", coinHandler.GetCoin)
}
