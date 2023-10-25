// go:build integration
package handlers_test

import (
	"fmt"
	"go-unit-test/handlers"
	"go-unit-test/repositories"
	"go-unit-test/services"
	"net/http/httptest"
	"testing"
)

func TestPromotionCalculateDiscountIntegrationService(T *testing.T) {

	amount := 100
	expected := 80

	promoRepo := repositories.NewPromotionRepositoryMock()
	promoRepo.On("GetPromotion").Return(repositories.Promotion{
		PurchaseMin:  100,
		Discount:     20,
		DiscountType: "percentage",
	}, nil)

	promoServ := services.NewPromotionService(promoRepo)

	promoHand := handlers.NewPromotionHandler(promoServ)

	app := setupApp()

	app.Get("/calculate-discount", promoHand.CalculateDiscount)

	res, err := app.Test(httptest.NewRequest("GET", fmt.Sprintf("/calculate-discount?amount=%d", amount), nil))

	assertError(T, err)

	assertResponse(T, res, 200, fmt.Sprintf(`{"discount":%d}`, expected))

}
