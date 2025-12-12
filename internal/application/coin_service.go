package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

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

func (s *CoinService) AddCoin(ctx context.Context, frontFile, backFile *multipart.FileHeader, groupName, userNotes, name, mint string, mintage int, modelName string, temperature float32) (*domain.Coin, error) {
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
	analysis, err := s.aiService.AnalyzeCoin(ctx, originalFrontPath, originalBackPath, modelName, temperature)
	if err != nil {
		// Log error but continue with empty analysis to not block creation
		fmt.Printf("⚠️ Gemini analysis failed: %v\n", err)
		analysis = &domain.CoinAnalysisResult{
			Description: "Analysis failed",
			RawDetails:  map[string]any{"error": err.Error()},
		}
	}
	// 4.1 Apply Rotation Correction
	if analysis.VerticalCorrectionAngleFront != 0 {
		if err := s.imageService.Rotate(processedFrontPath, analysis.VerticalCorrectionAngleFront); err != nil {
			fmt.Printf("⚠️ Failed to rotate front image: %v\n", err)
		}
	}
	if analysis.VerticalCorrectionAngleBack != 0 {
		if err := s.imageService.Rotate(processedBackPath, analysis.VerticalCorrectionAngleBack); err != nil {
			fmt.Printf("⚠️ Failed to rotate back image: %v\n", err)
		}
	}

	// Handle Group
	var groupID *int
	groupName = strings.TrimSpace(groupName)
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
		NumistaNumber:  analysis.NumistaNumber,
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

func (s *CoinService) ListCoins(ctx context.Context, filter domain.CoinFilter) ([]*domain.Coin, error) {
	return s.repo.List(ctx, filter)
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

	// Average Value
	avgValue, err := s.repo.GetAverageValue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get average value: %w", err)
	}
	stats.AverageValue = avgValue

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

	// New Analytics
	stats.CountryDistribution, _ = s.repo.GetCountryDistribution(ctx)

	// Century & Decade Distribution (derived from GetAllCoins)
	allCoins, err := s.repo.GetAllCoins(ctx)
	if err == nil {
		stats.AllCoins = make([]domain.Coin, len(allCoins))
		stats.CenturyDistribution = make(map[string]int)
		stats.DecadeDistribution = make(map[string]int)
		var totalSilver, totalGold float64

		for i, c := range allCoins {
			stats.AllCoins[i] = *c
			if c.Year > 0 {
				// Century
				century := (c.Year-1)/100 + 1
				key := fmt.Sprintf("S. %s", toRoman(century))
				stats.CenturyDistribution[key]++

				// Decade (e.g. 1995 -> 1990s)
				decade := (c.Year / 10) * 10
				decadeKey := fmt.Sprintf("%ds", decade)
				stats.DecadeDistribution[decadeKey]++
			}

			// Calculate weights in memory to allow exclusion of Nordic Gold
			mat := strings.ToLower(c.Material)
			if strings.Contains(mat, "silver") {
				totalSilver += c.WeightG
			}
			if strings.Contains(mat, "gold") && !strings.Contains(mat, "nordic gold") {
				totalGold += c.WeightG
			}
		}
		stats.TotalSilverWeight = totalSilver
		stats.TotalGoldWeight = totalGold
	}

	stats.OldestCoin, _ = s.repo.GetOldestCoin(ctx)

	rarest, err := s.repo.GetRarestCoins(ctx, 5)
	if err == nil {
		stats.RarestCoins = make([]domain.Coin, len(rarest))
		for i, c := range rarest {
			stats.RarestCoins[i] = *c
		}
	}

	stats.GroupDistribution, _ = s.repo.GetGroupDistribution(ctx)

	// Previously fetched weights here, now calculated above

	if heaviest, err := s.repo.GetHeaviestCoin(ctx); err == nil && heaviest != nil {
		stats.HeaviestCoin = heaviest
	}

	if smallest, err := s.repo.GetSmallestCoin(ctx); err == nil && smallest != nil {
		stats.SmallestCoin = smallest
	}

	if random, err := s.repo.GetRandomCoin(ctx); err == nil && random != nil {
		stats.RandomCoin = random
	}

	return stats, nil
}

func (s *CoinService) CreateGroup(ctx context.Context, name, description string) (*domain.Group, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("group name cannot be empty")
	}
	return s.groupRepo.Create(ctx, name, description)
}

func (s *CoinService) UpdateGroup(ctx context.Context, id int, name, description string) (*domain.Group, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("group name cannot be empty")
	}

	group := &domain.Group{
		ID:          id,
		Name:        name,
		Description: description,
	}

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *CoinService) DeleteGroup(ctx context.Context, id int) error {
	return s.groupRepo.Delete(ctx, id)
}

