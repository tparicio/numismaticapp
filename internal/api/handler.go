package api

import (
	"strconv"

	"github.com/antonioparicio/numismaticapp/internal/application"
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

	// Call service
	coin, err := h.service.AddCoin(c.Context(), frontFile, backFile, groupName, userNotes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(coin)
}

func (h *CoinHandler) ListGroups(c *fiber.Ctx) error {
	groups, err := h.service.ListGroups(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groups)
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

	coins, err := h.service.ListCoins(c.Context(), limit, offset)
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
