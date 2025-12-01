package rules

import (
	"testing"
	"time"

	"github.com/Badgain/book-discount/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFridayRule_Name(t *testing.T) {
	rule := NewFridayRule()
	assert.Equal(t, "FridayRule", rule.Name())
}

func TestFridayRule_BlocksOtherRules(t *testing.T) {
	rule := NewFridayRule()
	assert.True(t, rule.BlocksOtherRules())
}

func TestFridayRule_CanApply(t *testing.T) {
	tests := []struct {
		name     string
		weekday  time.Weekday
		expected bool
	}{
		{
			name:     "monday",
			weekday:  time.Monday,
			expected: false,
		},
		{
			name:     "tuesday",
			weekday:  time.Tuesday,
			expected: false,
		},
		{
			name:     "wednesday",
			weekday:  time.Wednesday,
			expected: false,
		},
		{
			name:     "thursday",
			weekday:  time.Thursday,
			expected: false,
		},
		{
			name:     "friday",
			weekday:  time.Friday,
			expected: true,
		},
		{
			name:     "saturday",
			weekday:  time.Saturday,
			expected: false,
		},
		{
			name:     "sunday",
			weekday:  time.Sunday,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewFridayRule()
			testTime := getDateForWeekday(tt.weekday)
			ruleCtx := domain.RuleContext{
				CustomerType: domain.CustomerTypeNew,
				Books: []domain.Book{
					{ID: "1", Price: 1000},
				},
				Time: testTime,
			}
			assert.Equal(t, tt.expected, rule.CanApply(ruleCtx))
		})
	}
}

func TestFridayRule_Apply_Success(t *testing.T) {
	tests := []struct {
		name             string
		books            []domain.Book
		expectedDiscount int64
	}{
		{
			name: "single book",
			books: []domain.Book{
				{ID: "1", Price: 1000},
			},
			expectedDiscount: 50,
		},
		{
			name: "multiple books same price",
			books: []domain.Book{
				{ID: "1", Price: 1000},
				{ID: "2", Price: 1000},
				{ID: "3", Price: 1000},
			},
			expectedDiscount: 150,
		},
		{
			name: "multiple books different prices",
			books: []domain.Book{
				{ID: "1", Price: 2000},
				{ID: "2", Price: 1500},
				{ID: "3", Price: 1000},
			},
			expectedDiscount: 225,
		},
		{
			name: "many books",
			books: []domain.Book{
				{ID: "1", Price: 1299},
				{ID: "2", Price: 1299},
				{ID: "3", Price: 1299},
				{ID: "4", Price: 500},
				{ID: "5", Price: 500},
			},
			expectedDiscount: 242,
		},
		{
			name:             "empty cart",
			books:            []domain.Book{},
			expectedDiscount: 0,
		},
	}

	friday := getDateForWeekday(time.Friday)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewFridayRule()
			ruleCtx := domain.RuleContext{
				CustomerType: domain.CustomerTypeNew,
				Books:        tt.books,
				Time:         friday,
			}

			result, err := rule.Apply(ruleCtx)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedDiscount, result.DiscountAmount)
			assert.Equal(t, len(tt.books), len(result.AppliedBooks))
			assert.Equal(t, 0, len(result.RemainingBooks))
			assert.Equal(t, "FridayRule", result.RuleName)
		})
	}
}

func TestFridayRule_Apply_ValidationErrors(t *testing.T) {
	friday := getDateForWeekday(time.Friday)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewFridayRule()
			ruleCtx := domain.RuleContext{
				CustomerType: domain.CustomerTypeNew,
				Books:        tt.books,
				Time:         friday,
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

func getDateForWeekday(weekday time.Weekday) time.Time {
	baseDate := time.Date(2025, 1, 13, 12, 0, 0, 0, time.UTC)

	currentWeekday := baseDate.Weekday()
	daysToAdd := int(weekday - currentWeekday)

	return baseDate.AddDate(0, 0, daysToAdd)
}
