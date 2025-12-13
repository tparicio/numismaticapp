package application_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/antonioparicio/numismaticapp/internal/application"
	"github.com/antonioparicio/numismaticapp/internal/application/mocks"
	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/numista"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (
	*application.CoinService,
	*mocks.MockCoinRepository,
	*mocks.MockGroupRepository,
	*mocks.MockImageService,
	*mocks.MockAIService,
	*mocks.MockStorageService,
	*mocks.MockBackgroundRemover,
	*mocks.MockNumistaService,
) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCoinRepository(ctrl)
	mockGroupRepo := mocks.NewMockGroupRepository(ctrl)
	mockImageService := mocks.NewMockImageService(ctrl)
	mockAIService := mocks.NewMockAIService(ctrl)
	mockStorage := mocks.NewMockStorageService(ctrl)
	mockBgRemover := mocks.NewMockBackgroundRemover(ctrl)
	mockNumistaClient := mocks.NewMockNumistaService(ctrl)

	service := application.NewCoinService(
		mockRepo,
		mockGroupRepo,
		mockImageService,
		mockAIService,
		mockStorage,
		mockBgRemover,
		mockNumistaClient,
	)

	return service, mockRepo, mockGroupRepo, mockImageService, mockAIService, mockStorage, mockBgRemover, mockNumistaClient
}

func TestListCoins(t *testing.T) {
	t.Run("Basic List", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		expectedCoins := []*domain.Coin{
			{ID: uuid.New(), Name: "Coin 1"},
			{ID: uuid.New(), Name: "Coin 2"},
		}
		mockRepo.EXPECT().List(ctx, domain.CoinFilter{}).Return(expectedCoins, nil)

		coins, err := service.ListCoins(ctx, domain.CoinFilter{})
		assert.NoError(t, err)
		assert.Equal(t, expectedCoins, coins)
	})

	t.Run("Filtered List (Year, Query)", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		year := 2000
		query := "Euro"
		filter := domain.CoinFilter{Year: &year, Query: &query}
		mockRepo.EXPECT().List(ctx, filter).Return([]*domain.Coin{}, nil)
		coins, err := service.ListCoins(ctx, filter)
		assert.NoError(t, err)
		assert.Empty(t, coins)
	})

	t.Run("Filtered List (Country)", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		country := "Spain"
		filter := domain.CoinFilter{Country: &country}
		mockRepo.EXPECT().List(ctx, filter).Return([]*domain.Coin{}, nil)
		coins, err := service.ListCoins(ctx, filter)
		assert.NoError(t, err)
		assert.Empty(t, coins)
	})
}

func TestGetCoin(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		id := uuid.New()
		expectedCoin := &domain.Coin{ID: id, Name: "Test Coin"}
		mockRepo.EXPECT().GetByID(ctx, id).Return(expectedCoin, nil)
		coin, err := service.GetCoin(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expectedCoin, coin)
	})

	t.Run("Error", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		id := uuid.New()
		mockRepo.EXPECT().GetByID(ctx, id).Return(nil, assert.AnError)
		_, err := service.GetCoin(ctx, id)
		assert.Error(t, err)
	})
}

func TestGetDashboardStats(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()

	mockRepo.EXPECT().Count(ctx).Return(int64(10), nil)
	mockRepo.EXPECT().GetTotalValue(ctx).Return(100.0, nil)
	mockRepo.EXPECT().GetAverageValue(ctx).Return(10.0, nil)
	mockRepo.EXPECT().ListTopValuable(ctx).Return([]*domain.Coin{}, nil)
	mockRepo.EXPECT().ListRecent(ctx).Return([]*domain.Coin{}, nil)
	mockRepo.EXPECT().GetMaterialDistribution(ctx).Return(map[string]int{"Gold": 5}, nil)
	mockRepo.EXPECT().GetGradeDistribution(ctx).Return(map[string]int{"XF": 5}, nil)
	mockRepo.EXPECT().GetAllValues(ctx).Return([]float64{10, 20}, nil)
	mockRepo.EXPECT().GetCountryDistribution(ctx).Return(map[string]int{"Spain": 5}, nil)

	mockRepo.EXPECT().GetAllCoins(ctx).Return([]*domain.Coin{
		{Year: 2000, Material: "Gold", WeightG: 10},
	}, nil)

	mockRepo.EXPECT().GetOldestCoin(ctx).Return(&domain.Coin{Year: 1800}, nil)
	mockRepo.EXPECT().GetRarestCoins(ctx, 5).Return([]*domain.Coin{}, nil)
	mockRepo.EXPECT().GetGroupDistribution(ctx).Return(map[string]int{"Group 1": 1}, nil)
	mockRepo.EXPECT().GetGroupStats(ctx).Return([]domain.GroupStat{}, nil)
	mockRepo.EXPECT().GetHeaviestCoin(ctx).Return(&domain.Coin{}, nil)
	mockRepo.EXPECT().GetSmallestCoin(ctx).Return(&domain.Coin{}, nil)
	mockRepo.EXPECT().GetRandomCoin(ctx).Return(&domain.Coin{}, nil)

	stats, err := service.GetDashboardStats(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(10), stats.TotalCoins)
}

