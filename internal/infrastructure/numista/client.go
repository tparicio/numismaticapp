package numista

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: "https://api.numista.com/v3",
		APIKey:  apiKey,
		HTTP:    http.DefaultClient,
	}
}

type TypeSearchResponse struct {
	Count int           `json:"count"`
	Types []NumistaType `json:"types"`
}

type NumistaType struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	Category         string `json:"category"`
	MinYear          int    `json:"min_year"`
	MaxYear          int    `json:"max_year"`
	ObverseThumbnail string `json:"obverse_thumbnail"`
	ReverseThumbnail string `json:"reverse_thumbnail"`
	Issuer           Issuer `json:"issuer"`
}

type Issuer struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (c *Client) SearchTypes(ctx context.Context, query, category, year, issuer string, count int) (*TypeSearchResponse, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("numista API key is not set")
	}

	u, err := url.Parse(c.BaseURL + "/types")
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	q := u.Query()
	q.Set("q", query)
	q.Set("category", category)
	if year != "" {
		q.Set("year", year)
	}
	if count > 0 {
		q.Set("count", fmt.Sprintf("%d", count))
	} else {
		q.Set("count", "10")
	}

	// Improve Issuer normalization if needed.
	// For now we rely on 'q' or simple params as implemented before.

	q.Set("lang", "es") // Default as requested

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Numista-API-Key", c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("numista api error: %s - %s", resp.Status, string(body))
	}

	var searchResp TypeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResp, nil
}

func (c *Client) GetType(ctx context.Context, id int) (map[string]any, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("numista API key is not set")
	}

	u, err := url.Parse(fmt.Sprintf("%s/types/%d", c.BaseURL, id))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	q := u.Query()
	q.Set("lang", "es")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Numista-API-Key", c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("numista api error: %s - %s", resp.Status, string(body))
	}

	var typeDetails map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&typeDetails); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return typeDetails, nil
}
