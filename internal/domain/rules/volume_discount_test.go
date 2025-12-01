package rules

import (
	"testing"
	"time"

	"github.com/Badgain/book-discount/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVolumeDiscountRule_Name(t *testing.T) {
	rule := NewVolumeDiscountRule([]VolumeRange{})
	assert.Equal(t, "VolumeDiscountRule", rule.Name())
}

func TestVolumeDiscountRule_BlocksOtherRules(t *testing.T) {
	rule := NewVolumeDiscountRule([]VolumeRange{})
	assert.False(t, rule.BlocksOtherRules())
}

func TestVolumeDiscountRule_CanApply(t *testing.T) {
	maxBooks5 := 5
	maxBooks10 := 10

	ranges := []VolumeRange{
		{
			CustomerType: domain.CustomerTypeNew,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 20,
		},
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 10,
		},
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     6,
			MaxBooks:     &maxBooks10,
			DiscountRate: 5,
		},
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     11,
			MaxBooks:     nil,
			DiscountRate: 2,
		},
	}

	rule := NewVolumeDiscountRule(ranges)

	tests := []struct {
		name         string
		customerType domain.CustomerType
		bookCount    int
		want         bool
	}{
		{
			name:         "new customer with 3 books (in range 2-5)",
			customerType: domain.CustomerTypeNew,
			bookCount:    3,
			want:         true,
		},
		{
			name:         "new customer with 1 book (below range)",
			customerType: domain.CustomerTypeNew,
			bookCount:    1,
			want:         false,
		},
		{
			name:         "new customer with 6 books (above range)",
			customerType: domain.CustomerTypeNew,
			bookCount:    6,
			want:         false,
		},
		{
			name:         "old customer with 3 books (in range 2-5)",
			customerType: domain.CustomerTypeOld,
			bookCount:    3,
			want:         true,
		},
		{
			name:         "old customer with 7 books (in range 6-10)",
			customerType: domain.CustomerTypeOld,
			bookCount:    7,
			want:         true,
		},
		{
			name:         "old customer with 15 books (in range 11+)",
			customerType: domain.CustomerTypeOld,
			bookCount:    15,
			want:         true,
		},
		{
			name:         "old customer with 1 book (below all ranges)",
			customerType: domain.CustomerTypeOld,
			bookCount:    1,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := make([]domain.Book, tt.bookCount)
			for i := 0; i < tt.bookCount; i++ {
				books[i] = domain.Book{ID: "book", Price: 1000}
			}

			ruleCtx := domain.RuleContext{
				CustomerType: tt.customerType,
				Books:        books,
				Time:         time.Now(),
			}

			got := rule.CanApply(ruleCtx)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVolumeDiscountRule_Apply_Success(t *testing.T) {
	maxBooks5 := 5

	ranges := []VolumeRange{
		{
			CustomerType: domain.CustomerTypeNew,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 20,
		},
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 10,
		},
	}

	rule := NewVolumeDiscountRule(ranges)

	tests := []struct {
		name             string
		customerType     domain.CustomerType
		books            []domain.Book
		wantDiscount     int64
		wantAppliedCount int
	}{
		{
			name:         "new customer, 3 books, 20% discount",
			customerType: domain.CustomerTypeNew,
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "2", Price: 1000},
				{ID: "3", Price: 1000},
			},
			wantDiscount:     600,
			wantAppliedCount: 3,
		},
		{
			name:         "old customer, 4 books, 10% discount",
			customerType: domain.CustomerTypeOld,
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "2", Price: 1000},
				{ID: "3", Price: 1000},
				{ID: "4", Price: 1000},
			},
			wantDiscount:     400,
			wantAppliedCount: 4,
		},
		{
			name:         "different prices",
			customerType: domain.CustomerTypeNew,
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "2", Price: 2000},
				{ID: "3", Price: 1500},
			},
			wantDiscount:     900,
			wantAppliedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ruleCtx := domain.RuleContext{
				CustomerType: tt.customerType,
				Books:        tt.books,
				Time:         time.Now(),
			}

			result, err := rule.Apply(ruleCtx)

			require.NoError(t, err)
			assert.Equal(t, tt.wantDiscount, result.DiscountAmount)
			assert.Equal(t, tt.wantAppliedCount, len(result.AppliedBooks))
			assert.Equal(t, 0, len(result.RemainingBooks))
			assert.Equal(t, "VolumeDiscountRule", result.RuleName)
		})
	}
}

