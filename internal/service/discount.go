package service

import (
	"context"

	"github.com/Badgain/book-discount/internal/domain"
)

/*
Бизнес правила:

1. Клиент еще не соверщаш покупку и берет от двух до пяти книг - скидка 20% на весь счет
2. Клиент уже совершал покупки в магазине, берет от 2 до 5 книг - скидка 10%
3. Клиент уже совершал покупки в магазине, берет от 6 до 10 книг - скидка 5%
4. Клиент уже совершал покупки в магазине, берет свяше 10 книг - скидка 2%
5. Если сегодня пятница - всем скидка 5% не зависимо от объема корзины или типа клиента

*/

// DiscountService реализует бизнес-логику расчета скидок с использованием правил
type DiscountService struct{}

// NewDiscountService создает новый экземпляр DiscountService с правилами по умолчанию
func NewDiscountService() *DiscountService {
	return &DiscountService{}
}

// Calculate вычисляет скидку на основе правил
func (s *DiscountService) Calculate(ctx context.Context, customer domain.CustomerType, books []domain.Book) (domain.Discount, error) {
	return domain.Discount{}, nil
}
