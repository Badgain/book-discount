package config

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type RuleConfig struct {
	Name     string         `yaml:"name"`
	Enabled  bool           `yaml:"enabled"`
	Priority int            `yaml:"priority"`
	Params   map[string]any `yaml:"params"`
}

type DiscountConfig struct {
	Rules []RuleConfig `yaml:"rules"`
}

func Load(configPath string) (*DiscountConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config DiscountConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config.Sort()

	return &config, nil
}

func Default() *DiscountConfig {
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
				Priority: 10,
				Params: map[string]any{
					"discountRate": 5,
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
						map[string]any{
							"customerType": "old",
							"minBooks":     2,
							"maxBooks":     5,
							"discountRate": 10,
						},
						map[string]any{
							"customerType": "old",
							"minBooks":     6,
							"maxBooks":     10,
							"discountRate": 5,
						},
						map[string]any{
							"customerType": "old",
							"minBooks":     11,
							"maxBooks":     nil,
							"discountRate": 2,
						},
					},
				},
			},
		},
	}

	cfg.Sort()
	return cfg
}

func (c *DiscountConfig) Validate() error {
	if len(c.Rules) == 0 {
		return fmt.Errorf("config must have at least one rule")
	}

	priorityMap := make(map[int]string)
	nameMap := make(map[string]bool)

	for i, rule := range c.Rules {
		if rule.Name == "" {
			return fmt.Errorf("rule %d: name cannot be empty", i)
		}

		if nameMap[rule.Name] {
			return fmt.Errorf("duplicate rule name: %s", rule.Name)
		}
		nameMap[rule.Name] = true

		if existingRule, exists := priorityMap[rule.Priority]; exists {
			return fmt.Errorf("duplicate priority %d: rules %s and %s", rule.Priority, existingRule, rule.Name)
		}
		priorityMap[rule.Priority] = rule.Name

		if err := validateRuleParams(rule); err != nil {
			return fmt.Errorf("rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

func validateRuleParams(rule RuleConfig) error {
	switch rule.Name {
	case "BulkSameBookRule":
		return validateBulkSameBookParams(rule.Params)
	case "FridayRule":
		return validateFridayParams(rule.Params)
	case "VolumeDiscountRule":
		return validateVolumeDiscountParams(rule.Params)
	default:
		return nil
	}
}

func validateBulkSameBookParams(params map[string]any) error {
	minBooks, ok := params["minBooks"].(int)
	if !ok {
		return fmt.Errorf("minBooks is required and must be int")
	}
	if minBooks < 0 {
		return fmt.Errorf("minBooks must be >= 0, got %d", minBooks)
	}

	discountRate, ok := params["discountRate"].(int)
	if !ok {
		return fmt.Errorf("discountRate is required and must be int")
	}
	if discountRate < 0 || discountRate > 100 {
		return fmt.Errorf("discountRate must be between 0 and 100, got %d", discountRate)
	}

	return nil
}

func validateFridayParams(params map[string]any) error {
	discountRate, ok := params["discountRate"].(int)
	if !ok {
		return fmt.Errorf("discountRate is required and must be int")
	}
	if discountRate < 0 || discountRate > 100 {
		return fmt.Errorf("discountRate must be between 0 and 100, got %d", discountRate)
	}

	return nil
}

func validateVolumeDiscountParams(params map[string]any) error {
	rangesRaw, ok := params["ranges"]
	if !ok {
		return fmt.Errorf("ranges is required")
	}

	rangesList, ok := rangesRaw.([]any)
	if !ok {
		return fmt.Errorf("ranges must be a list")
	}

	if len(rangesList) == 0 {
		return fmt.Errorf("ranges cannot be empty")
	}

	for i, rngRaw := range rangesList {
		rngMap, ok := rngRaw.(map[string]any)
		if !ok {
			return fmt.Errorf("range %d must be a map", i)
		}

		customerType, ok := rngMap["customerType"].(string)
		if !ok {
			return fmt.Errorf("range %d: customerType is required and must be string", i)
		}
		if customerType != "new" && customerType != "old" {
			return fmt.Errorf("range %d: customerType must be 'new' or 'old', got '%s'", i, customerType)
		}

		minBooks, ok := rngMap["minBooks"].(int)
		if !ok {
			return fmt.Errorf("range %d: minBooks is required and must be int", i)
		}
		if minBooks < 0 {
			return fmt.Errorf("range %d: minBooks must be >= 0, got %d", i, minBooks)
		}

		if maxBooksRaw, exists := rngMap["maxBooks"]; exists && maxBooksRaw != nil {
			maxBooks, ok := maxBooksRaw.(int)
			if !ok {
				return fmt.Errorf("range %d: maxBooks must be int or null", i)
			}
			if maxBooks < minBooks {
				return fmt.Errorf("range %d: maxBooks (%d) must be >= minBooks (%d)", i, maxBooks, minBooks)
			}
		}

		discountRate, ok := rngMap["discountRate"].(int)
		if !ok {
			return fmt.Errorf("range %d: discountRate is required and must be int", i)
		}
		if discountRate < 0 || discountRate > 100 {
			return fmt.Errorf("range %d: discountRate must be between 0 and 100, got %d", i, discountRate)
		}
	}

	return nil
}

func (c *DiscountConfig) Sort() {
	sort.Slice(c.Rules, func(i, j int) bool {
		return c.Rules[i].Priority < c.Rules[j].Priority
	})
}
