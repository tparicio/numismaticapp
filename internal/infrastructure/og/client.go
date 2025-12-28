package og

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Metadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

var (
	ogTitleRegex       = regexp.MustCompile(`<meta\s+property=["']og:title["']\s+content=["'](.*?)["']\s*/?>`)
	ogDescriptionRegex = regexp.MustCompile(`<meta\s+property=["']og:description["']\s+content=["'](.*?)["']\s*/?>`)
	ogImageRegex       = regexp.MustCompile(`<meta\s+property=["']og:image["']\s+content=["'](.*?)["']\s*/?>`)
	titleRegex         = regexp.MustCompile(`<title>(.*?)</title>`)
	descriptionRegex   = regexp.MustCompile(`<meta\s+name=["']description["']\s+content=["'](.*?)["']\s*/?>`)
)

func FetchMetadata(ctx context.Context, url string) (*Metadata, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// Set User-Agent to avoid being blocked by some sites
	req.Header.Set("User-Agent", "NumismaticApp/1.0")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read body (limit to 1MB to avoid memory issues with large pages)
	// We only need the head mostly.
	buf := make([]byte, 1024*1024)
	n, err := resp.Body.Read(buf)
	if err != nil && n == 0 {
		// Allow EOF if we read something
		if err.Error() != "EOF" {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
	}
	body := string(buf[:n])

	meta := &Metadata{}

	// Extract Title
	if match := ogTitleRegex.FindStringSubmatch(body); len(match) > 1 {
		meta.Title = match[1]
	} else if match := titleRegex.FindStringSubmatch(body); len(match) > 1 {
		meta.Title = match[1]
	}

	// Extract Description
	if match := ogDescriptionRegex.FindStringSubmatch(body); len(match) > 1 {
		meta.Description = match[1]
	} else if match := descriptionRegex.FindStringSubmatch(body); len(match) > 1 {
		meta.Description = match[1]
	}

	// Extract Image
	if match := ogImageRegex.FindStringSubmatch(body); len(match) > 1 {
		meta.Image = match[1]
	}

	return meta, nil
}
