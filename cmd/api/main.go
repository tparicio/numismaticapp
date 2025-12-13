package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/antonioparicio/numismaticapp/internal/api"
	"github.com/antonioparicio/numismaticapp/internal/application"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/gemini"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/image"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/numista"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 0. Logging Setup
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	// 1. Config
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		slog.Error("GEMINI_API_KEY is not set")
		os.Exit(1)
	}

	// 2. Database
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Initialize Infrastructure
	// Initialize Infrastructure
	coinRepo := infrastructure.NewPostgresCoinRepository(dbPool)
	groupRepo := infrastructure.NewPostgresGroupRepository(dbPool)

	geminiModel := os.Getenv("GEMINI_MODEL")
	geminiClient, err := gemini.NewGeminiService(ctx, os.Getenv("GEMINI_API_KEY"), geminiModel)
	if err != nil {
		slog.Error("Failed to create Gemini client", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := geminiClient.Close(); err != nil {
			slog.Error("Failed to close Gemini client", "error", err)
		}
	}()

	imageService := image.NewVipsImageService()
	storageService := storage.NewLocalFileStorage("storage")

	rembgURL := os.Getenv("REMBG_URL")
	if rembgURL == "" {
		rembgURL = "http://rembg:5000/api/remove" // Default for docker-compose
	}
	rembgClient := image.NewRembgClient(rembgURL)

	// Run Migrations
	migrator := infrastructure.NewMigrationService(dbPool)
	if err := migrator.RunMigrations(ctx); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize Numista Client
	numistaKey := os.Getenv("NUMISTA_API_KEY")
	numistaClient := numista.NewClient(numistaKey)

	// Initialize Application Services
	coinService := application.NewCoinService(coinRepo, groupRepo, imageService, geminiClient, storageService, rembgClient, numistaClient)

	// 5. API
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20MB limit for images
	})
	coinHandler := api.NewCoinHandler(coinService)
	healthHandler := api.NewHealthHandler(dbPool)
	api.SetupRouter(app, coinHandler, healthHandler)

	// 6. Start
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("Server starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
