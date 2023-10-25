package services_test

import (
	"errors"
	"go-unit-test/repositories"
	"go-unit-test/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromotionCalculateDiscount(t *testing.T) {

	type testCase struct {
		name         string
		purchaseMin  float64
		Discount     float64
		amount       float64
		expected     float64
		DiscountType string
	}

	cases := []testCase{
		{name: "applied 100", purchaseMin: 100, Discount: 20, amount: 100, expected: 80, DiscountType: "percentage"},
		{name: "applied 200", purchaseMin: 100, Discount: 20, amount: 200, expected: 160, DiscountType: "percentage"},
		{name: "applied 300", purchaseMin: 100, Discount: 20, amount: 300, expected: 240, DiscountType: "percentage"},
		{name: "not applied 50", purchaseMin: 100, Discount: 20, amount: 50, expected: 0, DiscountType: "percentage"},
		{name: "not applied 0", purchaseMin: 100, Discount: 20, amount: 0, expected: 0, DiscountType: "percentage"},

		// amount
		{name: "applied 100", purchaseMin: 100, Discount: 20, amount: 100, expected: 80, DiscountType: "amount"},
		{name: "applied 200", purchaseMin: 100, Discount: 20, amount: 200, expected: 180, DiscountType: "amount"},
		{name: "applied 300", purchaseMin: 100, Discount: 20, amount: 300, expected: 280, DiscountType: "amount"},
		{name: "not applied 50", purchaseMin: 100, Discount: 20, amount: 50, expected: 0, DiscountType: "amount"},
		{name: "not applied 0", purchaseMin: 100, Discount: 20, amount: 0, expected: 0, DiscountType: "amount"},

		// unknown
		{name: "unknown", purchaseMin: 100, Discount: 20, amount: 100, expected: 100, DiscountType: "unknown"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			//Arrage
			promoRepo := repositories.NewPromotionRepositoryMock()
			promoRepo.On("GetPromotion").Return(repositories.Promotion{
				ID:           1,
				PurchaseMin:  c.purchaseMin,
				Discount:     c.Discount,
				DiscountType: c.DiscountType,
			}, nil)

			promoService := services.NewPromotionService(promoRepo)

			//Act
			discount, _ := promoService.CalculateDiscount(c.amount)
			expected := c.expected

			//Assert
			assert.Equal(t, expected, discount)
		})
	}

	t.Run("purchase amount zero", func(t *testing.T) {
		//Arrage
		promoRepo := repositories.NewPromotionRepositoryMock()
		promoRepo.On("GetPromotion").Return(repositories.Promotion{
			ID:          1,
			PurchaseMin: 100,
			Discount:    20,
		}, nil)

		promoService := services.NewPromotionService(promoRepo)

		//Act
		_, err := promoService.CalculateDiscount(0)

		//Assert
		assert.ErrorIs(t, err, services.ErrZeroAmount)
		promoRepo.AssertNotCalled(t, "GetPromotion")
	})

	t.Run("repository error", func(t *testing.T) {
		//Arrage
		promoRepo := repositories.NewPromotionRepositoryMock()
		promoRepo.On("GetPromotion").Return(repositories.Promotion{}, errors.New(""))

		promoService := services.NewPromotionService(promoRepo)

		//Act
		_, err := promoService.CalculateDiscount(100)

		//Assert
		assert.ErrorIs(t, err, services.ErrRepository)
	})

}

func TestPromotionServiceMock_CalculateDiscount(t *testing.T) {
	// Create a new mock instance
	mockService := services.NewPromotionServiceMock()

	// Define test cases
	tests := []struct {
		name           string
		input          float64
		expectedResult float64
		expectedError  error
		mockReturn     []interface{}
	}{
		{
			name:           "Successful discount calculation",
			input:          100,
			expectedResult: 20,
			expectedError:  nil,
			mockReturn:     []interface{}{float64(20), nil},
		},
		{
			name:           "Error in discount calculation",
			input:          -10,
			expectedResult: 0,
			expectedError:  errors.New("invalid amount"),
			mockReturn:     []interface{}{float64(0), errors.New("invalid amount")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock behavior
			mockService.On("CalculateDiscount", tt.input).Return(tt.mockReturn...)

			// Call the method
			result, err := mockService.CalculateDiscount(tt.input)

			// Assertions
			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert that the expected method was called with the expected arguments
			mockService.AssertExpectations(t)
		})
	}
}