func TestAddCoin_Flows(t *testing.T) {
	frontData := []byte("f")
	backData := []byte("b")

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, mockGroupRepo, mockImageService, mockAIService, mockStorage, mockBgRemover, mockNumistaClient := setupTest(t)
		ctx := context.Background()
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p1", nil).Times(2)
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{Name: "C"}, nil)
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("b"), nil).Times(2)
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).Times(2)
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p2", nil).Times(2)
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), 300).Return("t", nil).Times(2)
		mockImageService.EXPECT().GetMetadata(gomock.Any()).Return(100, 100, int64(100), "image/png", nil).AnyTimes()
		mockGroupRepo.EXPECT().GetByName(gomock.Any(), "G").Return(&domain.Group{ID: 1}, nil)
		mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		// Async Numista might call Update or GetByID, allow it
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&domain.Coin{Year: 2024, FaceValue: "1"}, nil).AnyTimes()
		mockNumistaClient.EXPECT().SearchTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{}, nil).AnyTimes()

		_, err := service.AddCoin(ctx, bytes.NewReader(frontData), "f.jpg", bytes.NewReader(backData), "b.jpg", "G", "", "", "", 0, "m", 0)
		assert.NoError(t, err)
		// Wait slightly for async to potentially run? Not strict.
	})

	t.Run("Storage Err", func(t *testing.T) {
		service, _, _, _, _, mockStorage, _, _ := setupTest(t)
		ctx := context.Background()
		// Fail first save
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("", assert.AnError).MaxTimes(1)
		// Other calls might happen or not depending on race, allow them
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()

		_, err := service.AddCoin(ctx, bytes.NewReader(frontData), "f.jpg", bytes.NewReader(backData), "b.jpg", "G", "", "", "", 0, "m", 0)
		assert.Error(t, err)
	})

	t.Run("AI Err", func(t *testing.T) {

		service, mockRepo, mockGroupRepo, mockImageService, mockAIService, mockStorage, mockBgRemover, mockNumistaClient := setupTest(t)
		ctx := context.Background()
		// Storage succeeds
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()
		// AI Fails locally
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		// Parallel tasks B and C might start, allow their calls
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("b"), nil).AnyTimes()
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).AnyTimes()
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), 300).Return("t", nil).AnyTimes()
		mockImageService.EXPECT().GetMetadata(gomock.Any()).Return(100, 100, int64(100), "image/png", nil).AnyTimes()
		mockGroupRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(&domain.Group{ID: 1}, nil).AnyTimes()
		mockGroupRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Group{ID: 1}, nil).AnyTimes()

		// AddCoin proceeds on AI error (soft fail), so it Saves
		mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		// It might trigger async enrichment if fields are present (unlikely from nil AI result but plausible flow)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&domain.Coin{}, nil).AnyTimes()
		mockNumistaClient.EXPECT().SearchTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{}, nil).AnyTimes()

		coin, err := service.AddCoin(ctx, bytes.NewReader(frontData), "f.jpg", bytes.NewReader(backData), "b.jpg", "G", "", "", "", 0, "m", 0)
		assert.NoError(t, err)
		assert.NotNil(t, coin)
	})
}

