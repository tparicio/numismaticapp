package main

import (
	"context"
	"log"
	"os"

	"github.com/antonioparicio/numismaticapp/internal/api"
	"github.com/antonioparicio/numismaticapp/internal/application"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/gemini"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/image"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// 1. Config
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	// 2. Database
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Initialize Infrastructure
	// Initialize Infrastructure
	coinRepo := infrastructure.NewPostgresCoinRepository(dbPool)
	groupRepo := infrastructure.NewPostgresGroupRepository(dbPool)

	geminiModel := os.Getenv("GEMINI_MODEL")
	geminiClient, err := gemini.NewGeminiService(ctx, os.Getenv("GEMINI_API_KEY"), geminiModel)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	defer geminiClient.Close()

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
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Application Services
	coinService := application.NewCoinService(coinRepo, groupRepo, imageService, geminiClient, storageService, rembgClient)

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
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
