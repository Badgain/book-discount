package dto

import (
	"math"

	"github.com/Badgain/book-discount/internal/domain"
)

type Book struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

type DiscountRequest struct {
	CustomerID     string `json:"customer_id"`
	CustomerType   string `json:"customer_type"`
	CashRegisterID string `json:"cash_register_id"`
	Books          []Book `json:"books"`
}

type DiscountResponse struct {
	OriginalAmount  float64 `json:"original_amount"`
	DiscountPercent float64 `json:"discount_percent"`
	DiscountAmount  float64 `json:"discount_amount"`
	FinalAmount     float64 `json:"final_amount"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (d *DiscountRequest) CustomerTypeAsDomain() domain.CustomerType {
	return domain.CustomerType(d.CustomerType)
}

func (d *DiscountRequest) BooksAsDomain() []domain.Book {
	books := make([]domain.Book, len(d.Books))
	for i, book := range d.Books {
		books[i] = domain.Book{
			ID:    book.ID,
			Price: floatToCents(book.Price),
		}
	}
	return books
}

func NewDiscountResponse(discount domain.Discount) DiscountResponse {
	return DiscountResponse{
		OriginalAmount:  centsToFloat(discount.CartAmount),
		DiscountPercent: discount.DiscountPercent,
		DiscountAmount:  centsToFloat(discount.DiscountAmount),
		FinalAmount:     centsToFloat(discount.TotalCost),
	}
}

func floatToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func centsToFloat(cents int64) float64 {
	return float64(cents) / 100.0
}
