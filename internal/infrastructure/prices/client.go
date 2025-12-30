package prices

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type PriceClient interface {
	GetMetalPrices(ctx context.Context) (float64, float64, error) // Gold, Silver (EUR per gram)
}

type CoinGeckoResponse struct {
	PaxGold       map[string]float64 `json:"pax-gold"`
	KinesisSilver map[string]float64 `json:"kinesis-silver"`
}

type CoinGeckoPriceClient struct {
	client        *http.Client
	cache         map[string]float64
	cacheTime     time.Time
	mu            sync.RWMutex
	cacheDuration time.Duration
}

func NewCoinGeckoPriceClient() *CoinGeckoPriceClient {
	return &CoinGeckoPriceClient{
		client:        &http.Client{Timeout: 10 * time.Second},
		cache:         make(map[string]float64),
		cacheDuration: 10 * time.Minute, // Cache for 10 minutes to respect rate limits
	}
}

// GetMetalPrices returns price per Gram in EUR for Gold and Silver
func (c *CoinGeckoPriceClient) GetMetalPrices(ctx context.Context) (float64, float64, error) {
	c.mu.RLock()
	if time.Since(c.cacheTime) < c.cacheDuration {
		gold := c.cache["gold"]
		silver := c.cache["silver"]
		if gold > 0 && silver > 0 {
			c.mu.RUnlock()
			return gold, silver, nil
		}
	}
	c.mu.RUnlock()

	// Fetch Gold (PAX Gold = 1 troy oz Gold)
	goldPriceOz, err := c.fetchPrice(ctx, "pax-gold")
	if err != nil {
		slog.Error("Failed to fetch gold price", "error", err)
		// Return cached values if available, even if expired
		c.mu.RLock()
		gold := c.cache["gold"]
		silver := c.cache["silver"]
		c.mu.RUnlock()
		if gold > 0 && silver > 0 {
			return gold, silver, nil // fallback to old cache
		}
		// If no cache, return fallback values
		return 60.0, 0.75, nil
	}

	// Fetch Silver (Kinesis Silver = 1 oz Silver)
	silverPriceOz, err := c.fetchPrice(ctx, "kinesis-silver")
	if err != nil {
		slog.Error("Failed to fetch silver price", "error", err)
		silverPriceOz = 23.0 // Approximate fallback per oz
	}

	// Conversion constants
	const gramPerOz = 31.1034768

	goldPriceGram := goldPriceOz / gramPerOz
	silverPriceGram := silverPriceOz / gramPerOz

	c.mu.Lock()
	c.cache["gold"] = goldPriceGram
	c.cache["silver"] = silverPriceGram
	c.cacheTime = time.Now()
	c.mu.Unlock()

	slog.Info("GetMetalPrices", "gold_eur_g", goldPriceGram, "silver_eur_g", silverPriceGram)
	return goldPriceGram, silverPriceGram, nil
}

func (c *CoinGeckoPriceClient) fetchPrice(ctx context.Context, id string) (float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=eur", id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}
	// Add User-Agent as CoinGecko sometimes blocks generic Go-http-client
	req.Header.Set("User-Agent", "NumismaticApp/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	slog.Info("CoinGecko Response", "id", id, "status", resp.StatusCode, "body", string(body))

	// Parse dynamic JSON where key is the id
	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if priceMap, ok := result[id]; ok {
		if price, ok := priceMap["eur"]; ok {
			return price, nil
		}
	}

	return 0, fmt.Errorf("price not found in response")
}
