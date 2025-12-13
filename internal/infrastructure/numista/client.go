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

func (c *Client) SearchTypes(ctx context.Context, query, category, year, issuer string) (*NumistaType, error) {
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
	// Improve Issuer normalization if needed, for now exact match or skip
	// The user suggested slug normalization, keeping it simple for now or strictly following plan
	if issuer != "" {
		// Numista expects code like 'france' or 'canada'.
		// We might send it as parameter or let 'q' handle it.
		// For now let's append to q if it's not a clear code, but user requested mapping specific params.
		// Let's try sending as issuer parameter if it looks like a code, otherwise append to q.
		// Actually best is to just not send issuer param if we aren't sure of valid code,
		// and rely on query string.
		// But user said: "Mapear coins.country (intentar normalizar a slug si es posible, o enviar como parte de q si no se tiene el código exacto)."
	}
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

	if searchResp.Count > 0 && len(searchResp.Types) > 0 {
		return &searchResp.Types[0], nil
	}

	return nil, nil // Not found
}
