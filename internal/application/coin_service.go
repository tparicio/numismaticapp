package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

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

func (s *CoinService) AddCoin(ctx context.Context, frontFile, backFile *multipart.FileHeader, groupName, userNotes, name, mint string, mintage int) (*domain.Coin, error) {
	// 1. Process Images (Remove Background)
	// We use Rembg to remove background and get a clean PNG
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
	// Crop and center
	croppedFrontBytes, err := s.imageService.CropToContent(processedFrontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to crop front: %w", err)
	}
	processedFrontPath, err := s.storage.SaveFile(coinID, "processed_front.png", bytes.NewReader(croppedFrontBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to save processed front: %w", err)
	}

	// Process back
	processedBackBytes, err := s.bgRemover.RemoveBackground(ctx, backBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to remove background from back: %w", err)
	}
	// Crop and center
	croppedBackBytes, err := s.imageService.CropToContent(processedBackBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to crop back: %w", err)
	}
	processedBackPath, err := s.storage.SaveFile(coinID, "processed_back.png", bytes.NewReader(croppedBackBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to save processed back: %w", err)
	}

	// 4. Analyze with Gemini
	fmt.Println("🤖 Analyzing coin with Gemini...")
	analysis, err := s.aiService.AnalyzeCoin(ctx, originalFrontPath, originalBackPath)
	if err != nil {
		// Log error but continue with empty analysis to not block creation
		fmt.Printf("⚠️ Gemini analysis failed: %v\n", err)
		analysis = &domain.CoinAnalysisResult{
			Description: "Analysis failed",
			RawDetails:  map[string]any{"error": err.Error()},
		}
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
		ID:             coinID,
		Country:        analysis.Country,
		Year:           analysis.Year,
		FaceValue:      analysis.FaceValue,
		Currency:       analysis.Currency,
		Material:       analysis.Material,
		Description:    analysis.Description,
		KMCode:         analysis.KMCode,
		MinValue:       analysis.MinValue,
		MaxValue:       analysis.MaxValue,
		Grade:          normalizeGrade(analysis.Grade),
		TechnicalNotes: analysis.Notes,
		GeminiDetails:  analysis.RawDetails,
		Images:         []domain.CoinImage{},
		GroupID:        groupID,
		PersonalNotes:  userNotes,
		Name:           analysis.Name, // Use AI provided name
		Mint:           analysis.Mint,
		Mintage:        analysis.Mintage,
		WeightG:        analysis.WeightG,
		DiameterMM:     analysis.DiameterMM,
		ThicknessMM:    analysis.ThicknessMM,
		Edge:           analysis.Edge,
		Shape:          analysis.Shape,
		AcquiredAt:     nil,
		SoldAt:         nil,
		PricePaid:      0,
		SoldPrice:      0,
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

	// 5. Generate Thumbnails
	// Front
	thumbFrontPath, err := s.imageService.GenerateThumbnail(processedFrontPath, 300)
	if err != nil {
		return nil, fmt.Errorf("failed to generate front thumbnail: %w", err)
	}
	if err := addImage(thumbFrontPath, "thumbnail", "front", frontFile.Filename); err != nil {
		return nil, err
	}

	// Back
	thumbBackPath, err := s.imageService.GenerateThumbnail(processedBackPath, 300)
	if err != nil {
		return nil, fmt.Errorf("failed to generate back thumbnail: %w", err)
	}
	if err := addImage(thumbBackPath, "thumbnail", "back", backFile.Filename); err != nil {
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

func (s *CoinService) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	// Total Coins
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count coins: %w", err)
	}
	stats.TotalCoins = count

	// Total Value
	totalValue, err := s.repo.GetTotalValue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total value: %w", err)
	}
	stats.TotalValue = totalValue

	// Top Valuable
	topValuable, err := s.repo.ListTopValuable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list top valuable: %w", err)
	}
	// Convert pointers to values for struct
	stats.TopValuableCoins = make([]domain.Coin, len(topValuable))
	for i, c := range topValuable {
		stats.TopValuableCoins[i] = *c
	}

	// Recent
	recent, err := s.repo.ListRecent(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent: %w", err)
	}
	stats.RecentCoins = make([]domain.Coin, len(recent))
	for i, c := range recent {
		stats.RecentCoins[i] = *c
	}

	// Material Distribution
	matDist, err := s.repo.GetMaterialDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get material distribution: %w", err)
	}
	stats.MaterialDistribution = matDist

	// Grade Distribution
	gradeDist, err := s.repo.GetGradeDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get grade distribution: %w", err)
	}
	stats.GradeDistribution = gradeDist

	// Value Distribution
	allValues, err := s.repo.GetAllValues(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all values: %w", err)
	}

	stats.ValueDistribution = make(map[string]int)
	// Initialize buckets to ensure they exist even if 0
	buckets := []string{"0-10", "10-50", "50-100", "100-500", "500+"}
	for _, b := range buckets {
		stats.ValueDistribution[b] = 0
	}

	for _, v := range allValues {
		if v < 10 {
			stats.ValueDistribution["0-10"]++
		} else if v < 50 {
			stats.ValueDistribution["10-50"]++
		} else if v < 100 {
			stats.ValueDistribution["50-100"]++
		} else if v < 500 {
			stats.ValueDistribution["100-500"]++
		} else {
			stats.ValueDistribution["500+"]++
		}
	}

	return stats, nil
}

