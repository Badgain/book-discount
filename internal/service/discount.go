package service

import (
	"context"
	"errors"
	"time"

	"github.com/Badgain/book-discount/internal/domain"
)

/*
Бизнес правила:

1. Клиент еще не соверщаш покупку и берет от двух до пяти книг - скидка 20% на весь счет
2. Клиент уже совершал покупки в магазине, берет от 2 до 5 книг - скидка 10%
3. Клиент уже совершал покупки в магазине, берет от 6 до 10 книг - скидка 5%
4. Клиент уже совершал покупки в магазине, берет свяше 10 книг - скидка 2%
5. Если сегодня пятница - всем скидка 5% не зависимо от объема корзины или типа клиента

// Вопросы:
// 1. Кстати, интересный момент: если покупатель первый раз покупает больше 5 разных книг, то скидки ему не будет?
// 5. Здесь имеется в виду, что доп скидка 5% или общая скидка будет 5%? Я возьмусь трактовать это, как сумму скидок.

*/

var (
	ErrEmptyCart        = errors.New("shopping cart is empty")
	ErrInvalidBookPrice = errors.New("book price must be positive")
)

type DiscountService struct {
	now func() time.Time
}

// NewDiscountService создает новый экземпляр DiscountService с правилами по умолчанию
func NewDiscountService() *DiscountService {
	return &DiscountService{
		now: time.Now,
	}
}

type bookAggregate struct {
	count int
	price float64
	total float64
}

// Calculate вычисляет скидку на основе правил
func (s *DiscountService) Calculate(ctx context.Context, customer domain.CustomerType, books []domain.Book) (domain.Discount, error) {
	if err := s.validateBooks(books); err != nil {
		return domain.Discount{}, err
	}

	var originalAmount float64

	// Группируем книги по ID для определения оптовых скидок
	booksByID := make(map[string]*bookAggregate)
	for _, b := range books {
		originalAmount += b.Price
		if agg, ok := booksByID[b.ID]; ok {
			agg.count++
			agg.total += b.Price
		} else {
			booksByID[b.ID] = &bookAggregate{count: 1, price: b.Price, total: b.Price}
		}
	}

	var (
		regularCount  int
		regularAmount float64
		bulkDiscount  float64
	)

	// Разделяем книги на оптовые и обычные
	for _, agg := range booksByID {
		if agg.count >= minBooksForBulkDiscount {
			discountedCopies := agg.count / 2 // Каждая вторая книга получает скидку 40%
			bulkDiscount += float64(discountedCopies) * agg.price * discountPercentBulkBook
		} else {
			regularAmount += agg.total
			regularCount += agg.count
		}
	}

	percent := s.discountPercent(customer, regularCount)
	if s.isFriday() {
		percent += discountPercentFriday
	}

	regularDiscount := regularAmount * percent
	discountAmount := regularDiscount + bulkDiscount
	finalAmount := originalAmount - discountAmount

	var totalPercent float64
	if originalAmount > 0 {
		totalPercent = discountAmount / originalAmount
	}

	return domain.Discount{
		CartAmount:      originalAmount,
		DiscountPercent: totalPercent,
		DiscountAmount:  discountAmount,
		TotalCost:       finalAmount,
	}, nil
}

func (s *DiscountService) discountPercent(customer domain.CustomerType, booksCount int) float64 {
	switch customer {
	case domain.CustomerTypeNew:
		if booksCount >= minBooksForNewCustomerDiscount && booksCount <= maxBooksForNewCustomerDiscount {
			return discountPercentNewCustomer2to5
		}
	case domain.CustomerTypeOld:
		switch {
		case booksCount >= minBooksForOldCustomerSmallDiscount && booksCount <= maxBooksForOldCustomerSmallDiscount:
			return discountPercentOldCustomer2to5
		case booksCount >= minBooksForOldCustomerMediumDiscount && booksCount <= maxBooksForOldCustomerMediumDiscount:
			return discountPercentOldCustomer6to10
		case booksCount >= minBooksForOldCustomerLargeDiscount:
			return discountPercentOldCustomerOver10
		}
	}
	return 0
}

func (s *DiscountService) isFriday() bool {
	return s.now().Weekday() == time.Friday
}

func (s *DiscountService) validateBooks(books []domain.Book) error {
	if len(books) == 0 {
		return ErrEmptyCart
	}

	for _, book := range books {
		if book.Price <= 0 {
			return ErrInvalidBookPrice
		}
	}

	return nil
}
