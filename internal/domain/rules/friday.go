package rules

import (
	"time"

	"github.com/Badgain/book-discount/internal/domain"
)

const (
	DefaultFridayDiscountRate = 5
)

type FridayRule struct {
	discountRate int
}

func NewFridayRule() *FridayRule {
	return &FridayRule{
		discountRate: DefaultFridayDiscountRate,
	}
}

func NewFridayRuleWithParams(discountRate int) *FridayRule {
	return &FridayRule{
		discountRate: discountRate,
	}
}

func (r *FridayRule) Name() string {
	return "FridayRule"
}

func (r *FridayRule) CanApply(ruleCtx domain.RuleContext) bool {
	return ruleCtx.Time.Weekday() == time.Friday
}

func (r *FridayRule) Apply(ruleCtx domain.RuleContext) (domain.RuleResult, error) {
	if err := validateBooks(ruleCtx.Books); err != nil {
		return domain.RuleResult{}, err
	}

	var discountAmount int64
	for _, book := range ruleCtx.Books {
		discountAmount += book.Price * int64(r.discountRate) / 100
	}

	return domain.RuleResult{
		AppliedBooks:   ruleCtx.Books,
		RemainingBooks: []domain.Book{},
		DiscountAmount: discountAmount,
		RuleName:       r.Name(),
	}, nil
}

func (r *FridayRule) BlocksOtherRules() bool {
	return true
}
