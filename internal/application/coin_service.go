package application

import (
	"bytes"
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
	bgRemover    domain.BackgroundRemover
}

func NewCoinService(
	repo domain.CoinRepository,
	groupRepo domain.GroupRepository,
	imageService domain.ImageService,
	aiService domain.AIService,
	storage StorageService,
	bgRemover domain.BackgroundRemover,
) *CoinService {
	return &CoinService{
		repo:         repo,
		groupRepo:    groupRepo,
		imageService: imageService,
		aiService:    aiService,
		storage:      storage,
		bgRemover:    bgRemover,
	}
}

func (s *CoinService) AddCoin(ctx context.Context, frontFile, backFile *multipart.FileHeader, groupName, userNotes string) (*domain.Coin, error) {
	// 1. Generate ID
	coinID := uuid.New()

	// 2. Save Original Images & Process Background
	frontSrc, err := frontFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open front file: %w", err)
	}
	defer frontSrc.Close()

	if _, err := frontSrc.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek front file: %w", err)
	}

	frontBytes, err := io.ReadAll(frontSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to read front file: %w", err)
	}

	backSrc, err := backFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open back file: %w", err)
	}
	defer backSrc.Close()

	if _, err := backSrc.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek back file: %w", err)
	}

	backBytes, err := io.ReadAll(backSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to read back file: %w", err)
	}

	originalFrontPath, err := s.storage.SaveFile(coinID, "original_front.jpg", bytes.NewReader(frontBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to save original front: %w", err)
	}

	originalBackPath, err := s.storage.SaveFile(coinID, "original_back.jpg", bytes.NewReader(backBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to save original back: %w", err)
	}

	// 3. Remove Background (Rembg)
	// Process front
	processedFrontBytes, err := s.bgRemover.RemoveBackground(ctx, frontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to remove background from front: %w", err)
	}
	processedFrontPath, err := s.storage.SaveFile(coinID, "processed_front.png", bytes.NewReader(processedFrontBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to save processed front: %w", err)
	}

	// Process back
	processedBackBytes, err := s.bgRemover.RemoveBackground(ctx, backBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to remove background from back: %w", err)
	}
	processedBackPath, err := s.storage.SaveFile(coinID, "processed_back.png", bytes.NewReader(processedBackBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to save processed back: %w", err)
	}

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
	// finalFrontPath, err := s.imageService.Rotate(originalFrontPath, analysis.VerticalCorrectionAngle)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to rotate front: %w", err)
	// }

	// Assuming same rotation for back or 0. Let's use the same for now as discussed.
	// finalBackPath, err := s.imageService.Rotate(originalBackPath, analysis.VerticalCorrectionAngle)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to rotate back: %w", err)
	// }

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
		SampleImageURLFront: originalFrontPath, // Was finalFrontPath
		SampleImageURLBack:  originalBackPath,  // Was finalBackPath
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
	if err := addImage(processedFrontPath, "crop", "front", frontFile.Filename); err != nil {
		return nil, err
	}
	if err := addImage(processedBackPath, "crop", "back", backFile.Filename); err != nil {
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
