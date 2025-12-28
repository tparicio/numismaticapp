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
	params, err := toDBParams(coin)
	if err != nil {
		return err
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

func (r *PostgresCoinRepository) AddImage(ctx context.Context, img domain.CoinImage) error {
	imgParams := db.CreateCoinImageParams{
		CoinID:           pgtype.UUID{Bytes: img.CoinID, Valid: true},
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
	return nil
}

func (r *PostgresCoinRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Coin, error) {
	row, err := r.q.GetCoin(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get coin: %w", err)
	}

	return r.coinFromDB(ctx, row)
}

func (r *PostgresCoinRepository) List(ctx context.Context, filter domain.CoinFilter) ([]*domain.Coin, error) {
	params := db.ListCoinsParams{
		Limit:     int32(filter.Limit),
		Offset:    int32(filter.Offset),
		GroupID:   toNullInt4Ptr(filter.GroupID),
		Year:      toNullInt4Ptr(filter.Year),
		Country:   toNullStringPtr(filter.Country),
		Query:     toNullStringPtr(filter.Query),
		MinPrice:  toNullFloat8Ptr(filter.MinPrice),
		MaxPrice:  toNullFloat8Ptr(filter.MaxPrice),
		Grade:     toNullStringPtr(filter.Grade),
		Material:  toNullStringPtr(filter.Material),
		MinYear:   toNullInt4Ptr(filter.MinYear),
		MaxYear:   toNullInt4Ptr(filter.MaxYear),
		SortBy:    toNullStringPtr(filter.SortBy),
		SortOrder: toNullStringPtr(filter.SortOrder),
	}

	rows, err := r.q.ListCoins(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list coins: %w", err)
	}

	// Use the shared helper which generates coins and batch-fetches images
	return r.rowsToCoins(ctx, rows)
}

func (r *PostgresCoinRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountCoins(ctx)
}

func (r *PostgresCoinRepository) GetTotalValue(ctx context.Context) (float64, error) {
	return r.q.GetTotalValue(ctx)
}

func (r *PostgresCoinRepository) GetAverageValue(ctx context.Context) (float64, error) {
	return r.q.GetAverageValue(ctx)
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
		dist[row.Grade.String] = int(row.Count)
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

func (r *PostgresCoinRepository) GetCountryDistribution(ctx context.Context) (map[string]int, error) {
	rows, err := r.q.GetCountryDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get country distribution: %w", err)
	}
	dist := make(map[string]int)
	for _, row := range rows {
		dist[row.Country.String] = int(row.Count)
	}
	return dist, nil
}

func (r *PostgresCoinRepository) GetOldestCoin(ctx context.Context) (*domain.Coin, error) {
	row, err := r.q.GetOldestCoin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest coin: %w", err)
	}
	return toDomainCoin(row)
}

func (r *PostgresCoinRepository) GetRarestCoins(ctx context.Context, limit int) ([]*domain.Coin, error) {
	rows, err := r.q.GetRarestCoins(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get rarest coins: %w", err)
	}
	return r.rowsToCoins(ctx, rows)
}

func (r *PostgresCoinRepository) GetGroupDistribution(ctx context.Context) (map[string]int, error) {
	rows, err := r.q.GetGroupDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get group distribution: %w", err)
	}
	dist := make(map[string]int)
	for _, row := range rows {
		dist[row.GroupName] = int(row.Count)
	}
	return dist, nil
}

func (r *PostgresCoinRepository) GetGroupStats(ctx context.Context) ([]domain.GroupStat, error) {
	rows, err := r.q.GetGroupStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get group stats: %w", err)
	}

	stats := make([]domain.GroupStat, len(rows))
	for i, row := range rows {
		// Handle null group ID
		groupID := 0
		if row.GroupID.Valid {
			groupID = int(row.GroupID.Int32)
		}

		stats[i] = domain.GroupStat{
			GroupID:   groupID,
			GroupName: row.GroupName,
			Count:     row.Count,
			MinVal:    row.MinVal,
			MaxVal:    row.MaxVal,
		}
	}
	return stats, nil
}

func (r *PostgresCoinRepository) GetTotalWeightByMaterial(ctx context.Context, materialLike string) (float64, error) {
	weight, err := r.q.GetTotalWeightByMaterial(ctx, toNullString(materialLike))
	if err != nil {
		return 0, fmt.Errorf("failed to get total weight: %w", err)
	}
	return weight, nil
}

func (r *PostgresCoinRepository) GetHeaviestCoin(ctx context.Context) (*domain.Coin, error) {
	row, err := r.q.GetHeaviestCoin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get heaviest coin: %w", err)
	}
	return r.coinFromDB(ctx, row)
}

