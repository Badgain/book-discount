package service

import (
	"context"
	"fmt"

	"github.com/Badgain/book-discount/internal/config"
	"github.com/Badgain/book-discount/internal/domain"
)

type DiscountService struct {
	engine *RuleEngine
}

func NewDiscountService(timeProvider domain.TimeProvider) (*DiscountService, error) {
	return NewDiscountServiceWithConfig(timeProvider, "")
}

func NewDiscountServiceWithConfig(timeProvider domain.TimeProvider, configPath string) (*DiscountService, error) {
	if timeProvider == nil {
		timeProvider = &domain.RealTimeProvider{}
	}

	var cfg *config.DiscountConfig
	var err error

	if configPath != "" {
		cfg, err = config.Load(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		cfg = config.Default()
	}

	engine, err := NewRuleEngine(cfg, timeProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule engine: %w", err)
	}

	return &DiscountService{
		engine: engine,
	}, nil
}

func (s *DiscountService) Calculate(ctx context.Context, customerType domain.CustomerType, books []domain.Book) (domain.Discount, error) {
	return s.engine.Calculate(ctx, customerType, books)
}
