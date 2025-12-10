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
	// NOTE: The 'handler' for groups is not defined in the function signature.
	// This change assumes a 'handler' variable (e.g., groupHandler) is available
	// in the scope where this function is called, or that the function signature
	// and imports will be updated separately to include a GroupHandler.
	// For the purpose of this instruction, only the routes are added as specified.
	api.Get("/groups", coinHandler.ListGroups)
	api.Post("/groups", coinHandler.CreateGroup)
	api.Put("/groups/:id", coinHandler.UpdateGroup)
	api.Delete("/groups/:id", coinHandler.DeleteGroup)
	v1.Get("/coins", coinHandler.ListCoins)
	v1.Get("/coins/:id", coinHandler.GetCoin)
	v1.Put("/coins/:id", coinHandler.UpdateCoin)
	v1.Delete("/coins/:id", coinHandler.DeleteCoin)
	v1.Get("/dashboard", coinHandler.GetDashboardStats)

	// SPA Fallback: Serve index.html for any other route not handled above
	// This ensures that refreshing pages like /coins/add works
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./web/dist/index.html")
	})
}
