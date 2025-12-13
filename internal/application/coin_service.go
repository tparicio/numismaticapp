package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/numista" // Add import
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
	repo          domain.CoinRepository
	groupRepo     domain.GroupRepository
	imageService  domain.ImageService
	aiService     domain.AIService
	storage       StorageService
	bgRemover     domain.BackgroundRemover
	numistaClient *numista.Client
}

func NewCoinService(
	repo domain.CoinRepository,
	groupRepo domain.GroupRepository,
	imageService domain.ImageService,
	aiService domain.AIService,
	storage StorageService,
	bgRemover domain.BackgroundRemover,
	numistaClient *numista.Client,
) *CoinService {
	return &CoinService{
		repo:          repo,
		groupRepo:     groupRepo,
		imageService:  imageService,
		aiService:     aiService,
		storage:       storage,
		bgRemover:     bgRemover,
		numistaClient: numistaClient,
	}
}

func (s *CoinService) AddCoin(ctx context.Context, frontData io.Reader, frontFilename string, backData io.Reader, backFilename string, groupName, userNotes, name, mint string, mintage int, modelName string, temperature float32) (*domain.Coin, error) {
	coinID := uuid.New()
	// Start Log
	slog.Info("Starting AddCoin process", "coin_id", coinID)

	// 1. Sync: Read and Save Original Images
	frontBytes, err := io.ReadAll(frontData)
	if err != nil {
		return nil, fmt.Errorf("failed to read front file: %w", err)
	}
	backBytes, err := io.ReadAll(backData)
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

	// 2. Async: Launch Parallel Tasks
	var wg sync.WaitGroup

	// Channels for results and errors
	errChan := make(chan error, 3)

	// AI Result
	type aiResult struct {
		res *domain.CoinAnalysisResult
	}
	aiChan := make(chan aiResult, 1)

	// Image Proc Result
	type imgResult struct {
		processedFrontPath string
		processedBackPath  string
		thumbFrontPath     string
		thumbBackPath      string
	}
	imgChan := make(chan imgResult, 1)

	// Group Result
	type grpResult struct {
		id *int
	}
	grpChan := make(chan grpResult, 1)

	// Context for cancellation if one fails
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Task A: AI Analysis
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Starting Task A: AI Analysis", "coin_id", coinID, "model", modelName)
		// Assuming "es" as default language per previous refactor
		analysis, err := s.aiService.AnalyzeCoin(ctx, originalFrontPath, originalBackPath, modelName, temperature, "es")
		if err != nil {
			slog.Warn("Gemini analysis failed", "coin_id", coinID, "error", err)
			// Don't fail the whole process, just return empty/error result
			analysis = &domain.CoinAnalysisResult{
				Description: "Analysis failed",
				RawDetails:  map[string]any{"error": err.Error()},
			}
		}
		slog.Info("Completed Task A: AI Analysis", "coin_id", coinID)
		aiChan <- aiResult{res: analysis}
	}()

	// Task B: Image Processing
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Starting Task B: Image Processing", "coin_id", coinID)

		// Process Front
		pFrontBytes, err := s.bgRemover.RemoveBackground(ctx, frontBytes)
		if err != nil {
			slog.Error("Failed to remove background from front", "coin_id", coinID, "error", err)
			errChan <- fmt.Errorf("failed to bg remove front: %w", err)
			return
		}
		cFrontBytes, err := s.imageService.CropToContent(pFrontBytes)
		if err != nil {
			errChan <- fmt.Errorf("failed to crop front: %w", err)
			return
		}
		pFrontPath, err := s.storage.SaveFile(coinID, "processed_front.png", bytes.NewReader(cFrontBytes))
		if err != nil {
			errChan <- fmt.Errorf("failed to save processed front: %w", err)
			return
		}
		tFrontPath, err := s.imageService.GenerateThumbnail(pFrontPath, 300)
		if err != nil {
			errChan <- fmt.Errorf("failed to thumb front: %w", err)
			return
		}

		// Process Back
		pBackBytes, err := s.bgRemover.RemoveBackground(ctx, backBytes)
		if err != nil {
			slog.Error("Failed to remove background from back", "coin_id", coinID, "error", err)
			errChan <- fmt.Errorf("failed to bg remove back: %w", err)
			return
		}
		cBackBytes, err := s.imageService.CropToContent(pBackBytes)
		if err != nil {
			errChan <- fmt.Errorf("failed to crop back: %w", err)
			return
		}
		pBackPath, err := s.storage.SaveFile(coinID, "processed_back.png", bytes.NewReader(cBackBytes))
		if err != nil {
			errChan <- fmt.Errorf("failed to save processed back: %w", err)
			return
		}
		tBackPath, err := s.imageService.GenerateThumbnail(pBackPath, 300)
		if err != nil {
			errChan <- fmt.Errorf("failed to thumb back: %w", err)
			return
		}

		imgChan <- imgResult{
			processedFrontPath: pFrontPath,
			processedBackPath:  pBackPath,
			thumbFrontPath:     tFrontPath,
			thumbBackPath:      tBackPath,
		}
		slog.Info("Completed Task B: Image Processing", "coin_id", coinID)
	}()

	// Task C: Group Management
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Starting Task C: Group Management", "coin_id", coinID)
		groupName = strings.TrimSpace(groupName)
		if groupName == "" {
			slog.Info("No group name provided, skipping group creation", "coin_id", coinID)
			grpChan <- grpResult{id: nil}
			return
		}

		group, err := s.groupRepo.GetByName(ctx, groupName)
		if err != nil {
			// Try create
			group, err = s.groupRepo.Create(ctx, groupName, "")
			if err != nil {
				slog.Error("Failed to create group", "coin_id", coinID, "group_name", groupName, "error", err)
				errChan <- fmt.Errorf("failed to create group: %w", err)
				return
			}
		}
		slog.Info("Completed Task C: Group Management", "coin_id", coinID, "group_id", group.ID)
		grpChan <- grpResult{id: &group.ID}
	}()

	// Wait for all
	wg.Wait()
	close(errChan)
	close(aiChan)
	close(imgChan)
	close(grpChan)

	// Check for critical errors (Image Proc or Group)
	// AI error logic is handled inside to be non-critical
	for err := range errChan {
		if err != nil {
			return nil, err // Return first error
		}
	}

	analysisRes := (<-aiChan).res
	imgRes := <-imgChan
	grpRes := <-grpChan

	// 5. Assemble Coin Entity
	coin := &domain.Coin{
		ID:                coinID,
		Country:           analysisRes.Country,
		Year:              analysisRes.Year,
		FaceValue:         analysisRes.FaceValue,
		Currency:          analysisRes.Currency,
		Material:          analysisRes.Material,
		Description:       analysisRes.Description,
		KMCode:            analysisRes.KMCode,
		NumistaNumber:     analysisRes.NumistaNumber,
		MinValue:          analysisRes.MinValue,
		MaxValue:          analysisRes.MaxValue,
		Grade:             normalizeGrade(analysisRes.Grade),
		TechnicalNotes:    analysisRes.Notes,
		GeminiDetails:     analysisRes.RawDetails,
		Images:            []domain.CoinImage{},
		GroupID:           grpRes.id,
		PersonalNotes:     userNotes,
		Name:              analysisRes.Name,
		Mint:              analysisRes.Mint,
		Mintage:           analysisRes.Mintage,
		WeightG:           analysisRes.WeightG,
		DiameterMM:        analysisRes.DiameterMM,
		ThicknessMM:       analysisRes.ThicknessMM,
		Edge:              analysisRes.Edge,
		Shape:             analysisRes.Shape,
		AcquiredAt:        nil, // TODO
		SoldAt:            nil,
		PricePaid:         0,
		SoldPrice:         0,
		GeminiModel:       modelName,
		GeminiTemperature: float64(temperature),
	}

	// Helper to add image RECORD
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
			Extension:        "png", // processed are png
			Size:             size,
			Width:            w,
			Height:           h,
			MimeType:         mime,
			OriginalFilename: originalFilename,
		})
		return nil
	}
	// For original jpgs
	addOriginal := func(path, side, originalFilename string) error {
		w, h, size, mime, err := s.imageService.GetMetadata(path)
		if err != nil {
			return fmt.Errorf("failed to get metadata for original: %w", err)
		}
		coin.Images = append(coin.Images, domain.CoinImage{
			ID:               uuid.New(),
			CoinID:           coinID,
			ImageType:        "original",
			Side:             side,
			Path:             path,
			Extension:        "jpg",
			Size:             size,
			Width:            w,
			Height:           h,
			MimeType:         mime,
			OriginalFilename: originalFilename,
		})
		return nil
	}

	if err := addOriginal(originalFrontPath, "front", frontFilename); err != nil {
		return nil, err
	}
	if err := addOriginal(originalBackPath, "back", backFilename); err != nil {
		return nil, err
	}
	if err := addImage(imgRes.processedFrontPath, "crop", "front", frontFilename); err != nil {
		return nil, err
	}
	if err := addImage(imgRes.processedBackPath, "crop", "back", backFilename); err != nil {
		return nil, err
	}
	if err := addImage(imgRes.thumbFrontPath, "thumbnail", "front", frontFilename); err != nil {
		return nil, err
	}
	if err := addImage(imgRes.thumbBackPath, "thumbnail", "back", backFilename); err != nil {
		return nil, err
	}

	// 6. Persist
	if err := s.repo.Save(ctx, coin); err != nil {
		slog.Error("Failed to save coin to DB", "coin_id", coinID, "error", err)
		return nil, fmt.Errorf("failed to save coin to db: %w", err)
	}
	slog.Info("Successfully saved coin", "coin_id", coinID)

	// 7. Trigger Numista Enrichment (Async)
	if s.numistaClient != nil && s.numistaClient.APIKey != "" {
		go func(id uuid.UUID) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			if err := s.EnrichCoinWithNumista(bgCtx, id); err != nil {
				slog.Error("Failed to enrich coin with Numista", "coin_id", id, "error", err)
			} else {
				slog.Info("Successfully enriched coin with Numista", "coin_id", id)
			}
		}(coin.ID)
	}

	return coin, nil
}

