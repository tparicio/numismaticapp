package gemini

import "strings"

type PromptGenerator struct{}

func NewPromptGenerator() *PromptGenerator {
	return &PromptGenerator{}
}

func (p *PromptGenerator) GetPrompt(lang string) string {
	lang = strings.ToLower(lang)
	if strings.HasPrefix(lang, "en") {
		return p.getEnglishPrompt()
	}
	// Default to Spanish
	return p.getSpanishPrompt()
}

func (p *PromptGenerator) getSpanishPrompt() string {
	return `
	Actúa como un experto numismático y analista de imagen. Tu tarea es extraer datos básicos y descriptivos de una moneda.

	INSTRUCCIONES DE VISIÓN (CRÍTICO):
	1. Ignora el cartón, la cápsula de plástico o el fondo. Céntrate solo en el disco metálico.
	2. Intenta leer el texto legible en la moneda para identificar país, año y valor.

	INSTRUCCIONES DE SALIDA:
	Genera UNICAMENTE un objeto JSON válido. Sin markdown.
	
	Estructura JSON requerida:
	{
		"name": "Título descriptivo (ej: 25 Pesetas - Mundial 82)",
		"country": "País",
		"year": 0,
		"face_value": "Valor facial",
		"currency": "Unidad monetaria",
		"material": "Material (ej: Plata .925, Cobre)",
		"description": "Descripción visual detallada del anverso y reverso",
		"km_code": "Código KM (ej: KM# 819)",

		"mint": "Ceca (ej: Madrid, M coronada)",
		"mintage": 0,
		"min_value": 0.0,
		"max_value": 0.0,
		"grade": "Estado estimado (USAR SOLO: PROOF, FDC, SC, EBC, MBC, BC, RC, MC)",
		"notes": "Cualquier nota adicional relevante observada"
	}
	`
}

func (p *PromptGenerator) getEnglishPrompt() string {
	return `
	Act as a numismatic expert and image analyst. Your task is to extract basic and descriptive data from a coin.

	VISION INSTRUCTIONS (CRITICAL):
	1. Ignore the cardboard, plastic capsule, or background. Focus only on the metal disk.
	2. Try to read the legible text on the coin to identify country, year, and value.

	OUTPUT INSTRUCTIONS:
	Generate ONLY a valid JSON object. No markdown.

	Required JSON structure:
	{
		"name": "Descriptive title (e.g. 25 Pesetas - World Cup 82)",
		"country": "Country",
		"year": 0,
		"face_value": "Face Value",
		"currency": "Currency",
		"material": "Material (e.g. Silver .925, Copper)",
		"description": "Detailed visual description of obverse and reverse",
		"km_code": "KM Code (e.g. KM# 819)",

		"mint": "Mint (e.g. Madrid)",
		"mintage": 0,
		"min_value": 0.0,
		"max_value": 0.0,
		"grade": "Estimated Condition (USE ONLY: PROOF, UNC, XF, VF, F, VG, G, AG)",
		"notes": "Any additional relevant notes observed"
	}
	`
}
