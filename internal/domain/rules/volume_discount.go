package rules

import (
	"fmt"

	"github.com/Badgain/book-discount/internal/domain"
)

type VolumeRange struct {
	CustomerType domain.CustomerType
	MinBooks     int
	MaxBooks     *int
	DiscountRate int
}

type VolumeDiscountRule struct {
	ranges []VolumeRange
}

func NewVolumeDiscountRule(ranges []VolumeRange) *VolumeDiscountRule {
	return &VolumeDiscountRule{
		ranges: ranges,
	}
}

func (r *VolumeDiscountRule) Name() string {
	return "VolumeDiscountRule"
}

func (r *VolumeDiscountRule) CanApply(ruleCtx domain.RuleContext) bool {
	bookCount := len(ruleCtx.Books)

	for _, rng := range r.ranges {
		if rng.CustomerType != ruleCtx.CustomerType {
			continue
		}

		if bookCount < rng.MinBooks {
			continue
		}

		if rng.MaxBooks != nil && bookCount > *rng.MaxBooks {
			continue
		}

		return true
	}

	return false
}

func (r *VolumeDiscountRule) Apply(ruleCtx domain.RuleContext) (domain.RuleResult, error) {
	if err := validateBooks(ruleCtx.Books); err != nil {
		return domain.RuleResult{}, err
	}

	bookCount := len(ruleCtx.Books)
	discountRate := r.getDiscountRate(ruleCtx.CustomerType, bookCount)

	if discountRate == 0 {
		return domain.RuleResult{}, fmt.Errorf("no matching discount range for customer type %s and %d books", ruleCtx.CustomerType, bookCount)
	}

	var discountAmount int64
	for _, book := range ruleCtx.Books {
		discountAmount += book.Price * int64(discountRate) / 100
	}

	return domain.RuleResult{
		AppliedBooks:   ruleCtx.Books,
		RemainingBooks: []domain.Book{},
		DiscountAmount: discountAmount,
		RuleName:       r.Name(),
	}, nil
}

func (r *VolumeDiscountRule) BlocksOtherRules() bool {
	return false
}

func (r *VolumeDiscountRule) getDiscountRate(customerType domain.CustomerType, bookCount int) int {
	for _, rng := range r.ranges {
		if rng.CustomerType != customerType {
			continue
		}

		if bookCount < rng.MinBooks {
			continue
		}

		if rng.MaxBooks != nil && bookCount > *rng.MaxBooks {
			continue
		}

		return rng.DiscountRate
	}

	return 0
}
