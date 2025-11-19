package dto

import "github.com/Badgain/book-discount/internal/domain"

// Book представляет книгу в корзине
type (
	Book struct {
		ID    string  `json:"id"`
		Price float64 `json:"price"`
	}

	DiscountRequest struct {
		CustomerID     string `json:"customer_id"`
		CustomerType   string `json:"customer_type"`
		CashRegisterID string `json:"cash_register_id"`
		Books          []Book `json:"books"`
	}

	// DiscountResponse представляет ответ с расчетом скидки
	DiscountResponse struct {
		OriginalAmount  float64 `json:"original_amount"`
		DiscountPercent float64 `json:"discount_percent"`
		DiscountAmount  float64 `json:"discount_amount"`
		FinalAmount     float64 `json:"final_amount"`
	}
)

func (d *DiscountRequest) CustomerTypeAsDomain() domain.CustomerType {
	return domain.CustomerType(d.CustomerType)
}

func (d *DiscountRequest) BooksAsDomain() []domain.Book {
	books := make([]domain.Book, len(d.Books))
	for i, book := range d.Books {
		books[i] = domain.Book{ID: book.ID, Price: book.Price}
	}
	return books
}
