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

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.1) // Lower temperature for more deterministic results

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
	Actúa como un experto numismático y analista de imagen. Tu tarea es extraer datos técnicos, calcular la corrección de rotación de una moneda y buscar referencias visuales comparativas.

	INSTRUCCIONES DE VISIÓN Y ORIENTACIÓN (CRÍTICO):
	1. Ignora el cartón, la cápsula de plástico o el fondo. Céntrate solo en el disco metálico.
	2. Identifica la "parte superior" lógica del diseño (cabeza, escudo, texto).
	3. Imagina un reloj superpuesto. Si la parte superior está a las 3 en punto (derecha) -> Ángulo -90. A las 9 (izquierda) -> 90. A las 6 (abajo) -> 180.
	4. Calcula el ángulo exacto para dejar la moneda vertical.

	INSTRUCCIONES DE BÚSQUEDA DE REFERENCIA:
	1. Identifica el Código KM (Krause) de la moneda.
	2. Busca el número identificador de Numista para esta moneda (Numista Number).

	INSTRUCCIONES DE SALIDA:
	Genera UNICAMENTE un objeto JSON válido. Sin markdown.

	Estructura JSON requerida:
	{
		"_debug_orientation_front": "Descripción de la orientación visual del anverso",
		"vertical_correction_angle_front": 0.0,
		"_debug_orientation_back": "Descripción de la orientación visual del reverso",
		"vertical_correction_angle_back": 0.0,
		"name": "Título descriptivo (ej: 25 Pesetas - Mundial 82)",
		"country": "País",
		"year": 1980,
		"face_value": "Valor facial",
		"currency": "Unidad monetaria",
		"material": "Material",
		"description": "Descripción visual",
		"km_code": "Código KM#",
		"numista_number": 0,
		"min_value": 0.0,
		"max_value": 0.0,
		"grade": "Estado estimado",
		"reference_source_name": "Nombre de la fuente (ej: Numista, uCoin)",
		"notes": "Notas",
		"weight_g": 0.0,
		"diameter_mm": 0.0,
		"thickness_mm": 0.0,
		"edge": "Canto",
		"shape": "Forma",
		"mint": "Ceca",
		"mintage": 0
	}
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
