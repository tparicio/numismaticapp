package image

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/antonioparicio/numismaticapp/internal/domain"
)

type RembgClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewRembgClient(baseURL string) *RembgClient {
	return &RembgClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *RembgClient) RemoveBackground(ctx context.Context, image []byte) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(image); err != nil {
		return nil, fmt.Errorf("failed to write image to form: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to rembg: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rembg service returned status: %d", resp.StatusCode)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return result, nil
}

var _ domain.BackgroundRemover = (*RembgClient)(nil)
