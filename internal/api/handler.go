package api

import (
	"fmt"
	"strconv"

	"github.com/antonioparicio/numismaticapp/internal/application"
	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CoinHandler struct {
	service *application.CoinService
}

func NewCoinHandler(service *application.CoinService) *CoinHandler {
	return &CoinHandler{service: service}
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
		fmt.Sscanf(mintageStr, "%d", &mintage)
	}

	modelName := c.FormValue("model_name")
	temperatureStr := c.FormValue("temperature")
	var temperature float32 = 0.4 // Default
	if temperatureStr != "" {
		if val, err := strconv.ParseFloat(temperatureStr, 32); err == nil {
			temperature = float32(val)
		}
	}

	// Call service
	coin, err := h.service.AddCoin(c.Context(), frontFile, backFile, groupName, userNotes, name, mint, mintage, modelName, temperature)
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

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *CoinHandler) CreateGroup(c *fiber.Ctx) error {
	var req CreateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse body"})
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
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	filter := domain.CoinFilter{
		Limit:  limit,
		Offset: offset,
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

	if country := c.Query("country"); country != "" {
		filter.Country = &country
	}

	if q := c.Query("q"); q != "" {
		filter.Query = &q
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = &val
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = &val
		}
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
	ModelName   string  `json:"model_name"`
	Temperature float32 `json:"temperature"`
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

	if req.Temperature == 0 {
		req.Temperature = 0.4
	}

	coin, err := h.service.ReanalyzeCoin(c.Context(), id, req.ModelName, req.Temperature)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(coin)
}
