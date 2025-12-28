package application

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/numista" // Add import
	"github.com/google/uuid"
)

type NumistaService interface {
	SearchTypes(ctx context.Context, query, category, year, issuer string, count int) (*numista.TypeSearchResponse, error)
	GetType(ctx context.Context, id int) (map[string]any, error)
}

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
	DeleteCoinDirectory(coinID uuid.UUID) error
}

type CoinService struct {
	repo          domain.CoinRepository
	groupRepo     domain.GroupRepository
	imageService  domain.ImageService
	aiService     domain.AIService
	storage       StorageService
	bgRemover     domain.BackgroundRemover
	numistaClient NumistaService
}

func NewCoinService(
	repo domain.CoinRepository,
	groupRepo domain.GroupRepository,
	imageService domain.ImageService,
	aiService domain.AIService,
	storage StorageService,
	bgRemover domain.BackgroundRemover,
	numistaClient NumistaService,
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
	yearVO, err := domain.NewYear(analysisRes.Year)
	if err != nil {
		slog.Warn("AI returned invalid year", "coin_id", coinID, "year", analysisRes.Year, "error", err)
		// We can default to 0 (Unknown) or keep invalid if our validation was strict.
		// Since validation returned error, let's use 0.
		yearVO, _ = domain.NewYear(0)
	}

	kmVO, _ := domain.NewKMCode(analysisRes.KMCode)
	gradeVO, _ := domain.NewGrade(normalizeGrade(analysisRes.Grade))
	mintageVO, _ := domain.NewMintage(analysisRes.Mintage)

	coin := &domain.Coin{
		ID:                coinID,
		Country:           analysisRes.Country,
		Year:              yearVO,
		FaceValue:         analysisRes.FaceValue,
		Currency:          analysisRes.Currency,
		Material:          analysisRes.Material,
		Description:       analysisRes.Description,
		KMCode:            kmVO,
		NumistaNumber:     analysisRes.NumistaNumber,
		MinValue:          analysisRes.MinValue,
		MaxValue:          analysisRes.MaxValue,
		Grade:             gradeVO,
		TechnicalNotes:    analysisRes.Notes,
		GeminiDetails:     analysisRes.RawDetails,
		Images:            []domain.CoinImage{},
		GroupID:           grpRes.id,
		PersonalNotes:     userNotes,
		Name:              analysisRes.Name,
		Mint:              analysisRes.Mint,
		Mintage:           mintageVO,
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
	if s.numistaClient != nil {
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
	slog.Info("Starting Numista enrichment", "coin_id", coinID)
	// 1. Get Coin
	coin, err := s.repo.GetByID(ctx, coinID)
	if err != nil {
		return fmt.Errorf("failed to get coin: %w", err)
	}

	// 2. Search Numista (Top 10)
	queryParts := []string{coin.FaceValue, coin.Currency}
	if coin.Country != "" {
		queryParts = append(queryParts, coin.Country)
	}
	query := strings.Join(queryParts, " ")
	query = strings.Join(strings.Fields(query), " ")

	yearStr := ""
	if coin.Year.Int() > 0 {
		yearStr = fmt.Sprintf("%d", coin.Year.Int())
	}

	slog.Info("Searching Numista", "query", query, "year", yearStr)

	searchResult, err := s.numistaClient.SearchTypes(ctx, query, "coin", yearStr, "", 10)
	if err != nil {
		slog.Error("Numista search request failed", "error", err)
		return fmt.Errorf("numista search failed: %w", err)
	}

	if searchResult == nil {
		slog.Info("Numista returned nil result", "coin_id", coinID)
		return nil
	}

	slog.Info("Numista search completed", "count", searchResult.Count, "results_len", len(searchResult.Types))

	// 3. Save full search response (with query injected)
	searchMap := map[string]interface{}{
		"q":     query,
		"count": searchResult.Count,
		"types": searchResult.Types,
	}
	searchJSON, err := json.Marshal(searchMap)
	if err != nil {
		slog.Warn("Failed to marshal numista search result", "error", err)
	} else {
		coin.NumistaSearch = string(searchJSON)
		slog.Info("Set NumistaSearch field", "length", len(coin.NumistaSearch))
	}

	if searchResult.Count == 0 || len(searchResult.Types) == 0 {
		slog.Info("No types found in Numista, saving empty search result", "coin_id", coinID)
		if err := s.repo.Update(ctx, coin); err != nil {
			return fmt.Errorf("failed to update coin with numista search: %w", err)
		}
		return nil
	}

	// Helper to extract numeric value from string (e.g. "50 Pesetas" -> 50.0)
	parseNumeric := func(input string) float64 {
		var val float64
		// Try scanning first number
		_, err := fmt.Sscanf(input, "%f", &val)
		if err == nil {
			return val
		}
		// If failed, maybe regex or simple iteration?
		// Simple approach: replace comma with dot, extract first sequence of digits/dots
		// For now simple scan is a good start.
		return 0 // fail
	}

	targetValue := parseNumeric(coin.FaceValue)
	slog.Info("Target numeric value", "raw", coin.FaceValue, "parsed", targetValue)

	var bestMatchDetails map[string]any
	var bestMatchID int
	var fallbackMatchDetails map[string]any // Matches value but maybe not year
	var fallbackMatchID int

	// 4. Iterate Candidates
	for i, candidate := range searchResult.Types {
		slog.Info("Checking candidate", "index", i, "id", candidate.ID, "title", candidate.Title, "min_year", candidate.MinYear, "max_year", candidate.MaxYear)

		// Fast Filter: Year Range
		// Only strictly check year if we have one
		yearMatches := true
		if coin.Year.Int() > 0 {
			if coin.Year.Int() < candidate.MinYear || coin.Year.Int() > candidate.MaxYear {
				yearMatches = false
				slog.Debug("Year mismatch", "coin_year", coin.Year.Int(), "range", fmt.Sprintf("%d-%d", candidate.MinYear, candidate.MaxYear))
			}
		}

		// Fetch matches details to check value
		// Warning: This makes an API call per candidate in loop
		details, err := s.numistaClient.GetType(ctx, candidate.ID)
		if err != nil {
			slog.Warn("Failed to get details for candidate, skipping", "id", candidate.ID, "error", err)
			continue
		}

		// Check Value Match
		var detailsValue float64
		if valMap, ok := details["value"].(map[string]any); ok {
			if numVal, ok := valMap["numeric_value"].(float64); ok {
				detailsValue = numVal
			}
		}

		valuesMatch := func(v1, v2 float64) bool {
			if v1 <= 0 || v2 <= 0 {
				return false
			}
			epsilon := 0.001
			// Exact match
			if v1 > v2-epsilon && v1 < v2+epsilon {
				return true
			}
			// v1 is cents (20), v2 is unit (0.2) -> v1 = v2 * 100
			if v1 > (v2*100)-epsilon && v1 < (v2*100)+epsilon {
				return true
			}
			// v1 is unit (0.2), v2 is cents (20) -> v1 * 100 = v2
			if (v1*100) > v2-epsilon && (v1*100) < v2+epsilon {
				return true
			}
			return false
		}

		valueMatches := valuesMatch(detailsValue, targetValue)
		slog.Info("Details fetched", "id", candidate.ID, "details_value", detailsValue, "target_value", targetValue, "value_matches", valueMatches, "year_matches", yearMatches)

		if valueMatches && yearMatches {
			slog.Info("PERFECT MATCH FOUND", "id", candidate.ID)
			bestMatchDetails = details
			bestMatchID = candidate.ID
			break // Stop searching
		}

		if valueMatches && fallbackMatchDetails == nil {
			slog.Info("Fallback match found (Value matches, year mismatch)", "id", candidate.ID)
			fallbackMatchDetails = details
			fallbackMatchID = candidate.ID
			// Continue searching for a perfect match...
		}

		// Capture first result as ultimate fallback if we have nothing else yet?
		// User requirement: "si ninguna cumple ambos... nos quedaremos con la primera que cumpla face_value".
		// This implies if we finish loop and no perfect match, we use fallbackMatchDetails.

		// If we don't have even a fallback, should we take the first result?
		// User didn't explicitly say "default to first result if NOTHING matches",
		// but standard behavior usually implies keeping the "closest" or at least "something".
		// Implemented logic: Priority 1: Perfect. Priority 2: Fallback (Value only).
	}

	finalDetails := bestMatchDetails
	finalID := bestMatchID

	if finalDetails == nil {
		if fallbackMatchDetails != nil {
			slog.Info("No perfect match found. Using fallback (Value match).")
			finalDetails = fallbackMatchDetails
			finalID = fallbackMatchID
		} else {
			slog.Info("No match found for Value. No update performed for Details.")
			// If user wants default to first result, uncomment below:
			/*
				slog.Info("Defaulting to first result.")
				finalDetails, _ = s.numistaClient.GetType(ctx, searchResult.Types[0].ID)
				coin.NumistaNumber = searchResult.Types[0].ID
			*/
			// For now obeying "recordar en el bucle... primera que cumpla face_value".
			// If none comply, we do nothing (except saving parsing search results).
		}
	}

	if finalDetails != nil {
		coin.NumistaDetails = finalDetails
		coin.NumistaNumber = finalID

		s.mapNumistaDetails(coin, finalDetails)

		slog.Info("Persisting Numista details", "coin_id", coinID, "numista_id", coin.NumistaNumber)
		if err := s.repo.Update(ctx, coin); err != nil {
			slog.Error("Failed to persist coin updates", "error", err)
			return fmt.Errorf("failed to update coin with numista details: %w", err)
		}
	} else {
		// Even if no details applied, we must save the NumistaSearch field we set earlier
		if err := s.repo.Update(ctx, coin); err != nil {
			return fmt.Errorf("failed to update coin (search only): %w", err)
		}
	}

	slog.Info("Numista enrichment finished", "coin_id", coinID)

	return nil
}

func (s *CoinService) ApplyNumistaCandidate(ctx context.Context, coinID uuid.UUID, numistaID int) (*domain.Coin, error) {
	slog.Info("Applying Numista candidate manually", "coin_id", coinID, "numista_id", numistaID)

	coin, err := s.repo.GetByID(ctx, coinID)
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}

	details, err := s.numistaClient.GetType(ctx, numistaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get numista details: %w", err)
	}

	coin.NumistaDetails = details
	coin.NumistaNumber = numistaID

	s.mapNumistaDetails(coin, details)

	if err := s.repo.Update(ctx, coin); err != nil {
		return nil, fmt.Errorf("failed to update coin: %w", err)
	}

	return coin, nil
}

func (s *CoinService) mapNumistaDetails(coin *domain.Coin, details map[string]any) {
	// 1. Dimensions & Weight
	if v, ok := details["size"].(float64); ok {
		coin.DiameterMM = v
	}
	if v, ok := details["thickness"].(float64); ok {
		coin.ThicknessMM = v
	}
	if v, ok := details["weight"].(float64); ok {
		coin.WeightG = v
	}

	// 2. Shape
	if v, ok := details["shape"].(string); ok {
		coin.Shape = v
	}

	// 3. Material (Composition)
	if comp, ok := details["composition"].(map[string]any); ok {
		if text, ok := comp["text"].(string); ok {
			coin.Material = text
		}
	}

	// 4. Mints (Pick first)
	if mints, ok := details["mints"].([]any); ok && len(mints) > 0 {
		if firstMint, ok := mints[0].(map[string]any); ok {
			if name, ok := firstMint["name"].(string); ok {
				coin.Mint = name
			}
		}
	}

	// 5. KM Code (Loop references)
	if refs, ok := details["references"].([]any); ok {
		for _, r := range refs {
			if refMap, ok := r.(map[string]any); ok {
				if cat, ok := refMap["catalogue"].(map[string]any); ok {
					if code, ok := cat["code"].(string); ok && code == "KM" {
						if number, ok := refMap["number"].(string); ok {
							coin.KMCode, _ = domain.NewKMCode(fmt.Sprintf("KM# %s", number))
							break // Found KM
						}
					}
				}
			}
		}
	}

	// 6. Ruler (First from list)
	if rulers, ok := details["ruler"].([]any); ok && len(rulers) > 0 {
		if firstRuler, ok := rulers[0].(map[string]any); ok {
			if name, ok := firstRuler["name"].(string); ok {
				coin.Ruler = name
			}
		}
	}

	// 7. Orientation
	if v, ok := details["orientation"].(string); ok {
		coin.Orientation = v
	}

	// 8. Series
	if v, ok := details["series"].(string); ok {
		coin.Series = v
	}

	// 9. Commemorated Topic
	if v, ok := details["commemorated_topic"].(string); ok {
		coin.CommemoratedTopic = v
	}

	slog.Info("Mapped Numista details to coin fields",
		"diameter", coin.DiameterMM,
		"weight", coin.WeightG,
		"mint", coin.Mint,
		"km", coin.KMCode,
		"ruler", coin.Ruler,
		"series", coin.Series,
		"orientation", coin.Orientation,
		"commemorated_topic", coin.CommemoratedTopic,
	)
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
	groups, err := s.groupRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	stats, err := s.repo.GetGroupStats(ctx)
	if err != nil {
		slog.Warn("Failed to fetch group stats", "error", err)
		return groups, nil
	}

	statMap := make(map[int]int)
	for _, stat := range stats {
		statMap[stat.GroupID] = int(stat.Count)
	}

	for _, group := range groups {
		if count, ok := statMap[group.ID]; ok {
			group.CoinCount = count
		}
	}

	return groups, nil
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
		// var totalSilver, totalGold float64 // Removed unused vars

		stats.OldestCoin, _ = s.repo.GetOldestCoin(ctx)

		// Calculate Oldest High Grade Coin (>= EBC)
		// We already have allCoins, let's iterate.
		// We want the oldest year with grade rank >= EBC (rank 4)
		var oldestHighGrade *domain.Coin
		minHighGradeYear := 9999

		getGradeRank := func(grade string) int {
			switch normalizeGrade(grade) {
			case "FDC", "BU", "Proof":
				return 6
			case "SC", "UNC":
				return 5
			case "EBC", "XF", "AU":
				return 4
			case "MBC", "VF":
				return 3
			case "BC", "F":
				return 2
			case "RC", "VG", "G":
				return 1
			default:
				return 0
			}
		}

		for i := range allCoins {
			c := allCoins[i]
			// Copy to stats.AllCoins
			stats.AllCoins[i] = *c

			// Century
			if c.Year.Int() > 0 {
				century := (c.Year.Int() / 100) + 1
				stats.CenturyDistribution[fmt.Sprintf("%d", century)]++
			}

			// Decade (only last 200 years for clearer chart?) or all?
			if c.Year.Int() > 1800 {
				decade := (c.Year.Int() / 10) * 10
				stats.DecadeDistribution[fmt.Sprintf("%d", decade)]++
			}

			// Precious metals weight
			mat := strings.ToLower(c.Material)
			if strings.Contains(mat, "plata") || strings.Contains(mat, "silver") {
				stats.TotalSilverWeight += c.WeightG
			}
			if strings.Contains(mat, "oro") || strings.Contains(mat, "gold") {
				stats.TotalGoldWeight += c.WeightG
			}

			// Oldest High Grade Logic
			if c.Grade.String() != "" && getGradeRank(c.Grade.String()) >= 4 { // EBC or better
				if c.Year.Int() > 0 && c.Year.Int() < minHighGradeYear {
					minHighGradeYear = c.Year.Int()
					oldestHighGrade = c
				} else if c.Year.Int() == minHighGradeYear {
					// Tie-breaker? Maybe higher grade?
					if oldestHighGrade != nil && getGradeRank(c.Grade.String()) > getGradeRank(oldestHighGrade.Grade.String()) {
						oldestHighGrade = c
					}
				}
			}
		}
		if oldestHighGrade != nil {
			stats.OldestHighGradeCoin = oldestHighGrade
		}
	}

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
	coin.Mintage, _ = domain.NewMintage(params.Mintage)
	coin.Country = params.Country
	coin.Year, _ = domain.NewYear(params.Year)
	coin.FaceValue = params.FaceValue
	coin.Currency = params.Currency
	coin.Material = params.Material
	coin.Description = params.Description
	coin.KMCode, _ = domain.NewKMCode(params.KMCode)
	coin.MinValue = params.MinValue
	coin.MaxValue = params.MaxValue
	coin.Grade, _ = domain.NewGrade(normalizeGrade(params.Grade))
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
	// Delete from database first
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Then delete files from storage
	if err := s.storage.DeleteCoinDirectory(id); err != nil {
		// Log error but don't fail the deletion
		// The coin is already deleted from DB
		slog.Error("failed to delete coin files", "coin_id", id, "error", err)
	}

	return nil
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
			switch img.Side {
			case "front":
				frontPath = img.Path
			case "back":
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
	coin.Year, _ = domain.NewYear(analysis.Year)
	coin.FaceValue = analysis.FaceValue
	coin.Currency = analysis.Currency
	coin.Material = analysis.Material
	coin.Description = analysis.Description
	coin.KMCode, _ = domain.NewKMCode(analysis.KMCode)
	coin.NumistaNumber = analysis.NumistaNumber
	coin.MinValue = analysis.MinValue
	coin.MaxValue = analysis.MaxValue
	coin.Grade, _ = domain.NewGrade(normalizeGrade(analysis.Grade))
	coin.TechnicalNotes = analysis.Notes
	coin.GeminiDetails = analysis.RawDetails
	coin.Name = analysis.Name
	coin.Mint = analysis.Mint
	coin.Mintage, _ = domain.NewMintage(analysis.Mintage)
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

// MarkCoinAsSold marks a coin as sold
func (s *CoinService) MarkCoinAsSold(ctx context.Context, id uuid.UUID, soldPrice float64, saleChannel string) (*domain.Coin, error) {
	soldAt := time.Now()
	return s.repo.MarkAsSold(ctx, id, soldAt, soldPrice, saleChannel)
}

// GetSaleChannels returns list of distinct sale channels
func (s *CoinService) GetSaleChannels(ctx context.Context) ([]string, error) {
	return s.repo.GetSaleChannels(ctx)
}

func (s *CoinService) ExportCoinsCSV(ctx context.Context) ([]byte, error) {
	coins, err := s.repo.List(ctx, domain.CoinFilter{Limit: 1000000}) // Get all
	if err != nil {
		return nil, fmt.Errorf("failed to list coins: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Headers
	header := []string{
		"ID", "Name", "Country", "Year", "Face Value", "Currency", "Composition",
		"Grade", "Mintage", "Mint", "Weight (g)", "Diameter (mm)", "Acquired Date",
		"Price Paid", "Sold Date", "Sold Price", "Notes",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, c := range coins {
		acquired := ""
		if c.AcquiredAt != nil {
			acquired = c.AcquiredAt.Format("2006-01-02")
		}
		sold := ""
		if c.SoldAt != nil {
			sold = c.SoldAt.Format("2006-01-02")
		}
		record := []string{
			c.ID.String(),
			c.Name,
			c.Country,
			fmt.Sprintf("%d", c.Year.Int()),
			c.FaceValue,
			c.Currency,
			c.Material,
			c.Grade.String(),
			fmt.Sprintf("%d", c.Mintage.Int64()),
			c.Mint,
			fmt.Sprintf("%.2f", c.WeightG),
			fmt.Sprintf("%.2f", c.DiameterMM),
			acquired,
			fmt.Sprintf("%.2f", c.PricePaid),
			sold,
			fmt.Sprintf("%.2f", c.SoldPrice),
			c.PersonalNotes + " " + c.TechnicalNotes,
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *CoinService) ExportCoinsSQL(ctx context.Context) ([]byte, error) {
	coins, err := s.repo.List(ctx, domain.CoinFilter{Limit: 1000000})
	if err != nil {
		return nil, fmt.Errorf("failed to list coins: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("-- NumismaticApp Database Dump\n")
	sb.WriteString(fmt.Sprintf("-- Generated: %s\n\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("BEGIN;\n\n")

	// Helper to escape SQL strings
	escape := func(s string) string {
		return "'" + strings.ReplaceAll(s, "'", "''") + "'"
	}
	// Helper for NULLs
	sOrNull := func(s string) string {
		if s == "" {
			return "NULL"
		}
		return escape(s)
	}
	fOrNull := func(f float64) string {
		if f == 0 {
			return "NULL"
		}
		return fmt.Sprintf("%f", f)
	}
	iOrNull := func(i int64) string {
		return fmt.Sprintf("%d", i)
	}

	for _, c := range coins {
		vals := []string{
			escape(c.ID.String()),
			sOrNull(c.Name),
			sOrNull(c.Country),
			iOrNull(int64(c.Year.Int())), // Cast int to int64 for helper
			sOrNull(c.FaceValue),
			sOrNull(c.Currency),
			sOrNull(c.Material),
			sOrNull(c.Grade.String()),
			iOrNull(c.Mintage.Int64()),
			sOrNull(c.Mint),
			fOrNull(c.WeightG),
			fOrNull(c.DiameterMM),
			fOrNull(c.PricePaid),
			fOrNull(c.SoldPrice),
		}

		line := fmt.Sprintf("INSERT INTO coins (id, name, country, year, face_value, currency, material, grade, mintage, mint, weight_g, diameter_mm, price_paid, sold_price) VALUES (%s);\n",
			strings.Join(vals, ", "))
		sb.WriteString(line)
	}

	sb.WriteString("\nCOMMIT;\n")
	return []byte(sb.String()), nil
}
