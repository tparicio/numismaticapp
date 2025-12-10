package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
		ID:             pgtype.UUID{Bytes: coin.ID, Valid: true},
		Name:           toNullString(coin.Name),
		Mint:           toNullString(coin.Mint),
		Mintage:        toNullInt8(coin.Mintage),
		Country:        toNullString(coin.Country),
		Year:           toNullInt4(coin.Year),
		FaceValue:      toNullString(coin.FaceValue),
		Currency:       toNullString(coin.Currency),
		Material:       toNullString(coin.Material),
		Description:    toNullString(coin.Description),
		KmCode:         toNullString(coin.KMCode),
		MinValue:       toNumeric(coin.MinValue),
		MaxValue:       toNumeric(coin.MaxValue),
		Grade:          toNullGradeType(coin.Grade),
		TechnicalNotes: toNullString(coin.TechnicalNotes),
		GeminiDetails:  geminiDetailsBytes,
		GroupID:        toNullInt4Ptr(coin.GroupID),
		PersonalNotes:  toNullString(coin.PersonalNotes),
		WeightG:        toNumeric(coin.WeightG),
		DiameterMm:     toNumeric(coin.DiameterMM),
		ThicknessMm:    toNumeric(coin.ThicknessMM),
		Edge:           toNullString(coin.Edge),
		Shape:          toNullString(coin.Shape),
		AcquiredAt:     toNullDate(coin.AcquiredAt),
		SoldAt:         toNullDate(coin.SoldAt),
		PricePaid:      toNumeric(coin.PricePaid),
		SoldPrice:      toNumeric(coin.SoldPrice),
	}

	result, err := r.q.CreateCoin(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save coin: %w", err)
	}

	// coin.ID is already set correctly by service, and we forced it in DB.
	// But let's update metadata just in case.
	coin.CreatedAt = result.CreatedAt.Time
	coin.UpdatedAt = result.UpdatedAt.Time

	// Save Images
	for _, img := range coin.Images {
		imgParams := db.CreateCoinImageParams{
			CoinID:           pgtype.UUID{Bytes: coin.ID, Valid: true},
			ImageType:        db.ImageType(img.ImageType),
			Side:             db.CoinSide(img.Side),
			Path:             img.Path,
			Extension:        img.Extension,
			Size:             img.Size,
			Width:            int32(img.Width),
			Height:           int32(img.Height),
			MimeType:         img.MimeType,
			OriginalFilename: toNullString(img.OriginalFilename),
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

func (r *PostgresCoinRepository) GetTotalValue(ctx context.Context) (float64, error) {
	return r.q.GetTotalValue(ctx)
}

func (r *PostgresCoinRepository) ListTopValuable(ctx context.Context) ([]*domain.Coin, error) {
	rows, err := r.q.ListTopValuableCoins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list top valuable coins: %w", err)
	}
	return r.rowsToCoins(ctx, rows)
}

func (r *PostgresCoinRepository) ListRecent(ctx context.Context) ([]*domain.Coin, error) {
	rows, err := r.q.ListRecentCoins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent coins: %w", err)
	}
	return r.rowsToCoins(ctx, rows)
}

func (r *PostgresCoinRepository) GetMaterialDistribution(ctx context.Context) (map[string]int, error) {
	rows, err := r.q.GetMaterialDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get material distribution: %w", err)
	}
	dist := make(map[string]int)
	for _, row := range rows {
		dist[row.Material.String] = int(row.Count)
	}
	return dist, nil
}

func (r *PostgresCoinRepository) GetGradeDistribution(ctx context.Context) (map[string]int, error) {
	rows, err := r.q.GetGradeDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get grade distribution: %w", err)
	}
	dist := make(map[string]int)
	for _, row := range rows {
		dist[string(row.Grade.GradeType)] = int(row.Count)
	}
	return dist, nil
}

func (r *PostgresCoinRepository) GetAllValues(ctx context.Context) ([]float64, error) {
	rows, err := r.q.GetAllValues(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all values: %w", err)
	}
	values := make([]float64, len(rows))
	for i, val := range rows {
		f, _ := val.Float64Value()
		values[i] = f.Float64
	}
	return values, nil
}

// Helper to avoid duplication
func (r *PostgresCoinRepository) rowsToCoins(ctx context.Context, rows []db.Coin) ([]*domain.Coin, error) {
	coins := make([]*domain.Coin, len(rows))
	for i, row := range rows {
		c, err := toDomainCoin(row)
		if err != nil {
			return nil, err
		}
		// Fetch images (N+1, acceptable for small lists like top 3)
		images, err := r.q.ListCoinImagesByCoinID(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list coin images for coin %v: %w", row.ID, err)
		}
		c.Images = toDomainImages(images)
		coins[i] = c
	}
	return coins, nil
}

