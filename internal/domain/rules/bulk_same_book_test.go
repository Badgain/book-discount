package rules

import (
	"testing"
	"time"

	"github.com/Badgain/book-discount/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

func TestBulkSameBookRule_Name(t *testing.T) {
	rule := NewBulkSameBookRule()
	assert.Equal(t, "BulkSameBookRule", rule.Name())
}

func TestBulkSameBookRule_BlocksOtherRules(t *testing.T) {
	rule := NewBulkSameBookRule()
	assert.False(t, rule.BlocksOtherRules())
}

func TestBulkSameBookRule_CanApply(t *testing.T) {
	tests := []struct {
		name     string
		books    []domain.Book
		expected bool
	}{
		{
			name:     "empty books",
			books:    []domain.Book{},
			expected: false,
		},
		{
			name: "less than required",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
			},
			expected: false,
		},
		{
			name: "exactly at threshold",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
			},
			expected: false,
		},
		{
			name: "above threshold",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
			},
			expected: true,
		},
		{
			name: "multiple groups only one qualifies",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewBulkSameBookRule()
			ruleCtx := domain.RuleContext{
				CustomerType: domain.CustomerTypeNew,
				Books:        tt.books,
				Time:         fixedTime,
			}
			assert.Equal(t, tt.expected, rule.CanApply(ruleCtx))
		})
	}
}

func TestBulkSameBookRule_Apply_Success(t *testing.T) {
	tests := []struct {
		name                 string
		books                []domain.Book
		expectedDiscount     int64
		expectedAppliedCount int
		expectedRemainCount  int
	}{
		{
			name: "six books - three with 40 percent",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
			},
			expectedDiscount:     1200,
			expectedAppliedCount: 6,
			expectedRemainCount:  0,
		},
		{
			name: "seven books - three with 40 percent",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
			},
			expectedDiscount:     1200,
			expectedAppliedCount: 7,
			expectedRemainCount:  0,
		},
		{
			name: "ten books - five with 40 percent",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
			},
			expectedDiscount:     2000,
			expectedAppliedCount: 10,
			expectedRemainCount:  0,
		},
		{
			name: "multiple groups",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
				{ID: "3", Price: 800},
				{ID: "3", Price: 800},
				{ID: "3", Price: 800},
				{ID: "3", Price: 800},
				{ID: "3", Price: 800},
				{ID: "3", Price: 800},
				{ID: "3", Price: 800},
			},
			expectedDiscount:     2160,
			expectedAppliedCount: 13,
			expectedRemainCount:  3,
		},
		{
			name: "no qualifying groups",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "1", Price: 1000},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
				{ID: "2", Price: 500},
			},
			expectedDiscount:     0,
			expectedAppliedCount: 0,
			expectedRemainCount:  5,
		},
		{
			name:                 "empty books",
			books:                []domain.Book{},
			expectedDiscount:     0,
			expectedAppliedCount: 0,
			expectedRemainCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewBulkSameBookRule()
			ruleCtx := domain.RuleContext{
				CustomerType: domain.CustomerTypeNew,
				Books:        tt.books,
				Time:         fixedTime,
			}

			result, err := rule.Apply(ruleCtx)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedDiscount, result.DiscountAmount)
			assert.Equal(t, tt.expectedAppliedCount, len(result.AppliedBooks))
			assert.Equal(t, tt.expectedRemainCount, len(result.RemainingBooks))
			assert.Equal(t, "BulkSameBookRule", result.RuleName)
		})
	}
}

func TestBulkSameBookRule_Apply_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		books       []domain.Book
		expectError bool
	}{
		{
			name:        "nil books",
			books:       nil,
			expectError: true,
		},
		{
			name: "empty book ID",
			books: []domain.Book{
				{ID: "", Price: 1000},
			},
			expectError: true,
		},
		{
			name: "negative price",
			books: []domain.Book{
				{ID: "1", Price: -100},
			},
			expectError: true,
		},
		{
			name: "zero price is valid",
			books: []domain.Book{
				{ID: "1", Price: 0},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewBulkSameBookRule()
			ruleCtx := domain.RuleContext{
				CustomerType: domain.CustomerTypeNew,
				Books:        tt.books,
				Time:         fixedTime,
			}

			_, err := rule.Apply(ruleCtx)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBulkSameBookRule_Apply_Determinism(t *testing.T) {
	books := []domain.Book{
		{ID: "z", Price: 1000},
		{ID: "z", Price: 1000},
		{ID: "z", Price: 1000},
		{ID: "z", Price: 1000},
		{ID: "z", Price: 1000},
		{ID: "z", Price: 1000},
		{ID: "a", Price: 500},
		{ID: "a", Price: 500},
		{ID: "a", Price: 500},
		{ID: "a", Price: 500},
		{ID: "a", Price: 500},
		{ID: "a", Price: 500},
		{ID: "m", Price: 800},
		{ID: "m", Price: 800},
	}

	rule := NewBulkSameBookRule()
	ruleCtx := domain.RuleContext{
		CustomerType: domain.CustomerTypeNew,
		Books:        books,
		Time:         fixedTime,
	}

	var firstResult domain.RuleResult
	for i := 0; i < 10; i++ {
		result, err := rule.Apply(ruleCtx)
		require.NoError(t, err)

		if i == 0 {
			firstResult = result
		} else {
			assert.Equal(t, firstResult.DiscountAmount, result.DiscountAmount)
			assert.Equal(t, len(firstResult.AppliedBooks), len(result.AppliedBooks))
			assert.Equal(t, len(firstResult.RemainingBooks), len(result.RemainingBooks))
		}
	}
}
