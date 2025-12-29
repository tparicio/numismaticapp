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
	Mintage        Mintage        `json:"mintage"`
	Country        string         `json:"country"`
	Year           Year           `json:"year"`
	FaceValue      string         `json:"face_value"`
	Currency       string         `json:"currency"`
	Material       string         `json:"material"`
	Description    string         `json:"description"`
	KMCode         KMCode         `json:"km_code"`
	NumistaNumber  int            `json:"numista_number"`
	NumistaDetails map[string]any `json:"numista_details"`
	// Detailed fields
	Ruler             string         `json:"ruler"`
	Orientation       string         `json:"orientation"`
	Series            string         `json:"series"`
	CommemoratedTopic string         `json:"commemorated_topic"`
	MinValue          float64        `json:"min_value"`
	MaxValue          float64        `json:"max_value"`
	Grade             Grade          `json:"grade"`
	TechnicalNotes    string         `json:"technical_notes"`
	GeminiDetails     map[string]any `json:"gemini_details"` // Raw JSON from Gemini
	GeminiModel       string         `json:"gemini_model"`
	GeminiTemperature float64        `json:"gemini_temperature"`
	NumistaSearch     string         `json:"numista_search"`
	Images            []CoinImage    `json:"images"`
	GroupID           *int           `json:"group_id"`
	PersonalNotes     string         `json:"personal_notes"`
	WeightG           float64        `json:"weight_g"`
	DiameterMM        float64        `json:"diameter_mm"`
	ThicknessMM       float64        `json:"thickness_mm"`
	Edge              string         `json:"edge"`
	Shape             string         `json:"shape"`
	AcquiredAt        *time.Time     `json:"acquired_at"`
	SoldAt            *time.Time     `json:"sold_at"`
	PricePaid         float64        `json:"price_paid"`
	SoldPrice         float64        `json:"sold_price"`
	SaleChannel       string         `json:"sale_channel"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type Group struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CoinCount   int       `json:"coin_count"` // Populated for display purposes
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

// CoinFilter defines available filters for listing coins.
type CoinFilter struct {
	Limit     int
	Offset    int
	GroupID   *int
	Year      *int
	Country   *string
	Query     *string
	MinPrice  *float64
	MaxPrice  *float64
	Grade     *string
	Material  *string
	MinYear   *int
	MaxYear   *int
	SortBy    *string
	SortOrder *string
}

// CoinRepository defines the interface for persisting coins.
type CoinRepository interface {
	Save(ctx context.Context, coin *Coin) error
	GetByID(ctx context.Context, id uuid.UUID) (*Coin, error)
	List(ctx context.Context, filter CoinFilter) ([]*Coin, error)
	Count(ctx context.Context) (int64, error)
	GetTotalValue(ctx context.Context) (float64, error)
	GetAverageValue(ctx context.Context) (float64, error)
	ListTopValuable(ctx context.Context) ([]*Coin, error)
	ListRecent(ctx context.Context) ([]*Coin, error)
	GetMaterialDistribution(ctx context.Context) (map[string]int, error)
	GetGradeDistribution(ctx context.Context) (map[string]int, error)
	GetAllValues(ctx context.Context) ([]float64, error)
	Update(ctx context.Context, coin *Coin) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetCountryDistribution(ctx context.Context) (map[string]int, error)
	GetOldestCoin(ctx context.Context) (*Coin, error)
	GetRarestCoins(ctx context.Context, limit int) ([]*Coin, error)
	GetGroupDistribution(ctx context.Context) (map[string]int, error)
	GetGroupStats(ctx context.Context) ([]GroupStat, error)
	GetTotalWeightByMaterial(ctx context.Context, materialLike string) (float64, error)
	GetHeaviestCoin(ctx context.Context) (*Coin, error)
	GetSmallestCoin(ctx context.Context) (*Coin, error)
	GetRandomCoin(ctx context.Context) (*Coin, error)
	GetAllCoins(ctx context.Context) ([]*Coin, error)
	AddImage(ctx context.Context, image CoinImage) error
	// Sell operations
	MarkAsSold(ctx context.Context, id uuid.UUID, soldAt time.Time, soldPrice float64, saleChannel string) (*Coin, error)
	GetSaleChannels(ctx context.Context) ([]string, error)
	// Link operations
	AddLink(ctx context.Context, link *CoinLink) error
	RemoveLink(ctx context.Context, linkID uuid.UUID) error
	GetLink(ctx context.Context, linkID uuid.UUID) (*CoinLink, error)
	UpdateLink(ctx context.Context, link *CoinLink) error
	ListLinks(ctx context.Context, coinID uuid.UUID) ([]*CoinLink, error)
	// Bulk operations for export
	GetAllImages(ctx context.Context) ([]CoinImage, error)
	GetAllLinks(ctx context.Context) ([]*CoinLink, error)
}

// CoinLink represents an external link associated with a coin.
type CoinLink struct {
	ID            uuid.UUID `json:"id"`
	CoinID        uuid.UUID `json:"coin_id"`
	URL           string    `json:"url"`
	Name          string    `json:"name"`
	OGTitle       string    `json:"og_title"`
	OGDescription string    `json:"og_description"`
	OGImage       string    `json:"og_image"`
	CreatedAt     time.Time `json:"created_at"`
}

// GroupRepository defines the interface for persisting groups.
type GroupRepository interface {
	Create(ctx context.Context, name, description string) (*Group, error)
	GetByName(ctx context.Context, name string) (*Group, error)
	List(ctx context.Context) ([]*Group, error)
	Update(ctx context.Context, group *Group) error
	Delete(ctx context.Context, id int) error
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
	// Rotate rotates the image at the given path by the specified angle.
	Rotate(imagePath string, angle float64) error
}

type GeminiModelInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AIService defines the interface for AI analysis.
type AIService interface {
	// AnalyzeCoin analyzes the front and back images of a coin and returns metadata.
	AnalyzeCoin(ctx context.Context, frontImagePath, backImagePath, modelName string, temperature float32, lang string) (*CoinAnalysisResult, error)
	// ListModels returns a list of available Gemini models.
	ListModels(ctx context.Context) ([]GeminiModelInfo, error)
}

// CoinAnalysisResult contains the data extracted by the AI.
type CoinAnalysisResult struct {
	Country                      string         `json:"country"`
	Year                         int            `json:"year"`
	FaceValue                    string         `json:"face_value"`
	Currency                     string         `json:"currency"`
	Material                     string         `json:"material"`
	Description                  string         `json:"description"`
	KMCode                       string         `json:"km_code"`
	NumistaNumber                int            `json:"numista_number"`
	MinValue                     float64        `json:"min_value"`
	MaxValue                     float64        `json:"max_value"`
	Grade                        string         `json:"grade"`
	Name                         string         `json:"name"`
	Notes                        string         `json:"notes"`
	VerticalCorrectionAngleFront float64        `json:"vertical_correction_angle_front"`
	VerticalCorrectionAngleBack  float64        `json:"vertical_correction_angle_back"`
	WeightG                      float64        `json:"weight_g"`
	DiameterMM                   float64        `json:"diameter_mm"`
	ThicknessMM                  float64        `json:"thickness_mm"`
	Edge                         string         `json:"edge"`
	Shape                        string         `json:"shape"`
	Mint                         string         `json:"mint"`
	Mintage                      int64          `json:"mintage"`
	ReferenceSourceName          string         `json:"reference_source_name"`
	RawDetails                   map[string]any `json:"raw_details"`
}

type PriceClient interface {
	GetMetalPrices(ctx context.Context) (float64, float64, error)
}
