package rules

import (
	"testing"

	"github.com/Badgain/book-discount/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRule_CustomRule(t *testing.T) {
	customRuleCalled := false

	RegisterRule("TestCustomRule", func(params map[string]any) (domain.DiscountRule, error) {
		customRuleCalled = true
		discountRate := params["discountRate"].(int)
		return NewFridayRuleWithParams(discountRate), nil
	})

	params := map[string]any{
		"discountRate": 15,
	}

	rule, err := CreateRule("TestCustomRule", params)

	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.True(t, customRuleCalled)
	assert.Equal(t, "FridayRule", rule.Name())
}

func TestCreateRule_UnknownRule(t *testing.T) {
	params := map[string]any{}

	rule, err := CreateRule("NonExistentRule", params)

	assert.Error(t, err)
	assert.Nil(t, rule)
	assert.Contains(t, err.Error(), "unknown rule")
}

func TestCreateRule_BulkSameBookRule(t *testing.T) {
	params := map[string]any{
		"minBooks":     3,
		"discountRate": 25,
	}

	rule, err := CreateRule("BulkSameBookRule", params)

	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "BulkSameBookRule", rule.Name())
}

func TestCreateRule_FridayRule(t *testing.T) {
	params := map[string]any{
		"discountRate": 10,
	}

	rule, err := CreateRule("FridayRule", params)

	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "FridayRule", rule.Name())
}

func TestCreateRule_VolumeDiscountRule(t *testing.T) {
	params := map[string]any{
		"ranges": []any{
			map[string]any{
				"customerType": "new",
				"minBooks":     2,
				"maxBooks":     5,
				"discountRate": 15,
			},
		},
	}

	rule, err := CreateRule("VolumeDiscountRule", params)

	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "VolumeDiscountRule", rule.Name())
}
