package service

import (
	"context"
	"testing"
	"time"

	"github.com/Badgain/book-discount/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockTime(weekday time.Weekday) func() time.Time {
	return func() time.Time {
		baseDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		daysToAdd := int(weekday - time.Monday)
		if daysToAdd < 0 {
			daysToAdd += 7
		}
		return baseDate.AddDate(0, 0, daysToAdd)
	}
}

func TestDiscountService_Calculate_Validation(t *testing.T) {
	svc := NewDiscountService()
	ctx := context.Background()

	tests := []struct {
		name    string
		books   []domain.Book
		wantErr error
	}{
		{
			name:    "empty cart",
			books:   []domain.Book{},
			wantErr: ErrEmptyCart,
		},
		{
			name: "zero price",
			books: []domain.Book{
				{ID: "1", Price: 0},
			},
			wantErr: ErrInvalidBookPrice,
		},
		{
			name: "negative price",
			books: []domain.Book{
				{ID: "1", Price: -10},
			},
			wantErr: ErrInvalidBookPrice,
		},
		{
			name: "valid books",
			books: []domain.Book{
				{ID: "1", Price: 10},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Calculate(ctx, domain.CustomerTypeNew, tt.books)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDiscountService_Calculate_NewCustomer(t *testing.T) {
	svc := NewDiscountService()
	svc.now = mockTime(time.Monday)
	ctx := context.Background()

	tests := []struct {
		name         string
		books        []domain.Book
		wantPercent  float64
		wantDiscount float64
		wantFinal    float64
	}{
		{
			name: "1 book - no discount",
			books: []domain.Book{
				{ID: "1", Price: 100},
			},
			wantPercent:  0,
			wantDiscount: 0,
			wantFinal:    100,
		},
		{
			name: "2 books - 20% discount",
			books: []domain.Book{
				{ID: "1", Price: 100},
				{ID: "2", Price: 50},
			},
			wantPercent:  0.20,
			wantDiscount: 30,
			wantFinal:    120,
		},
		{
			name: "3 books - 20% discount",
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
			},
			wantPercent:  0.20,
			wantDiscount: 6,
			wantFinal:    24,
		},
		{
			name: "5 books - 20% discount",
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
				{ID: "4", Price: 10},
				{ID: "5", Price: 10},
			},
			wantPercent:  0.20,
			wantDiscount: 10,
			wantFinal:    40,
		},
		{
			name: "6 books - no discount (over limit)",
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
				{ID: "4", Price: 10},
				{ID: "5", Price: 10},
				{ID: "6", Price: 10},
			},
			wantPercent:  0,
			wantDiscount: 0,
			wantFinal:    60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Calculate(ctx, domain.CustomerTypeNew, tt.books)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantPercent, result.DiscountPercent, 0.001)
			assert.InDelta(t, tt.wantDiscount, result.DiscountAmount, 0.01)
			assert.InDelta(t, tt.wantFinal, result.TotalCost, 0.01)
		})
	}
}

func TestDiscountService_Calculate_OldCustomer(t *testing.T) {
	svc := NewDiscountService()
	svc.now = mockTime(time.Monday)
	ctx := context.Background()

	tests := []struct {
		name         string
		books        []domain.Book
		wantPercent  float64
		wantDiscount float64
		wantFinal    float64
	}{
		{
			name: "1 book - no discount",
			books: []domain.Book{
				{ID: "1", Price: 100},
			},
			wantPercent:  0,
			wantDiscount: 0,
			wantFinal:    100,
		},
		{
			name: "2 books - 10% discount",
			books: []domain.Book{
				{ID: "1", Price: 100},
				{ID: "2", Price: 50},
			},
			wantPercent:  0.10,
			wantDiscount: 15,
			wantFinal:    135,
		},
		{
			name: "5 books - 10% discount",
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
				{ID: "4", Price: 10},
				{ID: "5", Price: 10},
			},
			wantPercent:  0.10,
			wantDiscount: 5,
			wantFinal:    45,
		},
		{
			name: "6 books - 5% discount",
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
				{ID: "4", Price: 10},
				{ID: "5", Price: 10},
				{ID: "6", Price: 10},
			},
			wantPercent:  0.05,
			wantDiscount: 3,
			wantFinal:    57,
		},
		{
			name: "10 books - 5% discount",
			books: func() []domain.Book {
				var books []domain.Book
				for i := 0; i < 10; i++ {
					books = append(books, domain.Book{ID: string(rune('a' + i)), Price: 10})
				}
				return books
			}(),
			wantPercent:  0.05,
			wantDiscount: 5,
			wantFinal:    95,
		},
		{
			name: "11 books - 2% discount",
			books: func() []domain.Book {
				var books []domain.Book
				for i := 0; i < 11; i++ {
					books = append(books, domain.Book{ID: string(rune('a' + i)), Price: 10})
				}
				return books
			}(),
			wantPercent:  0.02,
			wantDiscount: 2.2,
			wantFinal:    107.8,
		},
		{
			name: "15 books - 2% discount",
			books: func() []domain.Book {
				var books []domain.Book
				for i := 0; i < 15; i++ {
					books = append(books, domain.Book{ID: string(rune('a' + i)), Price: 10})
				}
				return books
			}(),
			wantPercent:  0.02,
			wantDiscount: 3,
			wantFinal:    147,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Calculate(ctx, domain.CustomerTypeOld, tt.books)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantPercent, result.DiscountPercent, 0.001)
			assert.InDelta(t, tt.wantDiscount, result.DiscountAmount, 0.01)
			assert.InDelta(t, tt.wantFinal, result.TotalCost, 0.01)
		})
	}
}

