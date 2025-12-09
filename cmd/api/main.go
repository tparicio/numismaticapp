package main

import (
	"context"
	"log"
	"os"

	"github.com/antonioparicio/numismaticapp/internal/api"
	"github.com/antonioparicio/numismaticapp/internal/application"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
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
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 3. Infrastructure
	coinRepo := infrastructure.NewPostgresCoinRepository(dbPool)
	geminiService, err := infrastructure.NewGeminiService(context.Background(), geminiKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)
	}
	imageService := infrastructure.NewVipsImageService()
	storageService := infrastructure.NewLocalFileStorage("./storage")

	// 4. Application
	coinService := application.NewCoinService(coinRepo, imageService, geminiService, storageService)

	// 5. API
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20MB limit for images
	})
	coinHandler := api.NewCoinHandler(coinService)
	api.SetupRouter(app, coinHandler)

	// 6. Start
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