func (s *CoinService) EnrichCoinWithNumista(ctx context.Context, coinID uuid.UUID) error {
	// 1. Get Coin
	coin, err := s.repo.GetByID(ctx, coinID)
	if err != nil {
		return fmt.Errorf("failed to get coin: %w", err)
	}

	// 2. Search Numista (Top 10)
	query := fmt.Sprintf("%s %s %s", coin.FaceValue, coin.Currency, coin.Name)
	query = strings.TrimSpace(query)
	issuer := coin.Country
	yearStr := ""
	if coin.Year > 0 {
		yearStr = fmt.Sprintf("%d", coin.Year)
	}

	// Search requesting 10 results
	searchResult, err := s.numistaClient.SearchTypes(ctx, query, "coin", yearStr, issuer, 10)
	if err != nil {
		return fmt.Errorf("numista search failed: %w", err)
	}

	if searchResult == nil {
		slog.Info("Numista returned no results", "coin_id", coinID)
		return nil
	}

	// 3. Save full search response
	searchJSON, err := json.Marshal(searchResult)
	if err != nil {
		slog.Warn("Failed to marshal numista search result", "error", err)
	} else {
		coin.NumistaSearch = string(searchJSON)
	}

	// If no types found, just save the search result and exit
	if searchResult.Count == 0 || len(searchResult.Types) == 0 {
		if err := s.repo.Update(ctx, coin); err != nil {
			return fmt.Errorf("failed to update coin with numista search: %w", err)
		}
		return nil
	}

	// 4. Get Details for the first result
	firstTypeID := searchResult.Types[0].ID
	typeDetails, err := s.numistaClient.GetType(ctx, firstTypeID)
	if err != nil {
		return fmt.Errorf("failed to get numista type details: %w", err)
	}

	// 5. Update Coin
	coin.NumistaNumber = firstTypeID
	coin.NumistaDetails = typeDetails

	// Example mapping from Numista Details (structure depends on API response)
	// Usually "title", "description", "year", etc. are top level or nested.
	// For now we trust the raw map is saved and maybe we update some basic fields if empty.
	// We can refine this mapping based on the actual JSON structure of 'typeDetails'.

	if err := s.repo.Update(ctx, coin); err != nil {
		return fmt.Errorf("failed to update coin with numista details: %w", err)
	}

	return nil
}

