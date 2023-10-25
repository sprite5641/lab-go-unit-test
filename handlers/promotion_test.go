package handlers_test

import (
	"fmt"
	"go-unit-test/handlers"
	"go-unit-test/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPromotionCalculateDiscount(t *testing.T) {
	app := setupApp()

	t.Run("success", func(t *testing.T) {
		testSuccessScenario(t, app)
	})

	t.Run("amount is required", func(t *testing.T) {
		testAmountIsRequiredScenario(t, app)
	})

	t.Run("amount is not a number", func(t *testing.T) {
		testAmountIsNotANumberScenario(t, app)
	})

	t.Run("error in discount calculation", func(t *testing.T) {
		testErrorInDiscountCalculation(t)
	})
}

func setupApp() *fiber.App {
	return fiber.New()
}

func testSuccessScenario(t *testing.T, app *fiber.App) {
	amount := 100
	expected := 20

	promoServ := services.NewPromotionServiceMock()
	promoServ.On("CalculateDiscount", float64(amount)).Return(float64(expected), nil)
	promoHand := handlers.NewPromotionHandler(promoServ)
	app.Get("/calculate-discount", promoHand.CalculateDiscount)
	res, err := app.Test(httptest.NewRequest("GET", fmt.Sprintf("/calculate-discount?amount=%d", amount), nil))
	assertError(t, err)

	assertResponse(t, res, fiber.StatusOK, fmt.Sprintf(`{"discount":%d}`, expected))
}

func testAmountIsRequiredScenario(t *testing.T, app *fiber.App) {
	promoServ := services.NewPromotionServiceMock()
	promoServ.On("CalculateDiscount", float64(0)).Return(float64(0), nil)
	promoHand := handlers.NewPromotionHandler(promoServ)
	app.Get("/calculate-discount", promoHand.CalculateDiscount)
	res, err := app.Test(httptest.NewRequest("GET", "/calculate-discount?amount=", nil))
	assertError(t, err)
	assertResponse(t, res, fiber.StatusBadRequest, `{"message":"amount is required"}`)
}

func testAmountIsNotANumberScenario(t *testing.T, app *fiber.App) {
	promoServ := services.NewPromotionServiceMock()
	promoServ.On("CalculateDiscount", float64(0)).Return(float64(0), nil)
	promoHand := handlers.NewPromotionHandler(promoServ)
	app.Get("/calculate-discount", promoHand.CalculateDiscount)
	res, err := app.Test(httptest.NewRequest("GET", "/calculate-discount?amount=abc", nil))
	assertError(t, err)
	assertResponse(t, res, fiber.StatusBadRequest, `{"message":"strconv.ParseFloat: parsing \"abc\": invalid syntax"}`)
}

func testErrorInDiscountCalculation(t *testing.T) {
	promoServ := services.NewPromotionServiceMock()
	promoServ.On("CalculateDiscount", float64(100)).Return(float64(0), services.ErrRepository)

	promoHand := handlers.NewPromotionHandler(promoServ)

	app := fiber.New()
	app.Get("/calculate-discount", promoHand.CalculateDiscount)
	res, err := app.Test(httptest.NewRequest("GET", "/calculate-discount?amount=100", nil))
	assertError(t, err)
	assertResponse(t, res, fiber.StatusInternalServerError, `{"message":"repository error"}`)
}

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assertResponse(t *testing.T, res *http.Response, expectedStatusCode int, expectedBody string) {
	if assert.Equal(t, expectedStatusCode, res.StatusCode) {
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, expectedBody, string(body))
	}
}
