package application

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/google/uuid"
)

// StorageService interface defined below

// Need to import infrastructure to use LocalFileStorage if not using interface.
// To keep clean architecture, we should define Storage interface in domain.
// For now, let's assume we pass the storage interface or struct.
// Let's define a simple Storage interface here or in domain.
// In domain/coin.go we didn't define Storage interface. Let's add it locally or use the one from infrastructure via interface.
// I'll assume we inject the dependencies.

type StorageService interface {
	SaveFile(coinID uuid.UUID, filename string, content io.Reader) (string, error)
	EnsureDir(coinID uuid.UUID) (string, error)
}

type CoinService struct {
	repo         domain.CoinRepository
	imageService domain.ImageService
	aiService    domain.AIService
	storage      StorageService
}

func NewCoinService(
	repo domain.CoinRepository,
	imageService domain.ImageService,
	aiService domain.AIService,
	storage StorageService,
) *CoinService {
	return &CoinService{
		repo:         repo,
		imageService: imageService,
		aiService:    aiService,
		storage:      storage,
	}
}

func (s *CoinService) AddCoin(ctx context.Context, frontFile, backFile *multipart.FileHeader) (*domain.Coin, error) {
	// 1. Generate ID
	coinID := uuid.New()

	// 2. Save Original Images
	frontSrc, err := frontFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open front file: %w", err)
	}
	defer frontSrc.Close()

	backSrc, err := backFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open back file: %w", err)
	}
	defer backSrc.Close()

	originalFrontPath, err := s.storage.SaveFile(coinID, "original_front.jpg", frontSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to save original front: %w", err)
	}

	originalBackPath, err := s.storage.SaveFile(coinID, "original_back.jpg", backSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to save original back: %w", err)
	}

	// 3. Analyze with Gemini (to get rotation angle)
	// We send the ORIGINALS to Gemini first? Or cropped?
	// The prompt says: "Recibir front/back -> FileSystem -> Pre-Procesamiento (Crop) -> Inteligencia (Gemini) -> Post-Procesamiento (Rotate)"
	// So:
	// a. Crop originals to circle (store as temp or overwrite?)
	// Let's crop to a new file "cropped_front.jpg"

	croppedFrontPath, err := s.imageService.CropToCircle(originalFrontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to crop front: %w", err)
	}

	croppedBackPath, err := s.imageService.CropToCircle(originalBackPath)
	if err != nil {
		return nil, fmt.Errorf("failed to crop back: %w", err)
	}

	// b. Send CROPPED images to Gemini
	analysis, err := s.aiService.AnalyzeCoin(ctx, croppedFrontPath, croppedBackPath)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze coin: %w", err)
	}

	// 4. Post-Processing: Rotate based on Gemini's suggestion
	// Gemini returns "vertical_correction_angle"
	finalFrontPath, err := s.imageService.Rotate(croppedFrontPath, analysis.VerticalCorrectionAngle)
	if err != nil {
		return nil, fmt.Errorf("failed to rotate front: %w", err)
	}

	// Assuming same rotation for back or 0. Let's use the same for now as discussed.
	finalBackPath, err := s.imageService.Rotate(croppedBackPath, analysis.VerticalCorrectionAngle)
	if err != nil {
		return nil, fmt.Errorf("failed to rotate back: %w", err)
	}

	// 5. Create Coin Entity
	coin := &domain.Coin{
		ID:                  coinID,
		Country:             analysis.Country,
		Year:                analysis.Year,
		FaceValue:           analysis.FaceValue,
		Currency:            analysis.Currency,
		Material:            analysis.Material,
		Description:         analysis.Description,
		KMCode:              analysis.KMCode,
		MinValue:            analysis.MinValue,
		MaxValue:            analysis.MaxValue,
		Grade:               analysis.Grade,
		SampleImageURLFront: finalFrontPath, // Store path or URL? For now path relative to storage or absolute?
		// Ideally we store a relative path or URL. Let's store the filename relative to storage root if possible,
		// but here we have full paths. Let's store full path for now or clean it up.
		// The frontend will need to serve this. We'll handle serving in API.
		SampleImageURLBack: finalBackPath,
		Notes:              analysis.Notes,
		GeminiDetails:      analysis.RawDetails,
	}

	// 6. Persist
	if err := s.repo.Save(ctx, coin); err != nil {
		return nil, fmt.Errorf("failed to save coin to db: %w", err)
	}

	return coin, nil
}

func (s *CoinService) ListCoins(ctx context.Context, limit, offset int) ([]*domain.Coin, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *CoinService) GetCoin(ctx context.Context, id uuid.UUID) (*domain.Coin, error) {
	return s.repo.GetByID(ctx, id)
}