// Helper to bridge side string to enum if needed in internal logic,
// but for SaveFile we pass filename explicitly.
// Inside saveNumistaImage logic above I need to fix Side.
// Let's rewrite saveNumistaImage slightly to take filename and dbSide.

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

	// Group Stats for Widget
	stats.GroupStats, _ = s.repo.GetGroupStats(ctx)

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

func (s *CoinService) RotateCoinImage(ctx context.Context, coinID uuid.UUID, side string, angle float64) error {
	slog.Info("RotateCoinImage called", "coin_id", coinID, "side", side, "angle", angle)
	coin, err := s.repo.GetByID(ctx, coinID)
	if err != nil {
		return fmt.Errorf("failed to get coin: %w", err)
	}

	// Find the processed image for the side
	var targetImg *domain.CoinImage
	for i := range coin.Images {
		if coin.Images[i].Side == side && coin.Images[i].ImageType == "crop" {
			targetImg = &coin.Images[i]
			break
		}
	}

	if targetImg == nil {
		return fmt.Errorf("processed image not found for side %s", side)
	}

	slog.Info("Found image to rotate", "path", targetImg.Path)

	// Rotate the image
	// Note: Rotate modifies the file in place
	if err := s.imageService.Rotate(targetImg.Path, angle); err != nil {
		slog.Error("Rotation failed", "error", err)
		return fmt.Errorf("failed to rotate image: %w", err)
	}
	slog.Info("Rotation successful")

	// Regenerate thumbnail
	// VipsImageService.GenerateThumbnail derives name from input: path + "_thumb.png"
	// Our naming convention in AddCoin is "processed_front.png" -> "processed_front_thumb.png"
	thumbPath, err := s.imageService.GenerateThumbnail(targetImg.Path, 300)
	if err != nil {
		slog.Error("Thumbnail regeneration failed", "error", err)
		return fmt.Errorf("failed to regenerate thumbnail: %w", err)
	}
	slog.Info("Thumbnail regenerated", "path", thumbPath)

	return nil
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
	// Since I can't easily edit imports reliably weplace without context,
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

func (s *CoinService) ReanalyzeCoin(ctx context.Context, id uuid.UUID, modelName string, temperature float32) (*domain.Coin, error) {
	// 1. Get Coin
	coin, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("coin not found: %w", err)
	}

	// 2. Find Original Images
	var frontPath, backPath string
	for _, img := range coin.Images {
		if img.ImageType == "original" {
			if img.Side == "front" {
				frontPath = img.Path
			} else if img.Side == "back" {
				backPath = img.Path
			}
		}
	}

	if frontPath == "" || backPath == "" {
		return nil, fmt.Errorf("original images not found for this coin")
	}

	// 3. Analyze with Gemini
	slog.Info("Re-analyzing coin with Gemini", "coin_id", id)
	analysis, err := s.aiService.AnalyzeCoin(ctx, frontPath, backPath, modelName, temperature, "es")
	if err != nil {
		return nil, fmt.Errorf("gemini analysis failed: %w", err)
	}

	// 4. Update Coin Fields from Analysis
	coin.Country = analysis.Country
	coin.Year = analysis.Year
	coin.FaceValue = analysis.FaceValue
	coin.Currency = analysis.Currency
	coin.Material = analysis.Material
	coin.Description = analysis.Description
	coin.KMCode = analysis.KMCode
	coin.NumistaNumber = analysis.NumistaNumber
	coin.MinValue = analysis.MinValue
	coin.MaxValue = analysis.MaxValue
	coin.Grade = normalizeGrade(analysis.Grade)
	coin.TechnicalNotes = analysis.Notes
	coin.GeminiDetails = analysis.RawDetails
	coin.Name = analysis.Name
	coin.Mint = analysis.Mint
	coin.Mintage = analysis.Mintage
	coin.WeightG = analysis.WeightG
	coin.DiameterMM = analysis.DiameterMM
	coin.ThicknessMM = analysis.ThicknessMM
	coin.Edge = analysis.Edge
	coin.Shape = analysis.Shape
	coin.GeminiModel = modelName
	coin.GeminiTemperature = float64(temperature)
	// We don't overwrite UserNotes, AddedAt, etc.

	// 5. Update in Repo
	// We use the same Update method but need to orchestrate it via UpdateCoin or direct repo update.
	// Since UpdateCoin accepts a param struct, we might need a direct save or construct the params.
	// Let's us direct repo update since we modified the entity directly.
	if err := s.repo.Update(ctx, coin); err != nil {
		return nil, fmt.Errorf("failed to update coin: %w", err)
	}

	return coin, nil
}
