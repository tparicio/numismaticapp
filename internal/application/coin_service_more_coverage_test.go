package application_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/numista"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEnrichCoinWithNumista_MoreCoverage(t *testing.T) {
	coinID := uuid.New()
	t.Run("Candidate Detail Fetch Error", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient, _ := setupTest(t)
		ctx := context.Background()
		// Explicit expectations for originals to avoid generic matching issues
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID:        coinID,
			FaceValue: "10 USD",
			Year:      mustYear(2000),
		}, nil)

		// Search returns 2 candidates
		mockNumistaClient.EXPECT().SearchTypes(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{
			Count: 2,
			Types: []numista.NumistaType{
				{ID: 101, Title: "Error Coin", MinYear: 2000, MaxYear: 2000},
				{ID: 102, Title: "Success Coin", MinYear: 2000, MaxYear: 2000},
			},
		}, nil)

		// 1. First candidate fails GetType
		mockNumistaClient.EXPECT().GetType(ctx, 101).Return(nil, errors.New("api error"))

		// 2. Second candidate succeeds
		mockNumistaClient.EXPECT().GetType(ctx, 102).Return(map[string]any{
			"value": map[string]any{"numeric_value": 10.0},
		}, nil)

		// Update should happen with ID 102 (Perfect Match)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Equal(t, 102, c.NumistaNumber)
			return nil
		})

		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.NoError(t, err)
	})

	t.Run("Too Many Results", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{ID: coinID}, nil)

		mockNumistaClient.EXPECT().SearchTypes(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{Count: 50}, nil)

		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Equal(t, 0, c.NumistaNumber)
			return nil
		})

		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.NoError(t, err)
	})

	t.Run("Non-Numeric Face Value", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID:        coinID,
			FaceValue: "Unknown", // Non-numeric
			Year:      mustYear(2000),
		}, nil)

		// Search returns candidates
		mockNumistaClient.EXPECT().SearchTypes(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{
			Count: 1,
			Types: []numista.NumistaType{{ID: 101, MinYear: 2000, MaxYear: 2000}},
		}, nil)

		mockNumistaClient.EXPECT().GetType(ctx, 101).Return(map[string]any{"value": map[string]any{"numeric_value": 10.0}}, nil)

		// No update because value mismatch (0 vs 10)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Equal(t, 0, c.NumistaNumber)
			return nil
		})

		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.NoError(t, err)
	})
}

func TestAddCoin_GroupCreateError(t *testing.T) {
	t.Run("Group Create Failure", func(t *testing.T) {
		service, _, mockGroupRepo, mockImageService, mockAIService, mockStorage, mockBgRemover, _, _ := setupTest(t)
		ctx := context.Background()

		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).AnyTimes()
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), gomock.Any()).Return("t", nil).AnyTimes()
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("b"), nil).AnyTimes()
		mockImageService.EXPECT().GetMetadata(gomock.Any()).Return(100, 100, int64(10), "png", nil).AnyTimes()
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{Name: "C"}, nil).AnyTimes()

		mockGroupRepo.EXPECT().GetByName(gomock.Any(), "FailGroup").Return(nil, errors.New("not found"))
		mockGroupRepo.EXPECT().Create(gomock.Any(), "FailGroup", "").Return(nil, errors.New("create error"))

		_, err := service.AddCoin(ctx, bytes.NewReader([]byte("f")), "f.jpg", bytes.NewReader([]byte("b")), "b.jpg", "FailGroup", "", "", "", 0, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create group")
	})
}