func TestDiscountService_Calculate_FridayDiscount(t *testing.T) {
	svc := NewDiscountService()
	svc.now = mockTime(time.Friday) // пятница
	ctx := context.Background()

	tests := []struct {
		name         string
		customer     domain.CustomerType
		books        []domain.Book
		wantPercent  float64
		wantDiscount float64
		wantFinal    float64
	}{
		{
			name:     "new customer 1 book + friday 5%",
			customer: domain.CustomerTypeNew,
			books: []domain.Book{
				{ID: "1", Price: 100},
			},
			wantPercent:  0.05,
			wantDiscount: 5,
			wantFinal:    95,
		},
		{
			name:     "new customer 3 books + friday 5% (20% + 5%)",
			customer: domain.CustomerTypeNew,
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
			},
			wantPercent:  0.25, // 20% + 5% = 25%
			wantDiscount: 7.5,  // 30 * 0.25 = 7.5
			wantFinal:    22.5,
		},
		{
			name:     "old customer 1 book + friday 5%",
			customer: domain.CustomerTypeOld,
			books: []domain.Book{
				{ID: "1", Price: 100},
			},
			wantPercent:  0.05,
			wantDiscount: 5,
			wantFinal:    95,
		},
		{
			name:     "old customer 3 books + friday 5% (10% + 5%)",
			customer: domain.CustomerTypeOld,
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
			},
			wantPercent:  0.15, // 10% + 5% = 15%
			wantDiscount: 4.5,  // 30 * 0.15 = 4.5
			wantFinal:    25.5,
		},
		{
			name:     "old customer 7 books + friday 5% (5% + 5%)",
			customer: domain.CustomerTypeOld,
			books: func() []domain.Book {
				var books []domain.Book
				for i := 0; i < 7; i++ {
					books = append(books, domain.Book{ID: string(rune('a' + i)), Price: 10})
				}
				return books
			}(),
			wantPercent:  0.10, // 5% + 5% = 10%
			wantDiscount: 7,    // 70 * 0.10 = 7
			wantFinal:    63,
		},
		{
			name:     "old customer 11 books + friday 5% (2% + 5%)",
			customer: domain.CustomerTypeOld,
			books: func() []domain.Book {
				var books []domain.Book
				for i := 0; i < 11; i++ {
					books = append(books, domain.Book{ID: string(rune('a' + i)), Price: 10})
				}
				return books
			}(),
			wantPercent:  0.07, // 2% + 5% = 7%
			wantDiscount: 7.7,  // 110 * 0.07 = 7.7
			wantFinal:    102.3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Calculate(ctx, tt.customer, tt.books)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantPercent, result.DiscountPercent, 0.001)
			assert.InDelta(t, tt.wantDiscount, result.DiscountAmount, 0.01)
			assert.InDelta(t, tt.wantFinal, result.TotalCost, 0.01)
		})
	}
}

func TestDiscountService_Calculate_BulkDiscount(t *testing.T) {
	svc := NewDiscountService()
	svc.now = mockTime(time.Monday)
	ctx := context.Background()

	tests := []struct {
		name         string
		customer     domain.CustomerType
		books        []domain.Book
		wantPercent  float64
		wantDiscount float64
		wantFinal    float64
	}{
		{
			name:     "6 identical books - 3 books with 40% discount",
			customer: domain.CustomerTypeOld,
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
			},
			wantPercent:  0.20, // 3 книги * 10 * 0.4 = 12, всего 60, процент = 12/60 = 0.2
			wantDiscount: 12,   // 3 книги со скидкой 40%
			wantFinal:    48,
		},
		{
			name:     "7 identical books - 3 books with 40% discount",
			customer: domain.CustomerTypeOld,
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
				{ID: "1", Price: 10},
			},
			wantPercent:  0.17,
			wantDiscount: 12,
			wantFinal:    58,
		},
		{
			name:     "10 identical books - 5 books with 40% discount",
			customer: domain.CustomerTypeOld,
			books: func() []domain.Book {
				var books []domain.Book
				for i := 0; i < 10; i++ {
					books = append(books, domain.Book{ID: "1", Price: 10})
				}
				return books
			}(),
			wantPercent:  0.20,
			wantDiscount: 20,
			wantFinal:    80,
		},
		{
			name:     "bulk books don't count for regular discount",
			customer: domain.CustomerTypeOld,
			books: []domain.Book{
				// 6 одинаковых книг (оптовая скидка, не считаются в regularCount)
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				// 2 обычные книги (должны получить 10% скидку)
				{ID: "regular1", Price: 10},
				{ID: "regular2", Price: 10},
			},
			wantPercent:  0.18,
			wantDiscount: 14,
			wantFinal:    66,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Calculate(ctx, tt.customer, tt.books)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantPercent, result.DiscountPercent, 0.001)
			assert.InDelta(t, tt.wantDiscount, result.DiscountAmount, 0.01)
			assert.InDelta(t, tt.wantFinal, result.TotalCost, 0.01)
		})
	}
}

