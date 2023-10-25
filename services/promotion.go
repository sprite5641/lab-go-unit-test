package services

import (
	"go-unit-test/repositories"
)

type PromotionService interface {
	CalculateDiscount(amount float64) (float64, error)
}

type promotionService struct {
	promoRepo repositories.PromotionRepository
}

func NewPromotionService(promoRepo repositories.PromotionRepository) PromotionService {
	return promotionService{promoRepo: promoRepo}
}

func (s promotionService) CalculateDiscount(amount float64) (float64, error) {

	if amount <= 0 {
		return 0, ErrZeroAmount
	}

	promotion, err := s.promoRepo.GetPromotion()
	if err != nil {
		return 0, ErrRepository
	}

	if amount < promotion.PurchaseMin {
		return 0, nil
	}

	if promotion.DiscountType == "percentage" {
		return amount - (promotion.Discount * amount / 100), nil
	}

	if promotion.DiscountType == "amount" {
		return amount - promotion.Discount, nil
	}

	return amount, nil
}