func TestEnrichCoinWithNumista(t *testing.T) {
	coinID := uuid.New()

	t.Run("Match Found", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()
		coin := &domain.Coin{ID: coinID, FaceValue: "20 Euro Cent", Year: 2008}
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(coin, nil)

		mockNumistaClient.EXPECT().SearchTypes(ctx, gomock.Any(), gomock.Any(), "2008", gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{
			Count: 1,
			Types: []numista.NumistaType{
				{ID: 123, Title: "20 Cents", MinYear: 2007, MaxYear: 2009},
			},
		}, nil)

		mockNumistaClient.EXPECT().GetType(ctx, 123).Return(map[string]any{
			"value": map[string]any{"numeric_value": 0.2},
			"shape": "Round",
		}, nil)

		mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.NoError(t, err)
	})

	t.Run("No Match Found", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()
		coin := &domain.Coin{ID: coinID, FaceValue: "Rare Coin", Year: 1900}
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(coin, nil)
		mockNumistaClient.EXPECT().SearchTypes(ctx, gomock.Any(), gomock.Any(), "1900", gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{Count: 0}, nil)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.NoError(t, err)
	})

	t.Run("Repo Error Get", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(nil, assert.AnError)
		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.Error(t, err)
	})

	t.Run("Search Error", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{FaceValue: "V", Year: 2000}, nil)
		mockNumistaClient.EXPECT().SearchTypes(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.Error(t, err)
	})
}

func TestApplyNumistaCandidate(t *testing.T) {
	coinID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{ID: coinID}, nil)
		mockNumistaClient.EXPECT().GetType(ctx, 999).Return(map[string]any{"title": "Manual Selection"}, nil)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
		_, err := service.ApplyNumistaCandidate(ctx, coinID, 999)
		assert.NoError(t, err)
	})

	t.Run("Error GetType", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{}, nil)
		mockNumistaClient.EXPECT().GetType(ctx, 999).Return(nil, assert.AnError)
		_, err := service.ApplyNumistaCandidate(ctx, coinID, 999)
		assert.Error(t, err)
	})

	t.Run("Full Mapping", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()

		fullDetails := map[string]any{
			"title":       "Full Coin",
			"size":        25.0,
			"thickness":   2.0,
			"weight":      8.5,
			"shape":       "Round",
			"composition": map[string]any{"text": "Gold"},
			"mints": []any{
				map[string]any{"name": "Royal Mint"},
			},
			"references": []any{
				map[string]any{
					"catalogue": map[string]any{"code": "KM"},
					"number":    "123",
				},
			},
			"ruler": []any{
				map[string]any{"name": "King Charles"},
			},
			"orientation":        "Coin alignment",
			"series":             "Commemorative",
			"commemorated_topic": "Anniversary",
		}

		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{ID: coinID}, nil)
		mockNumistaClient.EXPECT().GetType(ctx, 999).Return(fullDetails, nil)

		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Equal(t, 25.0, c.DiameterMM)
			assert.Equal(t, 2.0, c.ThicknessMM)
			assert.Equal(t, 8.5, c.WeightG)
			assert.Equal(t, "Round", c.Shape)
			assert.Equal(t, "Gold", c.Material)
			assert.Equal(t, "Royal Mint", c.Mint)
			assert.Equal(t, "KM# 123", c.KMCode)
			assert.Equal(t, "King Charles", c.Ruler)
			assert.Equal(t, "Coin alignment", c.Orientation)
			assert.Equal(t, "Commemorative", c.Series)
			assert.Equal(t, "Anniversary", c.CommemoratedTopic)
			return nil
		})

		_, err := service.ApplyNumistaCandidate(ctx, coinID, 999)
		assert.NoError(t, err)
	})
}

func TestListGroups(t *testing.T) {
	service, _, mockGroupRepo, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()
	mockGroupRepo.EXPECT().List(ctx).Return([]*domain.Group{{Name: "G1"}}, nil)
	groups, err := service.ListGroups(ctx)
	assert.NoError(t, err)
	assert.Len(t, groups, 1)
}

func TestCreateGroup(t *testing.T) {
	service, _, mockGroupRepo, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()
	mockGroupRepo.EXPECT().Create(ctx, "G1", "Desc").Return(&domain.Group{ID: 1}, nil)
	g, err := service.CreateGroup(ctx, "G1", "Desc")
	assert.NoError(t, err)
	assert.Equal(t, 1, g.ID)
}

func TestUpdateGroup(t *testing.T) {
	service, _, mockGroupRepo, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()
	mockGroupRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	_, err := service.UpdateGroup(ctx, 1, "G2", "Desc2")
	assert.NoError(t, err)
}

func TestDeleteGroup(t *testing.T) {
	service, _, mockGroupRepo, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()
	mockGroupRepo.EXPECT().Delete(ctx, 1).Return(nil)
	err := service.DeleteGroup(ctx, 1)
	assert.NoError(t, err)
}