func (r *PostgresCoinRepository) GetSmallestCoin(ctx context.Context) (*domain.Coin, error) {
	row, err := r.q.GetSmallestCoin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get smallest coin: %w", err)
	}
	return r.coinFromDB(ctx, row)
}

func (r *PostgresCoinRepository) GetRandomCoin(ctx context.Context) (*domain.Coin, error) {
	row, err := r.q.GetRandomCoin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get random coin: %w", err)
	}
	return r.coinFromDB(ctx, row)
}

func (r *PostgresCoinRepository) GetAllCoins(ctx context.Context) ([]*domain.Coin, error) {
	rows, err := r.q.GetAllCoins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all coins: %w", err)
	}
	// For scatter plot, we might not need images, but let's keep it consistent or optimize later.
	// Optimization: Don't fetch images for scatter plot to avoid N+1 on large dataset.
	// The frontend only needs Year and Grade.
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

// Helper to avoid duplication
func (r *PostgresCoinRepository) rowsToCoins(ctx context.Context, rows []db.Coin) ([]*domain.Coin, error) {
	if len(rows) == 0 {
		return []*domain.Coin{}, nil
	}

	coins := make([]*domain.Coin, len(rows))
	coinIDs := make([]pgtype.UUID, len(rows))
	coinMap := make(map[uuid.UUID]*domain.Coin)

	for i, row := range rows {
		c, err := toDomainCoin(row)
		if err != nil {
			return nil, err
		}
		c.Images = []domain.CoinImage{} // Initialize empty slice
		coins[i] = c
		coinIDs[i] = pgtype.UUID{Bytes: c.ID, Valid: true}
		coinMap[c.ID] = c
	}

	// Batch fetch images
	images, err := r.q.ListCoinImagesByCoinIDs(ctx, coinIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to batch list coin images: %w", err)
	}

	// Map images to coins
	for _, img := range toDomainImages(images) {
		if coin, exists := coinMap[img.CoinID]; exists {
			coin.Images = append(coin.Images, img)
		}
	}

	return coins, nil
}

func (r *PostgresCoinRepository) Update(ctx context.Context, coin *domain.Coin) error {
	createParams, err := toDBParams(coin)
	if err != nil {
		return err
	}

	params := db.UpdateCoinParams(createParams)

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

func (r *PostgresGroupRepository) Update(ctx context.Context, group *domain.Group) error {
	row, err := r.q.UpdateGroup(ctx, db.UpdateGroupParams{
		ID:          int32(group.ID),
		Name:        group.Name,
		Description: toNullString(group.Description),
	})
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}
	// Update original struct with returned values (e.g. updated_at if we had it)
	group.Name = row.Name
	group.Description = row.Description.String
	return nil
}

func (r *PostgresGroupRepository) Delete(ctx context.Context, id int) error {
	return r.q.DeleteGroup(ctx, int32(id))
}

func (r *PostgresCoinRepository) coinFromDB(ctx context.Context, row db.Coin) (*domain.Coin, error) {
	coin, err := toDomainCoin(row)
	if err != nil {
		return nil, err
	}

	images, err := r.q.ListCoinImagesByCoinID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}
	coin.Images = toDomainImages(images)
	return coin, nil
}

func toDBParams(coin *domain.Coin) (db.CreateCoinParams, error) {
	geminiDetailsBytes, err := json.Marshal(coin.GeminiDetails)
	if err != nil {
		return db.CreateCoinParams{}, fmt.Errorf("failed to marshal gemini details: %w", err)
	}

	numistaDetailsBytes, err := json.Marshal(coin.NumistaDetails)
	if err != nil {
		return db.CreateCoinParams{}, fmt.Errorf("failed to marshal numista details: %w", err)
	}

	return db.CreateCoinParams{
		ID:                pgtype.UUID{Bytes: coin.ID, Valid: true},
		Name:              toNullString(coin.Name),
		Mint:              toNullString(coin.Mint),
		Mintage:           toNullInt8(coin.Mintage.Int64()),
		Country:           toNullString(coin.Country),
		Year:              toNullInt4(coin.Year.Int()),
		FaceValue:         toNullString(coin.FaceValue),
		Currency:          toNullString(coin.Currency),
		Material:          toNullString(coin.Material),
		Description:       toNullString(coin.Description),
		KmCode:            toNullString(coin.KMCode.String()),
		MinValue:          toNumeric(coin.MinValue),
		MaxValue:          toNumeric(coin.MaxValue),
		Grade:             toNullString(coin.Grade.String()),
		TechnicalNotes:    toNullString(coin.TechnicalNotes),
		GeminiDetails:     geminiDetailsBytes,
		GroupID:           toNullInt4Ptr(coin.GroupID),
		PersonalNotes:     toNullString(coin.PersonalNotes),
		WeightG:           toNumeric(coin.WeightG),
		DiameterMm:        toNumeric(coin.DiameterMM),
		ThicknessMm:       toNumeric(coin.ThicknessMM),
		Edge:              toNullString(coin.Edge),
		Shape:             toNullString(coin.Shape),
		AcquiredAt:        toNullDate(coin.AcquiredAt),
		SoldAt:            toNullDate(coin.SoldAt),
		PricePaid:         toNumeric(coin.PricePaid),
		SoldPrice:         toNumeric(coin.SoldPrice),
		NumistaNumber:     toNullInt4(coin.NumistaNumber),
		NumistaDetails:    numistaDetailsBytes,
		GeminiModel:       toNullString(coin.GeminiModel),
		GeminiTemperature: toNumeric(coin.GeminiTemperature),
		NumistaSearch:     toNullString(coin.NumistaSearch),
	}, nil
}

