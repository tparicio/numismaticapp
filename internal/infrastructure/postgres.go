package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/antonioparicio/numismaticapp/internal/infrastructure/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCoinRepository struct {
	q  *db.Queries
	db *pgxpool.Pool
}

func NewPostgresCoinRepository(pool *pgxpool.Pool) *PostgresCoinRepository {
	return &PostgresCoinRepository{
		q:  db.New(pool),
		db: pool,
	}
}

func (r *PostgresCoinRepository) Save(ctx context.Context, coin *domain.Coin) error {
	// Convert map to JSONB byte array
	geminiDetailsBytes, err := json.Marshal(coin.GeminiDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal gemini details: %w", err)
	}

	params := db.CreateCoinParams{
		Country:             coin.Country,
		Year:                int32(coin.Year),
		FaceValue:           coin.FaceValue,
		Currency:            coin.Currency,
		Material:            coin.Material,
		Description:         toNullString(coin.Description),
		KmCode:              toNullString(coin.KMCode),
		MinValue:            toNumeric(coin.MinValue),
		MaxValue:            toNumeric(coin.MaxValue),
		Grade:               toNullString(coin.Grade),
		SampleImageUrlFront: toNullString(coin.SampleImageURLFront),
		SampleImageUrlBack:  toNullString(coin.SampleImageURLBack),
		Notes:               toNullString(coin.Notes),
		GeminiDetails:       geminiDetailsBytes,
	}

	result, err := r.q.CreateCoin(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save coin: %w", err)
	}

	coin.ID = uuid.UUID(result.ID.Bytes)
	coin.CreatedAt = result.CreatedAt.Time
	coin.UpdatedAt = result.UpdatedAt.Time

	return nil
}

func (r *PostgresCoinRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Coin, error) {
	row, err := r.q.GetCoin(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}
	return toDomainCoin(row)
}

func (r *PostgresCoinRepository) List(ctx context.Context, limit, offset int) ([]*domain.Coin, error) {
	rows, err := r.q.ListCoins(ctx, db.ListCoinsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list coins: %w", err)
	}

	coins := make([]*domain.Coin, len(rows))
	for i, row := range rows {
		c, err := toDomainCoin(row)
		if err != nil {
			return nil, err
		}
		coins[i] = c
	}
	return coins, nil
}

func (r *PostgresCoinRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountCoins(ctx)
}

// Helper functions for conversion

func toDomainCoin(row db.Coin) (*domain.Coin, error) {
	var geminiDetails map[string]any
	if len(row.GeminiDetails) > 0 {
		if err := json.Unmarshal(row.GeminiDetails, &geminiDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal gemini details: %w", err)
		}
	}

	// Helper to convert numeric to float64 (simplified for this example)
	// In production, handle pgtype.Numeric carefully
	minVal, _ := row.MinValue.Float64Value()
	maxVal, _ := row.MaxValue.Float64Value()

	return &domain.Coin{
		ID:                  uuid.UUID(row.ID.Bytes),
		Country:             row.Country,
		Year:                int(row.Year),
		FaceValue:           row.FaceValue,
		Currency:            row.Currency,
		Material:            row.Material,
		Description:         row.Description.String,
		KMCode:              row.KmCode.String,
		MinValue:            minVal.Float64,
		MaxValue:            maxVal.Float64,
		Grade:               row.Grade.String,
		SampleImageURLFront: row.SampleImageUrlFront.String,
		SampleImageURLBack:  row.SampleImageUrlBack.String,
		Notes:               row.Notes.String,
		GeminiDetails:       geminiDetails,
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}, nil
}

func toNullString(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

func toNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(fmt.Sprintf("%f", f))
	return n
}
