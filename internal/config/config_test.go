package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscountConfig_Validate_EmptyConfig(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one rule")
}

func TestDiscountConfig_Validate_DuplicatePriority(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{
			{
				Name:     "BulkSameBookRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"minBooks":     5,
					"discountRate": 40,
				},
			},
			{
				Name:     "FridayRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"discountRate": 5,
				},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate priority")
}

func TestDiscountConfig_Validate_DuplicateName(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{
			{
				Name:     "BulkSameBookRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"minBooks":     5,
					"discountRate": 40,
				},
			},
			{
				Name:     "BulkSameBookRule",
				Enabled:  true,
				Priority: 2,
				Params: map[string]any{
					"minBooks":     3,
					"discountRate": 30,
				},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate rule name")
}

func TestDiscountConfig_Validate_InvalidDiscountRate(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{
			{
				Name:     "BulkSameBookRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"minBooks":     5,
					"discountRate": 150,
				},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "between 0 and 100")
}

func TestDiscountConfig_Validate_NegativeMinBooks(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{
			{
				Name:     "BulkSameBookRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"minBooks":     -5,
					"discountRate": 40,
				},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be >= 0")
}

func TestDiscountConfig_Validate_InvalidCustomerType(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{
			{
				Name:     "VolumeDiscountRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"ranges": []any{
						map[string]any{
							"customerType": "premium",
							"minBooks":     2,
							"maxBooks":     5,
							"discountRate": 20,
						},
					},
				},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be 'new' or 'old'")
}

func TestDiscountConfig_Validate_MaxBooksLessThanMinBooks(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{
			{
				Name:     "VolumeDiscountRule",
				Enabled:  true,
				Priority: 1,
				Params: map[string]any{
					"ranges": []any{
						map[string]any{
							"customerType": "new",
							"minBooks":     10,
							"maxBooks":     5,
							"discountRate": 20,
						},
					},
				},
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maxBooks")
	assert.Contains(t, err.Error(), "must be >= minBooks")
}

func TestDiscountConfig_Validate_ValidConfig(t *testing.T) {
	cfg := Default()
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestDiscountConfig_Sort(t *testing.T) {
	cfg := &DiscountConfig{
		Rules: []RuleConfig{
			{
				Name:     "Rule3",
				Priority: 30,
			},
			{
				Name:     "Rule1",
				Priority: 10,
			},
			{
				Name:     "Rule2",
				Priority: 20,
			},
		},
	}

	cfg.Sort()

	assert.Equal(t, "Rule1", cfg.Rules[0].Name)
	assert.Equal(t, "Rule2", cfg.Rules[1].Name)
	assert.Equal(t, "Rule3", cfg.Rules[2].Name)
}
