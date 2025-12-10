package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Coin represents a numismatic coin in the collection.
type Coin struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Mint           string         `json:"mint"`
	Mintage        int64          `json:"mintage"`
	Country        string         `json:"country"`
	Year           int            `json:"year"`
	FaceValue      string         `json:"face_value"`
	Currency       string         `json:"currency"`
	Material       string         `json:"material"`
	Description    string         `json:"description"`
	KMCode         string         `json:"km_code"`
	MinValue       float64        `json:"min_value"`
	MaxValue       float64        `json:"max_value"`
	Grade          string         `json:"grade"`
	TechnicalNotes string         `json:"technical_notes"`
	GeminiDetails  map[string]any `json:"gemini_details"` // Raw JSON from Gemini
	Images         []CoinImage    `json:"images"`
	GroupID        *int           `json:"group_id"`
	PersonalNotes  string         `json:"personal_notes"`
	WeightG        float64        `json:"weight_g"`
	DiameterMM     float64        `json:"diameter_mm"`
	ThicknessMM    float64        `json:"thickness_mm"`
	Edge           string         `json:"edge"`
	Shape          string         `json:"shape"`
	AcquiredAt     *time.Time     `json:"acquired_at"`
	SoldAt         *time.Time     `json:"sold_at"`
	PricePaid      float64        `json:"price_paid"`
	SoldPrice      float64        `json:"sold_price"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Group struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CoinImage struct {
	ID               uuid.UUID `json:"id"`
	CoinID           uuid.UUID `json:"coin_id"`
	ImageType        string    `json:"image_type"` // original, crop, thumbnail, sample
	Side             string    `json:"side"`       // front, back
	Path             string    `json:"path"`
	Extension        string    `json:"extension"`
	Size             int64     `json:"size"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	MimeType         string    `json:"mime_type"`
	OriginalFilename string    `json:"original_filename"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CoinRepository defines the interface for persisting coins.
type CoinRepository interface {
	Save(ctx context.Context, coin *Coin) error
	GetByID(ctx context.Context, id uuid.UUID) (*Coin, error)
	List(ctx context.Context, limit, offset int) ([]*Coin, error)
	Count(ctx context.Context) (int64, error)
}

// GroupRepository defines the interface for persisting groups.
type GroupRepository interface {
	Create(ctx context.Context, name, description string) (*Group, error)
	GetByName(ctx context.Context, name string) (*Group, error)
	List(ctx context.Context) ([]*Group, error)
}

// ImageService defines the interface for image processing operations.
type ImageService interface {
	// ProcessCoinImages takes raw front and back images, crops them to circle.
	ProcessCoinImages(frontPath, backPath string) (processedFrontPath, processedBackPath string, err error)
	// CropToCircle detects the coin and crops the image to a circle.
	CropToCircle(imagePath string) (string, error)
	GetMetadata(imagePath string) (width, height int, size int64, mimeType string, err error)
	// CropToContent crops the image to remove transparent borders and centers it.
	CropToContent(image []byte) ([]byte, error)
	// GenerateThumbnail creates a smaller version of the image preserving aspect ratio and transparency.
	GenerateThumbnail(imagePath string, width int) (string, error)
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
	WeightG                 float64        `json:"weight_g"`
	DiameterMM              float64        `json:"diameter_mm"`
	ThicknessMM             float64        `json:"thickness_mm"`
	Edge                    string         `json:"edge"`
	Shape                   string         `json:"shape"`
	Mint                    string         `json:"mint"`
	Mintage                 int64          `json:"mintage"`
	RawDetails              map[string]any `json:"raw_details"`
}
