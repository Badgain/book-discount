package domain

import (
	"context"
	"time"
)

type CustomerType string

const (
	CustomerTypeNew CustomerType = "new"
	CustomerTypeOld CustomerType = "old"
)

type Book struct {
	ID    string
	Price int64
}

type Discount struct {
	CartAmount      int64
	DiscountPercent float64
	TotalCost       int64
	DiscountAmount  int64
}

type DiscountCalculator interface {
	Calculate(ctx context.Context, customerType CustomerType, books []Book) (Discount, error)
}

type RuleContext struct {
	CustomerType CustomerType
	Books        []Book
	Time         time.Time
}

type RuleResult struct {
	AppliedBooks   []Book
	RemainingBooks []Book
	DiscountAmount int64
	RuleName       string
}

type DiscountRule interface {
	Name() string
	CanApply(ruleCtx RuleContext) bool
	Apply(ruleCtx RuleContext) (RuleResult, error)
	BlocksOtherRules() bool
}

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (r *RealTimeProvider) Now() time.Time {
	return time.Now()
}

type MockTimeProvider struct {
	CurrentTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.CurrentTime
}