func toRoman(n int) string {
	if n <= 0 {
		return ""
	}
	// Simplified for centuries (1-21)
	vals := []struct {
		val int
		sym string
	}{
		{21, "XXI"}, {20, "XX"}, {19, "XIX"}, {18, "XVIII"}, {17, "XVII"},
		{16, "XVI"}, {15, "XV"}, {14, "XIV"}, {13, "XIII"}, {12, "XII"},
		{11, "XI"}, {10, "X"}, {9, "IX"}, {8, "VIII"}, {7, "VII"},
		{6, "VI"}, {5, "V"}, {4, "IV"}, {3, "III"}, {2, "II"}, {1, "I"},
	}
	for _, v := range vals {
		if n == v.val {
			return v.sym
		}
	}
	return fmt.Sprintf("%d", n)
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

type UpdateCoinParams struct {
	Name           string     `json:"name"`
	Mint           string     `json:"mint"`
	Mintage        int64      `json:"mintage"`
	Country        string     `json:"country"`
	Year           int        `json:"year"`
	FaceValue      string     `json:"face_value"`
	Currency       string     `json:"currency"`
	Material       string     `json:"material"`
	Description    string     `json:"description"`
	KMCode         string     `json:"km_code"`
	MinValue       float64    `json:"min_value"`
	MaxValue       float64    `json:"max_value"`
	Grade          string     `json:"grade"`
	TechnicalNotes string     `json:"technical_notes"`
	PersonalNotes  string     `json:"personal_notes"`
	WeightG        float64    `json:"weight_g"`
	DiameterMM     float64    `json:"diameter_mm"`
	ThicknessMM    float64    `json:"thickness_mm"`
	Edge           string     `json:"edge"`
	Shape          string     `json:"shape"`
	AcquiredAt     *time.Time `json:"acquired_at"`
	SoldAt         *time.Time `json:"sold_at"`
	PricePaid      float64    `json:"price_paid"`
	SoldPrice      float64    `json:"sold_price"`
	GroupName      string     `json:"group_name"`
}

func (s *CoinService) UpdateCoin(ctx context.Context, id uuid.UUID, params UpdateCoinParams) (*domain.Coin, error) {
	coin, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}

	// Update fields
	coin.Name = params.Name
	coin.Mint = params.Mint
	coin.Mintage = params.Mintage
	coin.Country = params.Country
	coin.Year = params.Year
	coin.FaceValue = params.FaceValue
	coin.Currency = params.Currency
	coin.Material = params.Material
	coin.Description = params.Description
	coin.KMCode = params.KMCode
	coin.MinValue = params.MinValue
	coin.MaxValue = params.MaxValue
	coin.Grade = normalizeGrade(params.Grade)
	coin.TechnicalNotes = params.TechnicalNotes
	coin.PersonalNotes = params.PersonalNotes
	coin.WeightG = params.WeightG
	coin.DiameterMM = params.DiameterMM
	coin.ThicknessMM = params.ThicknessMM
	coin.Edge = params.Edge
	coin.Shape = params.Shape
	coin.AcquiredAt = params.AcquiredAt
	coin.SoldAt = params.SoldAt
	coin.PricePaid = params.PricePaid
	coin.SoldPrice = params.SoldPrice

	// Handle Group
	if params.GroupName != "" {
		group, err := s.groupRepo.GetByName(ctx, params.GroupName)
		if err != nil {
			// If not found (or other error), try to create
			group, err = s.groupRepo.Create(ctx, params.GroupName, "")
			if err != nil {
				return nil, fmt.Errorf("failed to create group: %w", err)
			}
		}
		coin.GroupID = &group.ID
	} else {
		coin.GroupID = nil
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

func (s *CoinService) downloadAndSaveImage(coinID uuid.UUID, url, filename string) (string, error) {
	// Basic implementation: fetch URL, save to storage
	// We need http client

	// START TEMPORARY FIX: Add import at top or use full package name if possible.
	// Go doesn't allow random imports. I must add "net/http" to imports.
	// Since I can't easily edit imports reliably with multi_replace without context,
	// I will assume "net/http" is available or handle it.
	// Actually, I'll just use http.Get

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	return s.storage.SaveFile(coinID, filename, resp.Body)
}

func (s *CoinService) GetGeminiModels(ctx context.Context) ([]domain.GeminiModelInfo, error) {
	return s.aiService.ListModels(ctx)
}
