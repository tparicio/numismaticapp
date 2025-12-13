package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GeminiService struct {
	client       *genai.Client
	mu           sync.RWMutex
	cachedModels []domain.GeminiModelInfo
	lastCache    time.Time
}

func NewGeminiService(ctx context.Context, apiKey string, _ string) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &GeminiService{
		client: client,
	}, nil
}

func (s *GeminiService) Close() error {
	return s.client.Close()
}

func (s *GeminiService) ListModels(ctx context.Context) ([]domain.GeminiModelInfo, error) {
	s.mu.RLock()
	// Cache for 1 hour
	if len(s.cachedModels) > 0 && time.Since(s.lastCache) < 1*time.Hour {
		defer s.mu.RUnlock()
		return s.cachedModels, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check
	if len(s.cachedModels) > 0 && time.Since(s.lastCache) < 1*time.Hour {
		return s.cachedModels, nil
	}

	iter := s.client.ListModels(ctx)
	var models []domain.GeminiModelInfo
	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list models: %w", err)
		}

		// Filter for Gemini models that generate content
		if strings.Contains(m.Name, "gemini") && (strings.Contains(m.SupportedGenerationMethods[0], "generateContent") || len(m.SupportedGenerationMethods) > 0) {
			isContent := false
			for _, method := range m.SupportedGenerationMethods {
				if method == "generateContent" {
					isContent = true
					break
				}
			}
			if isContent {
				cleanName := strings.TrimPrefix(m.Name, "models/")
				models = append(models, domain.GeminiModelInfo{
					Name:        cleanName,
					Description: m.Description,
				})
			}
		}
	}

	s.cachedModels = models
	s.lastCache = time.Now()

	return models, nil
}

func (s *GeminiService) AnalyzeCoin(ctx context.Context, frontImagePath, backImagePath string, modelName string, temperature float32, lang string) (*domain.CoinAnalysisResult, error) {
	frontData, err := os.ReadFile(frontImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read front image: %w", err)
	}

	backData, err := os.ReadFile(backImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read back image: %w", err)
	}

	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}

	// Ensure model has temperature set
	model := s.client.GenerativeModel(modelName)
	model.SetTemperature(temperature)

	promptGen := NewPromptGenerator()
	prompt := promptGen.GetPrompt(lang)

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.ImageData("jpeg", frontData),
		genai.ImageData("jpeg", backData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no content returned from gemini")
	}

	// Extract text from response
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Clean up markdown code blocks if present
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	var result domain.CoinAnalysisResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w. Response: %s", err, responseText)
	}

	// Store raw details for debugging/extra info
	var rawDetails map[string]any
	_ = json.Unmarshal([]byte(responseText), &rawDetails)
	result.RawDetails = rawDetails

	return &result, nil
}
