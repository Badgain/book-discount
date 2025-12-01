package service

import (
	"context"
	"testing"
	"time"

	"github.com/Badgain/book-discount/internal/config"
	"github.com/Badgain/book-discount/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleEngine_EmptyCart(t *testing.T) {
	engine := createTestEngine(time.Monday)

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, []domain.Book{})

	require.NoError(t, err)
	assert.Equal(t, int64(0), discount.CartAmount)
	assert.Equal(t, int64(0), discount.DiscountAmount)
	assert.Equal(t, int64(0), discount.TotalCost)
	assert.Equal(t, 0.0, discount.DiscountPercent)
}

func TestRuleEngine_NewCustomer_ThreeBooks_20Percent(t *testing.T) {
	engine := createTestEngine(time.Monday)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "2", Price: 1000},
		{ID: "3", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(3000), discount.CartAmount)
	assert.Equal(t, int64(600), discount.DiscountAmount)
	assert.Equal(t, int64(2400), discount.TotalCost)
	assert.Equal(t, 0.2, discount.DiscountPercent)
}

func TestRuleEngine_Friday_BlocksCustomerVolumeRule(t *testing.T) {
	engine := createTestEngine(time.Friday)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "2", Price: 1000},
		{ID: "3", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(3000), discount.CartAmount)
	assert.Equal(t, int64(150), discount.DiscountAmount)
	assert.Equal(t, int64(2850), discount.TotalCost)
	assert.InDelta(t, 0.05, discount.DiscountPercent, 0.001)
}

func TestRuleEngine_BulkBooks_Plus_RemainingBooks_NewCustomer(t *testing.T) {
	engine := createTestEngine(time.Monday)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "2", Price: 500},
		{ID: "3", Price: 500},
		{ID: "4", Price: 500},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(7500), discount.CartAmount)
	assert.Equal(t, int64(1500), discount.DiscountAmount)
	assert.Equal(t, int64(6000), discount.TotalCost)
	assert.Equal(t, 0.2, discount.DiscountPercent)
}

func TestRuleEngine_Friday_BulkBooks_Plus_RemainingBooks(t *testing.T) {
	engine := createTestEngine(time.Friday)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "2", Price: 500},
		{ID: "3", Price: 500},
		{ID: "4", Price: 500},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(7500), discount.CartAmount)
	assert.Equal(t, int64(1275), discount.DiscountAmount)
	assert.Equal(t, int64(6225), discount.TotalCost)
	assert.Equal(t, 0.17, discount.DiscountPercent)
}

func TestRuleEngine_CustomConfig_DisabledRule(t *testing.T) {
	cfg := &config.DiscountConfig{
		Rules: []config.RuleConfig{
			{
				Name:     "BulkSameBookRule",
				Enabled:  false,
				Priority: 1,
				Params: map[string]any{
					"minBooks":     5,
					"discountRate": 40,
				},
			},
			{
				Name:     "VolumeDiscountRule",
				Enabled:  true,
				Priority: 20,
				Params: map[string]any{
					"ranges": []any{
						map[string]any{
							"customerType": "new",
							"minBooks":     2,
							"maxBooks":     5,
							"discountRate": 20,
						},
					},
				},
			},
		},
	}

	timeProvider := &mockTimeProvider{
		currentTime: getDateForWeekday(time.Monday),
	}
	engine, err := NewRuleEngine(cfg, timeProvider)
	require.NoError(t, err)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
		{ID: "1", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(6000), discount.CartAmount)
	assert.Equal(t, int64(0), discount.DiscountAmount)
}

func TestRuleEngine_CustomConfig_ModifiedParameters(t *testing.T) {
	cfg := &config.DiscountConfig{
		Rules: []config.RuleConfig{
			{
				Name:     "VolumeDiscountRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"ranges": []any{
						map[string]any{
							"customerType": "new",
							"minBooks":     2,
							"maxBooks":     5,
							"discountRate": 30,
						},
					},
				},
			},
		},
	}

	timeProvider := &mockTimeProvider{
		currentTime: getDateForWeekday(time.Monday),
	}
	engine, err := NewRuleEngine(cfg, timeProvider)
	require.NoError(t, err)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "2", Price: 1000},
		{ID: "3", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(3000), discount.CartAmount)
	assert.Equal(t, int64(900), discount.DiscountAmount)
	assert.Equal(t, 0.3, discount.DiscountPercent)
}

func TestRuleEngine_Friday_5PercentDiscount(t *testing.T) {
	engine := createTestEngine(time.Friday)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "2", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeOld, books)

	require.NoError(t, err)
	assert.Equal(t, int64(2000), discount.CartAmount)
	assert.Equal(t, int64(100), discount.DiscountAmount)
	assert.Equal(t, int64(1900), discount.TotalCost)
	assert.InDelta(t, 0.05, discount.DiscountPercent, 0.001)
}

func TestRuleEngine_Friday_BlocksVolumeDiscountRule(t *testing.T) {
	engine := createTestEngine(time.Friday)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "2", Price: 1000},
		{ID: "3", Price: 1000},
		{ID: "4", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(4000), discount.CartAmount)
	assert.Equal(t, int64(200), discount.DiscountAmount)
	assert.InDelta(t, 0.05, discount.DiscountPercent, 0.001)
}

func TestRuleEngine_Friday_DoesNotBlockBulkRule(t *testing.T) {
	engine := createTestEngine(time.Friday)

	books := []domain.Book{
		{ID: "hp", Price: 1000},
		{ID: "hp", Price: 1000},
		{ID: "hp", Price: 1000},
		{ID: "hp", Price: 1000},
		{ID: "hp", Price: 1000},
		{ID: "hp", Price: 1000},
		{ID: "lotr", Price: 1000},
		{ID: "dune", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(8000), discount.CartAmount)
	assert.Equal(t, int64(1300), discount.DiscountAmount)
	assert.InDelta(t, 0.1625, discount.DiscountPercent, 0.001)
}

func TestRuleEngine_NotFriday_VolumeDiscountApplies(t *testing.T) {
	engine := createTestEngine(time.Monday)

	books := []domain.Book{
		{ID: "1", Price: 1000},
		{ID: "2", Price: 1000},
		{ID: "3", Price: 1000},
		{ID: "4", Price: 1000},
	}

	discount, err := engine.Calculate(context.Background(), domain.CustomerTypeNew, books)

	require.NoError(t, err)
	assert.Equal(t, int64(4000), discount.CartAmount)
	assert.Equal(t, int64(800), discount.DiscountAmount)
	assert.Equal(t, 0.2, discount.DiscountPercent)
}

func createTestEngine(weekday time.Weekday) *RuleEngine {
	cfg := config.Default()
	timeProvider := &mockTimeProvider{
		currentTime: getDateForWeekday(weekday),
	}

	engine, err := NewRuleEngine(cfg, timeProvider)
	if err != nil {
		panic("failed to create test engine: " + err.Error())
	}

	return engine
}

func getDateForWeekday(weekday time.Weekday) time.Time {
	baseDate := time.Date(2025, 1, 13, 12, 0, 0, 0, time.UTC)
	currentWeekday := baseDate.Weekday()
	daysToAdd := int(weekday - currentWeekday)
	return baseDate.AddDate(0, 0, daysToAdd)
}
