package rules

import (
	"fmt"

	"github.com/Badgain/book-discount/internal/domain"
)

type RuleFactory func(params map[string]any) (domain.DiscountRule, error)

var ruleFactories = make(map[string]RuleFactory)

func RegisterRule(name string, factory RuleFactory) {
	ruleFactories[name] = factory
}

func CreateRule(name string, params map[string]any) (domain.DiscountRule, error) {
	factory, ok := ruleFactories[name]
	if !ok {
		return nil, fmt.Errorf("unknown rule: %s", name)
	}
	return factory(params)
}

func createBulkSameBookRule(params map[string]any) (domain.DiscountRule, error) {
	minBooks := params["minBooks"].(int)
	discountRate := params["discountRate"].(int)
	return NewBulkSameBookRuleWithParams(minBooks, discountRate), nil
}

func createFridayRule(params map[string]any) (domain.DiscountRule, error) {
	discountRate := params["discountRate"].(int)
	return NewFridayRuleWithParams(discountRate), nil
}

func createVolumeDiscountRule(params map[string]any) (domain.DiscountRule, error) {
	rangesList := params["ranges"].([]any)

	ranges := make([]VolumeRange, 0, len(rangesList))
	for _, rngRaw := range rangesList {
		rngMap := rngRaw.(map[string]any)

		customerTypeStr := rngMap["customerType"].(string)
		minBooks := rngMap["minBooks"].(int)

		var maxBooks *int
		if maxBooksRaw, exists := rngMap["maxBooks"]; exists && maxBooksRaw != nil {
			maxBooksVal := maxBooksRaw.(int)
			maxBooks = &maxBooksVal
		}

		discountRate := rngMap["discountRate"].(int)

		ranges = append(ranges, VolumeRange{
			CustomerType: domain.CustomerType(customerTypeStr),
			MinBooks:     minBooks,
			MaxBooks:     maxBooks,
			DiscountRate: discountRate,
		})
	}

	return NewVolumeDiscountRule(ranges), nil
}

func init() {
	RegisterRule("BulkSameBookRule", createBulkSameBookRule)
	RegisterRule("FridayRule", createFridayRule)
	RegisterRule("VolumeDiscountRule", createVolumeDiscountRule)
}