func TestDiscountService_Calculate_CombinedDiscounts(t *testing.T) {
	svc := NewDiscountService()
	svc.now = mockTime(time.Friday) // пятница
	ctx := context.Background()

	tests := []struct {
		name         string
		customer     domain.CustomerType
		books        []domain.Book
		wantPercent  float64
		wantDiscount float64
		wantFinal    float64
	}{
		{
			name:     "bulk + regular + friday",
			customer: domain.CustomerTypeOld,
			books: []domain.Book{
				// 6 одинаковых книг (оптовая скидка)
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				// 3 обычные книги (10% + 5% пятница = 15%)
				{ID: "regular1", Price: 10},
				{ID: "regular2", Price: 10},
				{ID: "regular3", Price: 10},
			},
			wantPercent:  0.18,
			wantDiscount: 16.5,
			wantFinal:    73.5,
		},
		{
			name:     "new customer bulk + regular + friday",
			customer: domain.CustomerTypeNew,
			books: []domain.Book{
				// 6 одинаковых книг (оптовая скидка)
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				{ID: "bulk", Price: 10},
				// 2 обычные книги (20% + 5% пятница = 25%)
				{ID: "regular1", Price: 10},
				{ID: "regular2", Price: 10},
			},
			wantPercent:  0.21,
			wantDiscount: 17,
			wantFinal:    63,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Calculate(ctx, tt.customer, tt.books)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantPercent, result.DiscountPercent, 0.001)
			assert.InDelta(t, tt.wantDiscount, result.DiscountAmount, 0.01)
			assert.InDelta(t, tt.wantFinal, result.TotalCost, 0.01)
		})
	}
}

func TestDiscountService_Calculate_EdgeCases(t *testing.T) {
	svc := NewDiscountService()
	svc.now = mockTime(time.Monday)
	ctx := context.Background()

	tests := []struct {
		name         string
		customer     domain.CustomerType
		books        []domain.Book
		wantPercent  float64
		wantDiscount float64
		wantFinal    float64
	}{
		{
			name:     "multiple bulk book groups",
			customer: domain.CustomerTypeOld,
			books: []domain.Book{
				// 6 книг типа A
				{ID: "A", Price: 10},
				{ID: "A", Price: 10},
				{ID: "A", Price: 10},
				{ID: "A", Price: 10},
				{ID: "A", Price: 10},
				{ID: "A", Price: 10},
				// 6 книг типа B
				{ID: "B", Price: 20},
				{ID: "B", Price: 20},
				{ID: "B", Price: 20},
				{ID: "B", Price: 20},
				{ID: "B", Price: 20},
				{ID: "B", Price: 20},
			},
			wantPercent:  0.20,
			wantDiscount: 36,
			wantFinal:    144,
		},
		{
			name:     "exactly 5 books for new customer",
			customer: domain.CustomerTypeNew,
			books: []domain.Book{
				{ID: "1", Price: 10},
				{ID: "2", Price: 10},
				{ID: "3", Price: 10},
				{ID: "4", Price: 10},
				{ID: "5", Price: 10},
			},
			wantPercent:  0.20,
			wantDiscount: 10,
			wantFinal:    40,
		},
		{
			name:     "exactly 10 books for old customer",
			customer: domain.CustomerTypeOld,
			books: func() []domain.Book {
				var books []domain.Book
				for i := 0; i < 10; i++ {
					books = append(books, domain.Book{ID: string(rune('a' + i)), Price: 10})
				}
				return books
			}(),
			wantPercent:  0.05,
			wantDiscount: 5,
			wantFinal:    95,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Calculate(ctx, tt.customer, tt.books)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantPercent, result.DiscountPercent, 0.001)
			assert.InDelta(t, tt.wantDiscount, result.DiscountAmount, 0.01)
			assert.InDelta(t, tt.wantFinal, result.TotalCost, 0.01)
		})
	}
}
