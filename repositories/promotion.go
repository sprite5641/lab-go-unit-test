package repositories

type PromotionRepository interface {
	GetPromotion() (Promotion, error)
}

type Promotion struct {
	ID           int
	PurchaseMin  float64
	Discount     float64
	DiscountType string
}
