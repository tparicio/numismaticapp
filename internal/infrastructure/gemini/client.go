package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService(ctx context.Context, apiKey string) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.4) // Lower temperature for more deterministic results

	return &GeminiService{
		client: client,
		model:  model,
	}, nil
}

func (s *GeminiService) Close() error {
	return s.client.Close()
}

func (s *GeminiService) AnalyzeCoin(ctx context.Context, frontImagePath, backImagePath string) (*domain.CoinAnalysisResult, error) {
	frontData, err := os.ReadFile(frontImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read front image: %w", err)
	}

	backData, err := os.ReadFile(backImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read back image: %w", err)
	}

	prompt := `
	Analiza estas imágenes de una moneda (anverso y reverso). 
	Devuelve UNICAMENTE un objeto JSON válido (sin markdown, sin texto adicional) con la siguiente estructura:
	{
		"country": "País de origen",
		"year": 1999,
		"face_value": "Valor facial (ej: 1 Euro)",
		"currency": "Moneda (ej: Euro)",
		"material": "Material (ej: Oro, Plata, Cobre)",
		"description": "Descripción visual detallada",
		"km_code": "Código KM# si es identificable",
		"min_value": 10.5,
		"max_value": 20.0,
		"grade": "Estado de conservación estimado (ej: EBC, MBC)",
		"notes": "Notas técnicas o de conservación",
		"vertical_correction_angle": 0.0,
		"weight_g": 0.0,
		"diameter_mm": 0.0,
		"thickness_mm": 0.0,
		"edge": "Descripción del canto (estriado, liso, leyenda...)",
		"shape": "Forma (redonda, cuadrada...)",
		"mint": "Ceca o marca de ceca",
		"mintage": 0
	}
	
	IMPORTANTE: 
	1. El campo 'vertical_correction_angle' debe ser el ángulo de rotación en grados (positivo o negativo) necesario para que el anverso de la moneda quede perfectamente vertical.
	2. Los campos numéricos (weight_g, diameter_mm, thickness_mm, mintage) deben ser estimaciones si no se pueden determinar con exactitud, o 0 si son totalmente desconocidos.
	`

	resp, err := s.model.GenerateContent(ctx,
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
