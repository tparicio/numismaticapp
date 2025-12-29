package api

import (
	"fmt"
	"strconv"

	"github.com/antonioparicio/numismaticapp/internal/application"
	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CoinHandler struct {
	service  *application.CoinService
	validate *validator.Validate
}

func NewCoinHandler(service *application.CoinService) *CoinHandler {
	return &CoinHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *CoinHandler) AddCoin(c *fiber.Ctx) error {
	// Parse multipart form
	frontFile, err := c.FormFile("front_image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "front_image is required"})
	}

	backFile, err := c.FormFile("back_image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "back_image is required"})
	}

	groupName := c.FormValue("group_name")
	userNotes := c.FormValue("user_notes")
	name := c.FormValue("name")
	mint := c.FormValue("mint")
	mintageStr := c.FormValue("mintage")
	mintage := 0
	if mintageStr != "" {
		if _, err := fmt.Sscanf(mintageStr, "%d", &mintage); err != nil {
			// Just log debug, mintage 0 is fine default
			fmt.Printf("Failed to parse mintage: %v\n", err)
		}
	}

	modelName := c.FormValue("model_name")
	temperatureStr := c.FormValue("temperature")
	var temperature float32 = 0.1 // Default
	if temperatureStr != "" {
		if val, err := strconv.ParseFloat(temperatureStr, 32); err == nil {
			temperature = float32(val)
		}
	}

	// Open files
	frontSrc, err := frontFile.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to open front file"})
	}
	defer func() {
		if err := frontSrc.Close(); err != nil {
			fmt.Printf("Failed to close front file: %v\n", err)
		}
	}()

	backSrc, err := backFile.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to open back file"})
	}
	defer func() {
		if err := backSrc.Close(); err != nil {
			fmt.Printf("Failed to close back file: %v\n", err)
		}
	}()

	// Call service
	coin, err := h.service.AddCoin(c.Context(), frontSrc, frontFile.Filename, backSrc, backFile.Filename, groupName, userNotes, name, mint, mintage, modelName, temperature)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(coin)
}

func (h *CoinHandler) ListGeminiModels(c *fiber.Ctx) error {
	models, err := h.service.GetGeminiModels(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list models"})
	}
	return c.JSON(models)
}

func (h *CoinHandler) ListGroups(c *fiber.Ctx) error {
	groups, err := h.service.ListGroups(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groups)
}

type RotateCoinRequest struct {
	Side  string  `json:"side" validate:"required,oneof=front back"`
	Angle float64 `json:"angle" validate:"required"`
}

