package handlers

import (
	"go-unit-test/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PromotionHandler interface {
	CalculateDiscount(c *fiber.Ctx) error
}

type promotionHandler struct {
	promoService services.PromotionService
}

func NewPromotionHandler(promoService services.PromotionService) PromotionHandler {
	return promotionHandler{promoService: promoService}
}

func (h promotionHandler) CalculateDiscount(c *fiber.Ctx) error {
	amount := c.Query("amount")

	if amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "amount is required",
		})
	}

	amount64, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	discount, err := h.promoService.CalculateDiscount(amount64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"discount": discount,
	})
}