func normalizeGrade(input string) string {
	// Map of common variations to standard enum values
	// 'MC', 'RC', 'BC', 'MBC', 'EBC', 'SC', 'FDC', 'PROOF'
	input = strings.ToUpper(input)

	// Direct match check
	validGrades := []string{"MC", "RC", "BC", "MBC", "EBC", "SC", "FDC", "PROOF"}
	for _, g := range validGrades {
		if input == g {
			return g
		}
	}

	// Contains check for descriptive strings like "MBC (Muy Bien Conservada)"
	if strings.Contains(input, "PROOF") {
		return "PROOF"
	}
	if strings.Contains(input, "FDC") {
		return "FDC"
	}
	if strings.Contains(input, "SC") {
		return "SC"
	}
	if strings.Contains(input, "EBC") {
		return "EBC"
	}
	if strings.Contains(input, "MBC") {
		return "MBC"
	}
	if strings.Contains(input, "BC") {
		return "BC"
	}
	if strings.Contains(input, "RC") {
		return "RC"
	}
	if strings.Contains(input, "MC") {
		return "MC"
	}

	// Fallback to empty string which will be converted to NULL by repo
	return ""
}

func (s *CoinService) UpdateCoin(ctx context.Context, id uuid.UUID, groupName, userNotes, name, mint string, mintage int) (*domain.Coin, error) {
	coin, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}

	// Update fields
	if name != "" {
		coin.Name = name
	}
	if mint != "" {
		coin.Mint = mint
	}
	if mintage != 0 {
		coin.Mintage = int64(mintage)
	}
	coin.PersonalNotes = userNotes

	// Handle Group
	if groupName != "" {
		group, err := s.groupRepo.GetByName(ctx, groupName)
		if err != nil {
			// If not found (or other error), try to create
			group, err = s.groupRepo.Create(ctx, groupName, "")
			if err != nil {
				return nil, fmt.Errorf("failed to create group: %w", err)
			}
		}
		coin.GroupID = &group.ID
	}

	if err := s.repo.Update(ctx, coin); err != nil {
		return nil, fmt.Errorf("failed to update coin: %w", err)
	}

	return coin, nil
}

func (s *CoinService) DeleteCoin(ctx context.Context, id uuid.UUID) error {
	// Optional: Delete images from storage (omitted for MVP safety)
	return s.repo.Delete(ctx, id)
}