func TestUpdateCoin(t *testing.T) {
	id := uuid.New()
	params := application.UpdateCoinParams{Name: "Updated Name", GroupName: "New Group"}

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, mockGroupRepo, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, id).Return(&domain.Coin{ID: id}, nil)
		mockGroupRepo.EXPECT().GetByName(ctx, "New Group").Return(&domain.Group{ID: 2}, nil)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
		_, err := service.UpdateCoin(ctx, id, params)
		assert.NoError(t, err)
	})

	t.Run("Get Error", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, id).Return(nil, assert.AnError)
		_, err := service.UpdateCoin(ctx, id, params)
		assert.Error(t, err)
	})
}

func TestDeleteCoin(t *testing.T) {
	id := uuid.New()
	t.Run("Success", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().Delete(ctx, id).Return(nil)
		err := service.DeleteCoin(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().Delete(ctx, id).Return(assert.AnError)
		err := service.DeleteCoin(ctx, id)
		assert.Error(t, err)
	})
}

func TestRotateCoinImage(t *testing.T) {
	coinID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, _, mockImageService, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID:     coinID,
			Images: []domain.CoinImage{{Side: "front", Extension: ".png", ImageType: "crop", Path: "path/front.png"}},
		}, nil)

		mockImageService.EXPECT().Rotate("path/front.png", 90.0).Return(nil)
		mockImageService.EXPECT().GenerateThumbnail("path/front.png", 300).Return("path/thumb.png", nil)
		// No strict update check as it might be implicit

		err := service.RotateCoinImage(ctx, coinID, "front", 90.0)
		assert.NoError(t, err)
	})

	t.Run("Repo Get Error", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(nil, assert.AnError)
		err := service.RotateCoinImage(ctx, coinID, "front", 90.0)
		assert.Error(t, err)
	})
}

func TestReanalyzeCoin(t *testing.T) {
	coinID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, _, _, mockAIService, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID: coinID,
			Images: []domain.CoinImage{
				{Side: "front", ImageType: "original", Path: "front.jpg"},
				{Side: "back", ImageType: "original", Path: "back.jpg"},
			},
		}, nil)

		mockAIService.EXPECT().AnalyzeCoin(ctx, "front.jpg", "back.jpg", "gemini-pro", float32(0.1), "es").Return(&domain.CoinAnalysisResult{
			Name: "Reanalyzed",
		}, nil)

		mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		coin, err := service.ReanalyzeCoin(ctx, coinID, "gemini-pro", 0.1)
		assert.NoError(t, err)
		assert.Equal(t, "Reanalyzed", coin.Name)
	})

	t.Run("AI Error", func(t *testing.T) {
		service, mockRepo, _, _, mockAIService, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			Images: []domain.CoinImage{{Side: "front", ImageType: "original", Path: "p1"}, {Side: "back", ImageType: "original", Path: "p2"}},
		}, nil)
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		_, err := service.ReanalyzeCoin(ctx, coinID, "m", 0)
		assert.Error(t, err)
	})
}

func TestGetGeminiModels(t *testing.T) {
	service, _, _, _, mockAIService, _, _, _ := setupTest(t)
	ctx := context.Background()
	mockAIService.EXPECT().ListModels(ctx).Return([]domain.GeminiModelInfo{{Name: "gemini-pro"}}, nil)
	models, err := service.GetGeminiModels(ctx)
	assert.NoError(t, err)
	assert.Len(t, models, 1)
}

func TestUpdateCoin_Grades(t *testing.T) {
	t.Run("Normalize Grades", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		id := uuid.New()

		testCases := []struct {
			input    string
			expected string
		}{
			{"EBC", "EBC"},
			{"ebc", "EBC"},
			{"MBC (Muy Bien Conservada)", "MBC"},
			{"Unknown Grade", ""},
			{"Rare Coin (FDC)", "FDC"},
		}

		for _, tc := range testCases {
			mockRepo.EXPECT().GetByID(ctx, id).Return(&domain.Coin{ID: id}, nil)
			mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
				assert.Equal(t, tc.expected, c.Grade)
				return nil
			})

			params := application.UpdateCoinParams{
				Grade: tc.input,
			}
			_, err := service.UpdateCoin(ctx, id, params)
			assert.NoError(t, err)
		}
	})
}