// Helper functions for conversion

func toDomainCoin(row db.Coin) (*domain.Coin, error) {
	var geminiDetails map[string]any
	if len(row.GeminiDetails) > 0 {
		if err := json.Unmarshal(row.GeminiDetails, &geminiDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal gemini details: %w", err)
		}
	}

	var numistaDetails map[string]any
	if len(row.NumistaDetails) > 0 {
		if err := json.Unmarshal(row.NumistaDetails, &numistaDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal numista details: %w", err)
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
	geminiTemp, _ := row.GeminiTemperature.Float64Value()

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

	// Handle factory errors by logging or defaulting?
	// For database retrieval, we trust the DB somewhat, but we should handle invalid data.
	// Since NewYear/NewMintage return errors for invalid data, we should ideally handle them.
	// However, toDomainCoin has a signature to assume data is somewhat valid or we bubble error.

	yearVO, err := domain.NewYear(int(row.Year.Int32))
	if err != nil {
		// If DB has invalid year, maybe return error or 0?
		// Let's allow it but log? No, domain shouldn't depend on log.
		// NewYear(0) is valid.
		yearVO, _ = domain.NewYear(0)
	}

	mintageVO, _ := domain.NewMintage(row.Mintage.Int64)
	kmVO, _ := domain.NewKMCode(row.KmCode.String)
	gradeVO, _ := domain.NewGrade(row.Grade.String)

	return &domain.Coin{
		ID:                uuid.UUID(row.ID.Bytes),
		Name:              row.Name.String,
		Mint:              row.Mint.String,
		Mintage:           mintageVO,
		Country:           row.Country.String,
		Year:              yearVO,
		FaceValue:         row.FaceValue.String,
		Currency:          row.Currency.String,
		Material:          row.Material.String,
		Description:       row.Description.String,
		KMCode:            kmVO,
		NumistaNumber:     int(row.NumistaNumber.Int32),
		NumistaDetails:    numistaDetails,
		Ruler:             "", // populated if needed
		Orientation:       "",
		Series:            "",
		CommemoratedTopic: "",
		MinValue:          minVal.Float64,
		MaxValue:          maxVal.Float64,
		Grade:             gradeVO,
		TechnicalNotes:    row.TechnicalNotes.String,
		GeminiDetails:     geminiDetails,
		GroupID:           groupID,
		PersonalNotes:     row.PersonalNotes.String,
		WeightG:           weightG.Float64,
		DiameterMM:        diameterMM.Float64,
		ThicknessMM:       thicknessMM.Float64,
		Edge:              row.Edge.String,
		Shape:             row.Shape.String,
		AcquiredAt:        acquiredAt,
		SoldAt:            soldAt,
		PricePaid:         pricePaid.Float64,
		SoldPrice:         soldPrice.Float64,
		GeminiModel:       row.GeminiModel.String,
		GeminiTemperature: geminiTemp.Float64,
		NumistaSearch:     row.NumistaSearch.String,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
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
	if err := n.Scan(fmt.Sprintf("%f", f)); err != nil {
		return pgtype.Numeric{Valid: false}
	}
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

func toNullStringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *s,
		Valid:  true,
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

func toNullFloat8Ptr(f *float64) pgtype.Float8 {
	if f == nil {
		return pgtype.Float8{Valid: false}
	}
	return pgtype.Float8{
		Float64: *f,
		Valid:   true,
	}
}

// MarkAsSold marks a coin as sold with price and channel
func (r *PostgresCoinRepository) MarkAsSold(ctx context.Context, id uuid.UUID, soldAt time.Time, soldPrice float64, saleChannel string) (*domain.Coin, error) {
	row, err := r.q.MarkCoinAsSold(ctx, db.MarkCoinAsSoldParams{
		ID:          pgtype.UUID{Bytes: id, Valid: true},
		SoldAt:      pgtype.Date{Time: soldAt, Valid: true},
		SoldPrice:   toNumeric(soldPrice),
		SaleChannel: toNullString(saleChannel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mark coin as sold: %w", err)
	}
	return r.coinFromDB(ctx, row)
}

// GetSaleChannels returns list of distinct sale channels
func (r *PostgresCoinRepository) GetSaleChannels(ctx context.Context) ([]string, error) {
	rows, err := r.q.GetDistinctSaleChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sale channels: %w", err)
	}
	channels := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.Valid && row.String != "" {
			channels = append(channels, row.String)
		}
	}
	return channels, nil
}

// AddLink adds a new coin link
func (r *PostgresCoinRepository) AddLink(ctx context.Context, link *domain.CoinLink) error {
	params := db.AddCoinLinkParams{
		CoinID:        pgtype.UUID{Bytes: link.CoinID, Valid: true},
		Url:           link.URL,
		Name:          toNullString(link.Name),
		OgTitle:       toNullString(link.OGTitle),
		OgDescription: toNullString(link.OGDescription),
		OgImage:       toNullString(link.OGImage),
	}

	row, err := r.q.AddCoinLink(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add link: %w", err)
	}

	link.ID = uuid.UUID(row.ID.Bytes)
	link.CreatedAt = row.CreatedAt.Time
	return nil
}

// RemoveLink removes a coin link
func (r *PostgresCoinRepository) RemoveLink(ctx context.Context, linkID uuid.UUID) error {
	return r.q.DeleteCoinLink(ctx, pgtype.UUID{Bytes: linkID, Valid: true})
}

// ListLinks lists all links for a coin
func (r *PostgresCoinRepository) ListLinks(ctx context.Context, coinID uuid.UUID) ([]*domain.CoinLink, error) {
	rows, err := r.q.ListCoinLinks(ctx, pgtype.UUID{Bytes: coinID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list links: %w", err)
	}

	links := make([]*domain.CoinLink, len(rows))
	for i, row := range rows {
		links[i] = &domain.CoinLink{
			ID:            uuid.UUID(row.ID.Bytes),
			CoinID:        uuid.UUID(row.CoinID.Bytes),
			URL:           row.Url,
			Name:          row.Name.String,
			OGTitle:       row.OgTitle.String,
			OGDescription: row.OgDescription.String,
			OGImage:       row.OgImage.String,
			CreatedAt:     row.CreatedAt.Time,
		}
	}
	return links, nil
}

// GetAllImages returns all images for export features
func (r *PostgresCoinRepository) GetAllImages(ctx context.Context) ([]domain.CoinImage, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, coin_id, image_type, side, path, extension, size, width, height, mime_type, original_filename, created_at, updated_at
		FROM coin_images
		ORDER BY coin_id, created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list all images: %w", err)
	}
	defer rows.Close()

	var images []domain.CoinImage
	for rows.Next() {
		var i domain.CoinImage
		var idBytes, coinIDBytes [16]byte
		var createdAt, updatedAt time.Time
		var width, height int32

		// Using generic scanning since we don't have the generated struct easily accessible/compatible without mapping
		err := rows.Scan(
			&idBytes,
			&coinIDBytes,
			&i.ImageType,
			&i.Side,
			&i.Path,
			&i.Extension,
			&i.Size,
			&width,
			&height,
			&i.MimeType,
			&i.OriginalFilename,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		i.ID = uuid.UUID(idBytes)
		i.CoinID = uuid.UUID(coinIDBytes)
		i.Width = int(width)
		i.Height = int(height)
		i.CreatedAt = createdAt
		i.UpdatedAt = updatedAt
		images = append(images, i)
	}
	return images, rows.Err()
}

// GetAllLinks returns all links for export features
func (r *PostgresCoinRepository) GetAllLinks(ctx context.Context) ([]*domain.CoinLink, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, coin_id, url, name, og_title, og_description, og_image, created_at
		FROM coin_links
		ORDER BY coin_id, created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list all links: %w", err)
	}
	defer rows.Close()

	var links []*domain.CoinLink
	for rows.Next() {
		l := &domain.CoinLink{}
		var idBytes, coinIDBytes [16]byte
		var pName, pOgTitle, pOgDesc, pOgImage pgtype.Text

		err := rows.Scan(
			&idBytes,
			&coinIDBytes,
			&l.URL,
			&pName,
			&pOgTitle,
			&pOgDesc,
			&pOgImage,
			&l.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		l.ID = uuid.UUID(idBytes)
		l.CoinID = uuid.UUID(coinIDBytes)
		l.Name = pName.String
		l.OGTitle = pOgTitle.String
		l.OGDescription = pOgDesc.String
		l.OGImage = pOgImage.String

		links = append(links, l)
	}
	return links, rows.Err()
}
