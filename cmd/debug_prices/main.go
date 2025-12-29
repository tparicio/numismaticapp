package main

import (
	"context"
	"fmt"
	"log"

	"github.com/antonioparicio/numismaticapp/internal/infrastructure/prices"
)

func main() {
	client := prices.NewCoinGeckoPriceClient()
	gold, silver, err := client.GetMetalPrices(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Gold: %.2f EUR/g\nSilver: %.2f EUR/g\n", gold, silver)
}