func TestGetDashboardStats_Century(t *testing.T) {
	// Tests the fallback in toRoman for centuries > 21
	service, mockRepo, _, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()

	mockRepo.EXPECT().Count(ctx).Return(int64(1), nil)
	mockRepo.EXPECT().GetTotalValue(ctx).Return(100.0, nil)
	mockRepo.EXPECT().GetAverageValue(ctx).Return(10.0, nil)
	mockRepo.EXPECT().ListTopValuable(ctx).Return([]*domain.Coin{}, nil)
	mockRepo.EXPECT().ListRecent(ctx).Return([]*domain.Coin{}, nil)
	mockRepo.EXPECT().GetMaterialDistribution(ctx).Return(map[string]int{}, nil)
	mockRepo.EXPECT().GetGradeDistribution(ctx).Return(map[string]int{}, nil)
	mockRepo.EXPECT().GetAllValues(ctx).Return([]float64{}, nil)
	mockRepo.EXPECT().GetCountryDistribution(ctx).Return(map[string]int{}, nil)
	mockRepo.EXPECT().GetOldestCoin(ctx).Return(&domain.Coin{Year: 1800}, nil)
	mockRepo.EXPECT().GetRarestCoins(ctx, 5).Return([]*domain.Coin{}, nil)
	mockRepo.EXPECT().GetGroupDistribution(ctx).Return(map[string]int{}, nil)
	mockRepo.EXPECT().GetGroupStats(ctx).Return([]domain.GroupStat{}, nil)
	mockRepo.EXPECT().GetHeaviestCoin(ctx).Return(&domain.Coin{}, nil)
	mockRepo.EXPECT().GetSmallestCoin(ctx).Return(&domain.Coin{}, nil)
	mockRepo.EXPECT().GetRandomCoin(ctx).Return(&domain.Coin{}, nil)

	// Inject a futuristic coin to trigger toRoman fallback (22nd century)
	mockRepo.EXPECT().GetAllCoins(ctx).Return([]*domain.Coin{
		{Year: 2199, Material: "Gold", WeightG: 10}, // 22nd Century -> "22"
	}, nil)

	stats, err := service.GetDashboardStats(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	// We don't verify the specific Century string as it's private and deep in stats,
	// but running this ensures coverage of that line.
}

func TestRotateCoinImage_Errors(t *testing.T) {
	t.Run("Thumbnail Error", func(t *testing.T) {
		service, mockRepo, _, mockImageService, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		coinID := uuid.New()

		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID:     coinID,
			Images: []domain.CoinImage{{Side: "front", Extension: ".png", ImageType: "crop", Path: "path/front.png"}},
		}, nil)

		mockImageService.EXPECT().Rotate("path/front.png", 90.0).Return(nil)
		mockImageService.EXPECT().GenerateThumbnail("path/front.png", 300).Return("", errors.New("thumb error"))

		err := service.RotateCoinImage(ctx, coinID, "front", 90.0)
		assert.Error(t, err)
	})

	t.Run("Image Not Found", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		coinID := uuid.New()

		mockRepo.EXPECT().GetByID(ctx, coinID).Return(&domain.Coin{
			ID:     coinID,
			Images: []domain.CoinImage{}, // No images
		}, nil)

		err := service.RotateCoinImage(ctx, coinID, "front", 90.0)
		assert.Error(t, err)
	})
}

func TestCreateGroup_Errors(t *testing.T) {
	t.Run("Repo Error", func(t *testing.T) {
		service, _, mockGroupRepo, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		mockGroupRepo.EXPECT().Create(ctx, "G1", "Desc").Return(nil, errors.New("db error"))
		_, err := service.CreateGroup(ctx, "G1", "Desc")
		assert.Error(t, err)
	})
	t.Run("Validation Error", func(t *testing.T) {
		service, _, _, _, _, _, _, _ := setupTest(t)
		ctx := context.Background()
		_, err := service.CreateGroup(ctx, "", "Desc")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group name cannot be empty")
	})
}

