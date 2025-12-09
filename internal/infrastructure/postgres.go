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
		Country:             toNullString(coin.Country),
		Year:                toNullInt4(coin.Year),
		FaceValue:           toNullString(coin.FaceValue),
		Currency:            toNullString(coin.Currency),
		Material:            toNullString(coin.Material),
		Description:         toNullString(coin.Description),
		KmCode:              toNullString(coin.KMCode),
		MinValue:            toNumeric(coin.MinValue),
		MaxValue:            toNumeric(coin.MaxValue),
		Grade:               toNullGradeType(coin.Grade),
		SampleImageUrlFront: toNullString(coin.SampleImageURLFront),
		SampleImageUrlBack:  toNullString(coin.SampleImageURLBack),
		Notes:               toNullString(coin.Notes),
		GeminiDetails:       geminiDetailsBytes,
		GroupID:             toNullInt4Ptr(coin.GroupID),
		UserNotes:           toNullString(coin.UserNotes),
	}

	result, err := r.q.CreateCoin(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save coin: %w", err)
	}

	coin.ID = uuid.UUID(result.ID.Bytes)
	coin.CreatedAt = result.CreatedAt.Time
	coin.UpdatedAt = result.UpdatedAt.Time

	// Save Images
	for _, img := range coin.Images {
		imgParams := db.CreateCoinImageParams{
			CoinID:    pgtype.UUID{Bytes: coin.ID, Valid: true},
			ImageType: db.ImageType(img.ImageType),
			Side:      db.CoinSide(img.Side),
			Path:      img.Path,
			Extension: img.Extension,
			Size:      img.Size,
			Width:     int32(img.Width),
			Height:    int32(img.Height),
			MimeType:  img.MimeType,
		}
		if _, err := r.q.CreateCoinImage(ctx, imgParams); err != nil {
			return fmt.Errorf("failed to save coin image: %w", err)
		}
	}

	return nil
}

func (r *PostgresCoinRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Coin, error) {
	row, err := r.q.GetCoin(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}

	coin, err := toDomainCoin(row)
	if err != nil {
		return nil, err
	}

	// Fetch images
	images, err := r.q.ListCoinImagesByCoinID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list coin images: %w", err)
	}

	coin.Images = toDomainImages(images)

	return coin, nil
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

		// Ideally we should batch fetch images or use a join, but for now N+1 is acceptable for MVP
		images, err := r.q.ListCoinImagesByCoinID(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list coin images for coin %v: %w", row.ID, err)
		}
		c.Images = toDomainImages(images)

		coins[i] = c
	}
	return coins, nil
}

func (r *PostgresCoinRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountCoins(ctx)
}

// Group Repository Implementation

type PostgresGroupRepository struct {
	q *db.Queries
}

func NewPostgresGroupRepository(pool *pgxpool.Pool) *PostgresGroupRepository {
	return &PostgresGroupRepository{
		q: db.New(pool),
	}
}

func (r *PostgresGroupRepository) Create(ctx context.Context, name, description string) (*domain.Group, error) {
	row, err := r.q.CreateGroup(ctx, db.CreateGroupParams{
		Name:        name,
		Description: toNullString(description),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}
	return toDomainGroup(row), nil
}

func (r *PostgresGroupRepository) GetByName(ctx context.Context, name string) (*domain.Group, error) {
	row, err := r.q.GetGroupByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get group by name: %w", err)
	}
	return toDomainGroup(row), nil
}

func (r *PostgresGroupRepository) List(ctx context.Context) ([]*domain.Group, error) {
	rows, err := r.q.ListGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	groups := make([]*domain.Group, len(rows))
	for i, row := range rows {
		groups[i] = toDomainGroup(row)
	}
	return groups, nil
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

	var groupID *int
	if row.GroupID.Valid {
		id := int(row.GroupID.Int32)
		groupID = &id
	}

	return &domain.Coin{
		ID:                  uuid.UUID(row.ID.Bytes),
		Country:             row.Country.String,
		Year:                int(row.Year.Int32),
		FaceValue:           row.FaceValue.String,
		Currency:            row.Currency.String,
		Material:            row.Material.String,
		Description:         row.Description.String,
		KMCode:              row.KmCode.String,
		MinValue:            minVal.Float64,
		MaxValue:            maxVal.Float64,
		Grade:               string(row.Grade.GradeType),
		SampleImageURLFront: row.SampleImageUrlFront.String,
		SampleImageURLBack:  row.SampleImageUrlBack.String,
		Notes:               row.Notes.String,
		GeminiDetails:       geminiDetails,
		GroupID:             groupID,
		UserNotes:           row.UserNotes.String,
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}, nil
}

func toDomainImages(rows []db.CoinImage) []domain.CoinImage {
	images := make([]domain.CoinImage, len(rows))
	for i, row := range rows {
		images[i] = domain.CoinImage{
			ID:        uuid.UUID(row.ID.Bytes),
			CoinID:    uuid.UUID(row.CoinID.Bytes),
			ImageType: string(row.ImageType),
			Side:      string(row.Side),
			Path:      row.Path,
			Extension: row.Extension,
			Size:      row.Size,
			Width:     int(row.Width),
			Height:    int(row.Height),
			MimeType:  row.MimeType,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}
	return images
}

func toDomainGroup(row db.Group) *domain.Group {
	return &domain.Group{
		ID:          int(row.ID),
		Name:        row.Name,
		Description: row.Description.String,
		CreatedAt:   row.CreatedAt.Time,
	}
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

func toNullInt4(i int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(i),
		Valid: i != 0, // Assuming 0 means null/unset for int fields in this domain
	}
}

func toNullInt4Ptr(i *int) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{
		Int32: int32(*i),
		Valid: true,
	}
}

func toNullGradeType(s string) db.NullGradeType {
	return db.NullGradeType{
		GradeType: db.GradeType(s),
		Valid:     s != "",
	}
}
