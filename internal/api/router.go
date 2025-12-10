package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter(app *fiber.App, coinHandler *CoinHandler, healthHandler *HealthHandler) {
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Static files (Images)
	// Assuming storage is at ./storage relative to execution
	app.Static("/storage", "./storage")

	// Serve Frontend Static Files
	app.Static("/", "./web/dist")

	// API Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health Check
	v1.Get("/health", healthHandler.HealthCheck)

	v1.Post("/coins", coinHandler.AddCoin)
	v1.Get("/coins", coinHandler.ListCoins)
	v1.Get("/coins/:id", coinHandler.GetCoin)
	v1.Get("/dashboard", coinHandler.GetDashboardStats)

	// SPA Fallback: Serve index.html for any other route not handled above
	// This ensures that refreshing pages like /coins/add works
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./web/dist/index.html")
	})
}
