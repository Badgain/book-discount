package domain

import "context"

// DiscountCalculator интерфейс для расчета скидок
type DiscountCalculator interface {
	Calculate(ctx context.Context, customerType CustomerType, books []Book) (Discount, error)
}