func TestAddCoin_ImageProcErrors(t *testing.T) {
	t.Run("Crop Error Front", func(t *testing.T) {
		service, _, _, mockImageService, mockAIService, mockStorage, mockBgRemover, _, _ := setupTest(t)
		ctx := context.Background()
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("p"), nil).AnyTimes()
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{}, nil).AnyTimes()

		// Crop fails
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return(nil, errors.New("crop fail")).Times(1)

		_, err := service.AddCoin(ctx, bytes.NewReader([]byte("f")), "f.jpg", bytes.NewReader([]byte("b")), "b.jpg", "", "", "", "", 0, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to crop front")
	})

	t.Run("Save Processed Front Error", func(t *testing.T) {
		service, _, _, mockImageService, mockAIService, mockStorage, mockBgRemover, _, _ := setupTest(t)
		ctx := context.Background()
		// Note: SaveFile is called for original first (2 calls), then processed.
		// Use explicit filenames to distinguish or order.
		mockStorage.EXPECT().SaveFile(gomock.Any(), "original_front.jpg", gomock.Any()).Return("p", nil)
		mockStorage.EXPECT().SaveFile(gomock.Any(), "original_back.jpg", gomock.Any()).Return("p", nil)

		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("p"), nil).AnyTimes()
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).AnyTimes()
		// Processed Save Fails
		mockStorage.EXPECT().SaveFile(gomock.Any(), "processed_front.png", gomock.Any()).Return("", errors.New("save fail")).Times(1)

		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{}, nil).AnyTimes()

		_, err := service.AddCoin(ctx, bytes.NewReader([]byte("f")), "f.jpg", bytes.NewReader([]byte("b")), "b.jpg", "", "", "", "", 0, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save processed front")
	})

	t.Run("Generate Thumbnail Front Error", func(t *testing.T) {
		service, _, _, mockImageService, mockAIService, mockStorage, mockBgRemover, _, _ := setupTest(t)
		ctx := context.Background()
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("p"), nil).AnyTimes()
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).AnyTimes()
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{}, nil).AnyTimes()

		// Thumb Fails
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), 300).Return("", errors.New("thumb fail")).Times(1)

		_, err := service.AddCoin(ctx, bytes.NewReader([]byte("f")), "f.jpg", bytes.NewReader([]byte("b")), "b.jpg", "", "", "", "", 0, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to thumb front")
	})

	t.Run("Remove Background Back Error", func(t *testing.T) {
		service, _, _, mockImageService, mockAIService, mockStorage, mockBgRemover, _, _ := setupTest(t)
		ctx := context.Background()
		// Explicit expectations for originals to avoid generic matching issues
		mockStorage.EXPECT().SaveFile(gomock.Any(), "original_front.jpg", gomock.Any()).Return("p", nil)
		mockStorage.EXPECT().SaveFile(gomock.Any(), "original_back.jpg", gomock.Any()).Return("p", nil)

		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{}, nil).AnyTimes()

		// 1. BgRemove Front OK
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("p"), nil).Times(1) // Front
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).Times(1)               // Front
		mockStorage.EXPECT().SaveFile(gomock.Any(), "processed_front.png", gomock.Any()).Return("p", nil).Times(1)
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), 300).Return("t", nil).Times(1)

		// 2. BgRemove Back Fail
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return(nil, errors.New("bg back fail")).Times(1) // Back

		_, err := service.AddCoin(ctx, bytes.NewReader([]byte("f")), "f.jpg", bytes.NewReader([]byte("b")), "b.jpg", "", "", "", "", 0, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to bg remove back")
	})
}

func TestAddCoin_AIAnalysisError(t *testing.T) {
	t.Run("AI Returns Error", func(t *testing.T) {
		service, mockRepo, mockGroupRepo, mockImageService, mockAIService, mockStorage, mockBgRemover, mockNumistaClient, _ := setupTest(t)
		ctx := context.Background()

		// Standard setup for success except AI
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).AnyTimes()
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), gomock.Any()).Return("t", nil).AnyTimes()
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("b"), nil).AnyTimes()
		mockGroupRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(&domain.Group{ID: 1}, nil).AnyTimes()
		mockImageService.EXPECT().GetMetadata(gomock.Any()).Return(100, 100, int64(10), "png", nil).AnyTimes()
		mockNumistaClient.EXPECT().SearchTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

		// AI Fails
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("ai error"))

		// Expect Save with fallback (Description="Analysis failed")
		mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Contains(t, c.Description, "Analysis failed")
			// Details should contain error
			return nil
		})

		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&domain.Coin{}, nil).AnyTimes()

		_, err := service.AddCoin(ctx, bytes.NewReader([]byte("f")), "f.jpg", bytes.NewReader([]byte("b")), "b.jpg", "G", "", "", "", 0, "m", 0)
		assert.NoError(t, err)
		time.Sleep(50 * time.Millisecond) // Wait for async
	})
}

func TestReanalyzeCoin_MoreCoverage(t *testing.T) {
	coinID := uuid.New()
	t.Run("Missing Images", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		// Coin has NO original images
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID: coinID,
			Images: []domain.CoinImage{
				{Side: "front", ImageType: "crop", Path: "p.png"}, // Only crop
			},
		}, nil)

		_, err := service.ReanalyzeCoin(ctx, coinID, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "original images not found")
	})

	t.Run("Update Repo Error", func(t *testing.T) {
		service, mockRepo, _, _, mockAIService, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID: coinID,
			Images: []domain.CoinImage{
				{Side: "front", ImageType: "original", Path: "f.jpg"},
				{Side: "back", ImageType: "original", Path: "b.jpg"},
			},
		}, nil)

		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{}, nil)

		// Update Fail
		mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db fail"))

		_, err := service.ReanalyzeCoin(ctx, coinID, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update coin")
	})
}