func TestAddCoin_ErrorPaths(t *testing.T) {
	frontData := []byte("f")
	backData := []byte("b")

	t.Run("Group Create Error", func(t *testing.T) {
		service, _, mockGroupRepo, mockImageService, mockAIService, mockStorage, mockBgRemover, _ := setupTest(t)
		ctx := context.Background()

		// Storage succeeds
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()

		// AI succeeds
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{Name: "C"}, nil)

		// Image Processing succeeds
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("b"), nil).AnyTimes()
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).AnyTimes()
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), 300).Return("t", nil).AnyTimes()
		mockImageService.EXPECT().GetMetadata(gomock.Any()).Return(100, 100, int64(100), "image/png", nil).AnyTimes()

		// Group Fails
		// GetByName returns error (triggers create)
		mockGroupRepo.EXPECT().GetByName(gomock.Any(), "G_Fail").Return(nil, errors.New("not found"))
		// Create returns error
		mockGroupRepo.EXPECT().Create(gomock.Any(), "G_Fail", gomock.Any()).Return(nil, errors.New("create error"))

		// Because Group failed, AddCoin should return error
		// Note: Clean up (DeleteCoin) might be called if implemented, or it just errors out.
		// The current implementation returns error on first error from channels.

		_, err := service.AddCoin(ctx, bytes.NewReader(frontData), "f.jpg", bytes.NewReader(backData), "b.jpg", "G_Fail", "", "", "", 0, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create group")
	})

	t.Run("Image Process Error", func(t *testing.T) {
		service, _, _, _, mockAIService, mockStorage, mockBgRemover, _ := setupTest(t)
		ctx := context.Background()

		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{Name: "C"}, nil)

		// Image Processing Fails at bg removal
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return(nil, errors.New("bg error")).Times(1)
		// Note: RemoveBackground is called twice (front/back). If first fails, it returns error.

		_, err := service.AddCoin(ctx, bytes.NewReader(frontData), "f.jpg", bytes.NewReader(backData), "b.jpg", "", "", "", "", 0, "m", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to bg remove")
	})
}

func TestUpdateGroup_Validation(t *testing.T) {
	service, _, _, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()
	_, err := service.UpdateGroup(ctx, 1, "", "Desc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "group name cannot be empty")
}

func TestEnrichCoinWithNumista_Complex(t *testing.T) {
	coinID := uuid.New()
	t.Run("Fallback to Value Match when Year Mismatches", func(t *testing.T) {
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()

		// Coin: 2005, Value 1 (Unit)
		coin := &domain.Coin{ID: coinID, FaceValue: "1 Euro", Year: 2005, Currency: "Euro", Country: "Spain"}
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(coin, nil)

		// Search returns:
		// 1. Year Mismatch (2000-2002), Value Match (1)
		// 2. Year Match (2004-2006), Value Mismatch (2)
		mockNumistaClient.EXPECT().SearchTypes(ctx, "1 Euro Euro Spain", "coin", "2005", "", 10).Return(&numista.TypeSearchResponse{
			Count: 2,
			Types: []numista.NumistaType{
				{ID: 101, Title: "1 Euro (Old)", MinYear: 2000, MaxYear: 2002},
				{ID: 102, Title: "2 Euro", MinYear: 2004, MaxYear: 2006},
			},
		}, nil)

		// Helper to return details
		mockNumistaClient.EXPECT().GetType(ctx, 101).Return(map[string]any{
			"value": map[string]any{"numeric_value": 1.0},
		}, nil)
		mockNumistaClient.EXPECT().GetType(ctx, 102).Return(map[string]any{
			"value": map[string]any{"numeric_value": 2.0},
		}, nil)

		// Expect update with ID 101 (Fallback) because it matched value even if year mismatched, and no perfect match was found.
		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Equal(t, 101, c.NumistaNumber)
			return nil
		})

		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.NoError(t, err)
	})

	t.Run("Value Unit Conversion Match", func(t *testing.T) {
		// Test 20 Cents vs 0.2
		service, mockRepo, _, _, _, _, _, mockNumistaClient := setupTest(t)
		ctx := context.Background()

		coin := &domain.Coin{ID: coinID, FaceValue: "20", Year: 2000}
		mockRepo.EXPECT().GetByID(ctx, coinID).Return(coin, nil)

		mockNumistaClient.EXPECT().SearchTypes(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&numista.TypeSearchResponse{
			Count: 1,
			Types: []numista.NumistaType{{ID: 201, MinYear: 2000, MaxYear: 2000}},
		}, nil)

		mockNumistaClient.EXPECT().GetType(ctx, 201).Return(map[string]any{
			"value": map[string]any{"numeric_value": 0.20}, // 0.20 unit matches 20 face value
		}, nil)

		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Equal(t, 201, c.NumistaNumber)
			return nil
		})

		err := service.EnrichCoinWithNumista(ctx, coinID)
		assert.NoError(t, err)
	})
}