func (r *PostgresCoinRepository) Update(ctx context.Context, coin *domain.Coin) error {
	// Convert map to JSONB byte array
	geminiDetailsBytes, err := json.Marshal(coin.GeminiDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal gemini details: %w", err)
	}

	params := db.UpdateCoinParams{
		ID:             pgtype.UUID{Bytes: coin.ID, Valid: true},
		Name:           toNullString(coin.Name),
		Mint:           toNullString(coin.Mint),
		Mintage:        toNullInt8(coin.Mintage),
		Country:        toNullString(coin.Country),
		Year:           toNullInt4(coin.Year),
		FaceValue:      toNullString(coin.FaceValue),
		Currency:       toNullString(coin.Currency),
		Material:       toNullString(coin.Material),
		Description:    toNullString(coin.Description),
		KmCode:         toNullString(coin.KMCode),
		MinValue:       toNumeric(coin.MinValue),
		MaxValue:       toNumeric(coin.MaxValue),
		Grade:          toNullGradeType(coin.Grade),
		TechnicalNotes: toNullString(coin.TechnicalNotes),
		GeminiDetails:  geminiDetailsBytes,
		GroupID:        toNullInt4Ptr(coin.GroupID),
		PersonalNotes:  toNullString(coin.PersonalNotes),
		WeightG:        toNumeric(coin.WeightG),
		DiameterMm:     toNumeric(coin.DiameterMM),
		ThicknessMm:    toNumeric(coin.ThicknessMM),
		Edge:           toNullString(coin.Edge),
		Shape:          toNullString(coin.Shape),
		AcquiredAt:     toNullDate(coin.AcquiredAt),
		SoldAt:         toNullDate(coin.SoldAt),
		PricePaid:      toNumeric(coin.PricePaid),
		SoldPrice:      toNumeric(coin.SoldPrice),
	}

	result, err := r.q.UpdateCoin(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update coin: %w", err)
	}

	coin.UpdatedAt = result.UpdatedAt.Time
	return nil
}

func (r *PostgresCoinRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteCoin(ctx, pgtype.UUID{Bytes: id, Valid: true})
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
	weightG, _ := row.WeightG.Float64Value()
	diameterMM, _ := row.DiameterMm.Float64Value()
	thicknessMM, _ := row.ThicknessMm.Float64Value()
	pricePaid, _ := row.PricePaid.Float64Value()
	soldPrice, _ := row.SoldPrice.Float64Value()

	var groupID *int
	if row.GroupID.Valid {
		id := int(row.GroupID.Int32)
		groupID = &id
	}

	var acquiredAt *time.Time
	if row.AcquiredAt.Valid {
		t := row.AcquiredAt.Time
		acquiredAt = &t
	}

	var soldAt *time.Time
	if row.SoldAt.Valid {
		t := row.SoldAt.Time
		soldAt = &t
	}

	return &domain.Coin{
		ID:             uuid.UUID(row.ID.Bytes),
		Name:           row.Name.String,
		Mint:           row.Mint.String,
		Mintage:        row.Mintage.Int64,
		Country:        row.Country.String,
		Year:           int(row.Year.Int32),
		FaceValue:      row.FaceValue.String,
		Currency:       row.Currency.String,
		Material:       row.Material.String,
		Description:    row.Description.String,
		KMCode:         row.KmCode.String,
		MinValue:       minVal.Float64,
		MaxValue:       maxVal.Float64,
		Grade:          string(row.Grade.GradeType),
		TechnicalNotes: row.TechnicalNotes.String,
		GeminiDetails:  geminiDetails,
		GroupID:        groupID,
		PersonalNotes:  row.PersonalNotes.String,
		WeightG:        weightG.Float64,
		DiameterMM:     diameterMM.Float64,
		ThicknessMM:    thicknessMM.Float64,
		Edge:           row.Edge.String,
		Shape:          row.Shape.String,
		AcquiredAt:     acquiredAt,
		SoldAt:         soldAt,
		PricePaid:      pricePaid.Float64,
		SoldPrice:      soldPrice.Float64,
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}, nil
}

func toDomainImages(rows []db.CoinImage) []domain.CoinImage {
	images := make([]domain.CoinImage, len(rows))
	for i, row := range rows {
		images[i] = domain.CoinImage{
			ID:               uuid.UUID(row.ID.Bytes),
			CoinID:           uuid.UUID(row.CoinID.Bytes),
			ImageType:        string(row.ImageType),
			Side:             string(row.Side),
			Path:             row.Path,
			Extension:        row.Extension,
			Size:             row.Size,
			Width:            int(row.Width),
			Height:           int(row.Height),
			MimeType:         row.MimeType,
			OriginalFilename: row.OriginalFilename.String,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        row.UpdatedAt.Time,
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
	if f == 0 {
		return pgtype.Numeric{Valid: false}
	}
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

func toNullInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: i,
		Valid: i != 0,
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

func toNullDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{
		Time:  *t,
		Valid: true,
	}
}
