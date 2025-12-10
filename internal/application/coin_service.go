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
	groupRepo    domain.GroupRepository
	imageService domain.ImageService
	aiService    domain.AIService
	storage      StorageService
}

func NewCoinService(
	repo domain.CoinRepository,
	groupRepo domain.GroupRepository,
	imageService domain.ImageService,
	aiService domain.AIService,
	storage StorageService,
) *CoinService {
	return &CoinService{
		repo:         repo,
		groupRepo:    groupRepo,
		imageService: imageService,
		aiService:    aiService,
		storage:      storage,
	}
}

func (s *CoinService) AddCoin(ctx context.Context, frontFile, backFile *multipart.FileHeader, groupName, userNotes string) (*domain.Coin, error) {
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
	// Let's crop to a new file "crop_front.png"
	croppedFrontPath, err := s.imageService.CropToCircle(originalFrontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to crop front: %w", err)
	}

	croppedBackPath, err := s.imageService.CropToCircle(originalBackPath)
	if err != nil {
		return nil, fmt.Errorf("failed to crop back: %w", err)
	}

	// b. Send CROPPED images to Gemini
	// analysis, err := s.aiService.AnalyzeCoin(ctx, croppedFrontPath, croppedBackPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to analyze coin: %w", err)
	// }

	// TEMPORARY: Disable Gemini API
	fmt.Println("⚠️ Gemini API disabled. Using dummy analysis.")
	analysis := &domain.CoinAnalysisResult{
		Country:                 "",
		Year:                    0,
		FaceValue:               "",
		Currency:                "",
		Material:                "",
		Description:             "Gemini analysis disabled",
		KMCode:                  "",
		MinValue:                0,
		MaxValue:                0,
		Grade:                   "",
		Notes:                   "",
		VerticalCorrectionAngle: 0,
		RawDetails:              map[string]any{"info": "Gemini disabled"},
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

	// Handle Group
	var groupID *int
	if groupName != "" {
		group, err := s.groupRepo.GetByName(ctx, groupName)
		if err != nil {
			// If not found (or other error), try to create
			// Ideally check specific error, but for MVP assuming not found
			group, err = s.groupRepo.Create(ctx, groupName, "")
			if err != nil {
				return nil, fmt.Errorf("failed to create group: %w", err)
			}
		}
		groupID = &group.ID
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
		SampleImageURLFront: finalFrontPath,
		SampleImageURLBack:  finalBackPath,
		Notes:               analysis.Notes,
		GeminiDetails:       analysis.RawDetails,
		Images:              []domain.CoinImage{},
		GroupID:             groupID,
		UserNotes:           userNotes,
	}

	// Helper to add image
	addImage := func(path, imgType, side, originalFilename string) error {
		w, h, size, mime, err := s.imageService.GetMetadata(path)
		if err != nil {
			return fmt.Errorf("failed to get metadata for %s: %w", imgType, err)
		}
		coin.Images = append(coin.Images, domain.CoinImage{
			ID:               uuid.New(),
			CoinID:           coinID,
			ImageType:        imgType,
			Side:             side,
			Path:             path,
			Extension:        "png", // TODO: Detect extension
			Size:             size,
			Width:            w,
			Height:           h,
			MimeType:         mime,
			OriginalFilename: originalFilename,
		})
		return nil
	}

	// Add all images
	if err := addImage(originalFrontPath, "original", "front", frontFile.Filename); err != nil {
		return nil, err
	}
	if err := addImage(originalBackPath, "original", "back", backFile.Filename); err != nil {
		return nil, err
	}
	if err := addImage(finalFrontPath, "crop", "front", frontFile.Filename); err != nil {
		return nil, err
	}
	if err := addImage(finalBackPath, "crop", "back", backFile.Filename); err != nil {
		return nil, err
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

func (s *CoinService) ListGroups(ctx context.Context) ([]*domain.Group, error) {
	return s.groupRepo.List(ctx)
}
