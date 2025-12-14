package application_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antonioparicio/numismaticapp/internal/application"
	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddCoin_AICoverage(t *testing.T) {
	t.Run("AI Returns Invalid Year/Mintage", func(t *testing.T) {
		service, mockRepo, mockGroupRepo, mockImageService, mockAIService, mockStorage, mockBgRemover, mockNumistaClient := setupTest(t)
		ctx := context.Background()

		// Standard setup
		mockStorage.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("p", nil).AnyTimes()
		mockImageService.EXPECT().CropToContent(gomock.Any()).Return([]byte("c"), nil).AnyTimes()
		mockImageService.EXPECT().GenerateThumbnail(gomock.Any(), gomock.Any()).Return("t", nil).AnyTimes()
		mockBgRemover.EXPECT().RemoveBackground(gomock.Any(), gomock.Any()).Return([]byte("b"), nil).AnyTimes()
		mockGroupRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(&domain.Group{ID: 1}, nil).AnyTimes()
		mockImageService.EXPECT().GetMetadata(gomock.Any()).Return(100, 100, int64(10), "png", nil).AnyTimes()
		mockNumistaClient.EXPECT().SearchTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

		// AI returns invalid year (50000) and mintage (-1) to trigger fallback/warnings
		mockAIService.EXPECT().AnalyzeCoin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.CoinAnalysisResult{
			Year:    50000,
			Mintage: -1,
		}, nil)

		// Expect Save with defaulted values
		mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Coin) error {
			assert.Equal(t, 0, c.Year.Int())             // Defaulted
			assert.Equal(t, int64(0), c.Mintage.Int64()) // Defaulted/Zero value
			return nil
		})

		// Async update
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&domain.Coin{}, nil).AnyTimes()

		_, err := service.AddCoin(ctx, bytes.NewReader([]byte("f")), "f.jpg", bytes.NewReader([]byte("b")), "b.jpg", "G", "", "", "", 0, "m", 0)
		assert.NoError(t, err)
		time.Sleep(50 * time.Millisecond) // Wait for async
	})
}

func TestUpdateCoin_GroupCreateError(t *testing.T) {
	service, mockRepo, mockGroupRepo, _, _, _, _, _ := setupTest(t)
	ctx := context.Background()
	id := uuid.New()

	mockRepo.EXPECT().GetByID(ctx, id).Return(&domain.Coin{ID: id}, nil)

	// Group not found, try create, create fails
	mockGroupRepo.EXPECT().GetByName(ctx, "NewGroup").Return(nil, errors.New("not found"))
	mockGroupRepo.EXPECT().Create(ctx, "NewGroup", "").Return(nil, errors.New("create failed"))

	_, err := service.UpdateCoin(ctx, id, application.UpdateCoinParams{GroupName: "NewGroup"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create group")
}