func (h *CoinHandler) RotateCoin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid UUID"})
	}

	var req RotateCoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse body"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.RotateCoinImage(c.Context(), id, req.Side, req.Angle); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=200"`
}

func (h *CoinHandler) CreateGroup(c *fiber.Ctx) error {
	var req CreateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse body"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	group, err := h.service.CreateGroup(c.Context(), req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *CoinHandler) UpdateGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req CreateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse body"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	group, err := h.service.UpdateGroup(c.Context(), id, req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(group)
}

func (h *CoinHandler) DeleteGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.service.DeleteGroup(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CoinHandler) ListCoins(c *fiber.Ctx) error {
	// Parse filters
	limit := 50
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	filter := domain.CoinFilter{
		Limit:     limit,
		Offset:    offset,
		Query:     strPtr(c.Query("q")),
		Country:   strPtr(c.Query("country")),
		Grade:     strPtr(c.Query("grade")),
		Material:  strPtr(c.Query("material")),
		SortBy:    strPtr(c.Query("sort_by")),
		SortOrder: strPtr(c.Query("order")),
	}

	if g := c.Query("group_id"); g != "" {
		if val, err := strconv.Atoi(g); err == nil {
			filter.GroupID = &val
		}
	}

	if y := c.Query("year"); y != "" {
		if val, err := strconv.Atoi(y); err == nil {
			filter.Year = &val
		}
	}

	if my := c.Query("min_year"); my != "" {
		if val, err := strconv.Atoi(my); err == nil {
			filter.MinYear = &val
		}
	}

	if my := c.Query("max_year"); my != "" {
		if val, err := strconv.Atoi(my); err == nil {
			filter.MaxYear = &val
		}
	}

	if mp := c.Query("min_price"); mp != "" {
		if val, err := strconv.ParseFloat(mp, 64); err == nil {
			filter.MinPrice = &val
		}
	}

	if mp := c.Query("max_price"); mp != "" {
		if val, err := strconv.ParseFloat(mp, 64); err == nil {
			filter.MaxPrice = &val
		}
	}

	if so := c.Query("order"); so != "" {
		filter.SortOrder = &so
	}

	coins, err := h.service.ListCoins(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(coins)
}

func (h *CoinHandler) GetCoin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	coin, err := h.service.GetCoin(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "coin not found"})
	}

	return c.JSON(coin)
}

func (h *CoinHandler) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.service.GetDashboardStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

func (h *CoinHandler) UpdateCoin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	var req application.UpdateCoinParams
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	coin, err := h.service.UpdateCoin(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(coin)
}

func (h *CoinHandler) DeleteCoin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	if err := h.service.DeleteCoin(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type ReanalyzeRequest struct {
	ModelName   string  `json:"model_name" validate:"required"`
	Temperature float32 `json:"temperature" validate:"gte=0,lte=1"`
}

func (h *CoinHandler) ReanalyzeCoin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	var req ReanalyzeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Temperature == 0 {
		req.Temperature = 0.1
	}

	coin, err := h.service.ReanalyzeCoin(c.Context(), id, req.ModelName, req.Temperature)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(coin)
}

func (h *CoinHandler) ReprocessNumista(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	if err := h.service.EnrichCoinWithNumista(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *CoinHandler) ApplyNumistaResult(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	numistaIDStr := c.Params("numista_id")
	numistaID, err := strconv.Atoi(numistaIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid numista id"})
	}

	coin, err := h.service.ApplyNumistaCandidate(c.Context(), id, numistaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(coin)
}

type SellCoinRequest struct {
	SoldPrice   float64 `json:"sold_price" validate:"required,gt=0"`
	SaleChannel string  `json:"sale_channel" validate:"required,min=1"`
}

func (h *CoinHandler) SellCoin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	var req SellCoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	coin, err := h.service.MarkCoinAsSold(c.Context(), id, req.SoldPrice, req.SaleChannel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(coin)
}

func (h *CoinHandler) GetSaleChannels(c *fiber.Ctx) error {
	channels, err := h.service.GetSaleChannels(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(channels)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (h *CoinHandler) ExportCSV(c *fiber.Ctx) error {
	data, err := h.service.ExportCoinsCSV(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="coins.csv"`)
	return c.Send(data)
}

func (h *CoinHandler) ExportSQL(c *fiber.Ctx) error {
	data, err := h.service.ExportCoinsSQL(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "application/sql")
	c.Set("Content-Disposition", `attachment; filename="backup.sql"`)
	return c.Send(data)
}

type AddLinkRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func (h *CoinHandler) ListCoinLinks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	links, err := h.service.GetLinks(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(links)
}

func (h *CoinHandler) AddCoinLink(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	}

	var req AddLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	link, err := h.service.AddLink(c.Context(), id, req.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(link)
}

func (h *CoinHandler) DeleteCoinLink(c *fiber.Ctx) error {
	linkIDStr := c.Params("link_id")
	linkID, err := uuid.Parse(linkIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid link uuid"})
	}

	if err := h.service.RemoveLink(c.Context(), linkID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CoinHandler) RefreshCoinLink(c *fiber.Ctx) error {
	linkIDStr := c.Params("link_id")
	linkID, err := uuid.Parse(linkIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid link uuid"})
	}

	link, err := h.service.RefreshLink(c.Context(), linkID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(link)
}

func (h *CoinHandler) AddGroupImage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image is required"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to open image"})
	}
	defer src.Close()

	if err := h.service.AddGroupImage(c.Context(), groupID, src, file.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *CoinHandler) RemoveGroupImage(c *fiber.Ctx) error {
	imgIDStr := c.Params("image_id")
	imgID, err := uuid.Parse(imgIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid image uuid"})
	}

	if err := h.service.RemoveGroupImage(c.Context(), imgID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CoinHandler) AddCoinGalleryImage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	coinID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid coin uuid"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image is required"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to open image"})
	}
	defer src.Close()

	if err := h.service.AddCoinGalleryImage(c.Context(), coinID, src, file.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *CoinHandler) RemoveCoinGalleryImage(c *fiber.Ctx) error {
	imgIDStr := c.Params("image_id")
	imgID, err := uuid.Parse(imgIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid image uuid"})
	}

	if err := h.service.RemoveCoinGalleryImage(c.Context(), imgID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CoinHandler) ListGroupImages(c *fiber.Ctx) error {
	idStr := c.Params("id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}

	images, err := h.service.ListGroupImages(c.Context(), groupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(images)
}

func (h *CoinHandler) ListCoinGalleryImages(c *fiber.Ctx) error {
	idStr := c.Params("id")
	coinID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid coin uuid"})
	}

	images, err := h.service.ListCoinGalleryImages(c.Context(), coinID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(images)
}
