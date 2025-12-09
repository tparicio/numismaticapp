package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Coin represents a numismatic coin in the collection.
type Coin struct {
	ID                  uuid.UUID       `json:"id"`
	Country             string          `json:"country"`
	Year                int             `json:"year"`
	FaceValue           string          `json:"face_value"`
	Currency            string          `json:"currency"`
	Material            string          `json:"material"`
	Description         string          `json:"description"`
	KMCode              string          `json:"km_code"`
	MinValue            float64         `json:"min_value"`
	MaxValue            float64         `json:"max_value"`
	Grade               string          `json:"grade"`
	SampleImageURLFront string          `json:"sample_image_url_front"`
	SampleImageURLBack  string          `json:"sample_image_url_back"`
	Notes               string          `json:"notes"`
	GeminiDetails       map[string]any  `json:"gemini_details"` // Raw JSON from Gemini
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

// CoinRepository defines the interface for persisting coins.
type CoinRepository interface {
	Save(ctx context.Context, coin *Coin) error
	GetByID(ctx context.Context, id uuid.UUID) (*Coin, error)
	List(ctx context.Context, limit, offset int) ([]*Coin, error)
	Count(ctx context.Context) (int64, error)
}

// ImageService defines the interface for image processing operations.
type ImageService interface {
	// ProcessCoinImages takes raw front and back images, crops them to circle,
	// and applies the given rotation angle.
	ProcessCoinImages(frontPath, backPath string, rotationAngle float64) (processedFrontPath, processedBackPath string, err error)
	// CropToCircle detects the coin and crops the image to a circle.
	CropToCircle(imagePath string) (string, error)
	// Rotate rotates the image by the given angle.
	Rotate(imagePath string, angle float64) (string, error)
}

// AIService defines the interface for AI analysis.
type AIService interface {
	// AnalyzeCoin analyzes the front and back images of a coin and returns metadata.
	AnalyzeCoin(ctx context.Context, frontImagePath, backImagePath string) (*CoinAnalysisResult, error)
}

// CoinAnalysisResult contains the data extracted by the AI.
type CoinAnalysisResult struct {
	Country                 string         `json:"country"`
	Year                    int            `json:"year"`
	FaceValue               string         `json:"face_value"`
	Currency                string         `json:"currency"`
	Material                string         `json:"material"`
	Description             string         `json:"description"`
	KMCode                  string         `json:"km_code"`
	MinValue                float64        `json:"min_value"`
	MaxValue                float64        `json:"max_value"`
	Grade                   string         `json:"grade"`
	Notes                   string         `json:"notes"`
	VerticalCorrectionAngle float64        `json:"vertical_correction_angle"`
	RawDetails              map[string]any `json:"raw_details"`
}
