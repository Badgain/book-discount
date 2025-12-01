package service

import (
	"context"
	"fmt"

	"github.com/Badgain/book-discount/internal/config"
	"github.com/Badgain/book-discount/internal/domain"
	"github.com/Badgain/book-discount/internal/domain/rules"
)

type RuleEngine struct {
	rules        []domain.DiscountRule
	timeProvider domain.TimeProvider
}

func NewRuleEngine(cfg *config.DiscountConfig, timeProvider domain.TimeProvider) (*RuleEngine, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if timeProvider == nil {
		return nil, fmt.Errorf("timeProvider cannot be nil")
	}

	var engineRules []domain.DiscountRule
	for _, ruleConfig := range cfg.Rules {
		if !ruleConfig.Enabled {
			continue
		}

		rule, err := rules.CreateRule(ruleConfig.Name, ruleConfig.Params)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule %s: %w", ruleConfig.Name, err)
		}

		engineRules = append(engineRules, rule)
	}

	return &RuleEngine{
		rules:        engineRules,
		timeProvider: timeProvider,
	}, nil
}

func (e *RuleEngine) Calculate(ctx context.Context, customerType domain.CustomerType, books []domain.Book) (domain.Discount, error) {
	if len(books) == 0 {
		return domain.Discount{
			CartAmount:      0,
			DiscountPercent: 0,
			TotalCost:       0,
			DiscountAmount:  0,
		}, nil
	}

	currentTime := e.timeProvider.Now()

	var totalDiscountAmount int64
	var cartAmount int64
	remainingBooks := books

	for _, book := range books {
		cartAmount += book.Price
	}

	for _, rule := range e.rules {
		if len(remainingBooks) == 0 {
			break
		}

		currentCtx := domain.RuleContext{
			CustomerType: customerType,
			Books:        remainingBooks,
			Time:         currentTime,
		}

		if !rule.CanApply(currentCtx) {
			continue
		}

		result, err := rule.Apply(currentCtx)
		if err != nil {
			return domain.Discount{}, fmt.Errorf("rule %s failed: %w", rule.Name(), err)
		}

		totalDiscountAmount += result.DiscountAmount
		remainingBooks = result.RemainingBooks

		if rule.BlocksOtherRules() {
			break
		}
	}

	finalAmount := cartAmount - totalDiscountAmount

	var discountPercent float64
	if cartAmount > 0 {
		discountPercent = float64(totalDiscountAmount) / float64(cartAmount)
	}

	return domain.Discount{
		CartAmount:      cartAmount,
		DiscountPercent: discountPercent,
		TotalCost:       finalAmount,
		DiscountAmount:  totalDiscountAmount,
	}, nil
}