func TestVolumeDiscountRule_Apply_ValidationErrors(t *testing.T) {
	maxBooks5 := 5

	ranges := []VolumeRange{
		{
			CustomerType: domain.CustomerTypeNew,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 20,
		},
	}

	rule := NewVolumeDiscountRule(ranges)

	tests := []struct {
		name          string
		customerType  domain.CustomerType
		books         []domain.Book
		wantErrSubstr string
	}{
		{
			name:          "nil books",
			customerType:  domain.CustomerTypeNew,
			books:         nil,
			wantErrSubstr: "cannot be nil",
		},
		{
			name:         "empty book ID",
			customerType: domain.CustomerTypeNew,
			books: []domain.Book{
				{ID: "", Price: 1000},
				{ID: "2", Price: 1000},
				{ID: "3", Price: 1000},
			},
			wantErrSubstr: "ID cannot be empty",
		},
		{
			name:         "negative price",
			customerType: domain.CustomerTypeNew,
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "2", Price: -100},
				{ID: "3", Price: 1000},
			},
			wantErrSubstr: "price cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ruleCtx := domain.RuleContext{
				CustomerType: tt.customerType,
				Books:        tt.books,
				Time:         time.Now(),
			}

			result, err := rule.Apply(ruleCtx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrSubstr)
			assert.Equal(t, int64(0), result.DiscountAmount)
		})
	}
}

func TestVolumeDiscountRule_Apply_NoMatchingRange(t *testing.T) {
	maxBooks5 := 5

	ranges := []VolumeRange{
		{
			CustomerType: domain.CustomerTypeNew,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 20,
		},
	}

	rule := NewVolumeDiscountRule(ranges)

	ruleCtx := domain.RuleContext{
		CustomerType: domain.CustomerTypeNew,
		Books: []domain.Book{
			{ID: "1", Price: 1000},
		},
		Time: time.Now(),
	}

	result, err := rule.Apply(ruleCtx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching discount range")
	assert.Equal(t, int64(0), result.DiscountAmount)
}

func TestVolumeDiscountRule_getDiscountRate(t *testing.T) {
	maxBooks5 := 5
	maxBooks10 := 10

	ranges := []VolumeRange{
		{
			CustomerType: domain.CustomerTypeNew,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 20,
		},
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     2,
			MaxBooks:     &maxBooks5,
			DiscountRate: 10,
		},
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     6,
			MaxBooks:     &maxBooks10,
			DiscountRate: 5,
		},
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     11,
			MaxBooks:     nil,
			DiscountRate: 2,
		},
	}

	rule := NewVolumeDiscountRule(ranges)

	tests := []struct {
		name         string
		customerType domain.CustomerType
		bookCount    int
		wantRate     int
	}{
		{
			name:         "new customer, 3 books → 20%",
			customerType: domain.CustomerTypeNew,
			bookCount:    3,
			wantRate:     20,
		},
		{
			name:         "new customer, 1 book → 0%",
			customerType: domain.CustomerTypeNew,
			bookCount:    1,
			wantRate:     0,
		},
		{
			name:         "old customer, 3 books → 10%",
			customerType: domain.CustomerTypeOld,
			bookCount:    3,
			wantRate:     10,
		},
		{
			name:         "old customer, 7 books → 5%",
			customerType: domain.CustomerTypeOld,
			bookCount:    7,
			wantRate:     5,
		},
		{
			name:         "old customer, 15 books → 2%",
			customerType: domain.CustomerTypeOld,
			bookCount:    15,
			wantRate:     2,
		},
		{
			name:         "old customer, 1 book → 0%",
			customerType: domain.CustomerTypeOld,
			bookCount:    1,
			wantRate:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.getDiscountRate(tt.customerType, tt.bookCount)
			assert.Equal(t, tt.wantRate, got)
		})
	}
}

func TestVolumeDiscountRule_UnlimitedMaxBooks(t *testing.T) {
	ranges := []VolumeRange{
		{
			CustomerType: domain.CustomerTypeOld,
			MinBooks:     10,
			MaxBooks:     nil,
			DiscountRate: 2,
		},
	}

	rule := NewVolumeDiscountRule(ranges)

	tests := []struct {
		name      string
		bookCount int
		want      bool
	}{
		{
			name:      "9 books (below min) → false",
			bookCount: 9,
			want:      false,
		},
		{
			name:      "10 books (at min) → true",
			bookCount: 10,
			want:      true,
		},
		{
			name:      "100 books (way above min) → true",
			bookCount: 100,
			want:      true,
		},
		{
			name:      "1000 books (unlimited) → true",
			bookCount: 1000,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := make([]domain.Book, tt.bookCount)
			for i := 0; i < tt.bookCount; i++ {
				books[i] = domain.Book{ID: "book", Price: 100}
			}

			ruleCtx := domain.RuleContext{
				CustomerType: domain.CustomerTypeOld,
				Books:        books,
				Time:         time.Now(),
			}

			got := rule.CanApply(ruleCtx)
			assert.Equal(t, tt.want, got)

			if got {
				result, err := rule.Apply(ruleCtx)
				require.NoError(t, err)
				expectedDiscount := int64(tt.bookCount) * 100 * 2 / 100
				assert.Equal(t, expectedDiscount, result.DiscountAmount)
			}
		})
	}
}
