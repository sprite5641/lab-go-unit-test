package services

import "github.com/stretchr/testify/mock"

type PromotionServiceMock struct {
	mock.Mock
}

func NewPromotionServiceMock() *PromotionServiceMock {
	return &PromotionServiceMock{}
}

func (m *PromotionServiceMock) CalculateDiscount(amount float64) (float64, error) {
	args := m.Called(amount)
	return args.Get(0).(float64), args.Error(1)
}
