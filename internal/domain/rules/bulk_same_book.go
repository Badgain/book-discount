package rules

import (
	"errors"
	"sort"

	"github.com/Badgain/book-discount/internal/domain"
)

const (
	DefaultMinBooksForBulkDiscount = 5
	DefaultBulkDiscountRate        = 40
)

type BulkSameBookRule struct {
	minBooks     int
	discountRate int
}

func NewBulkSameBookRule() *BulkSameBookRule {
	return &BulkSameBookRule{
		minBooks:     DefaultMinBooksForBulkDiscount,
		discountRate: DefaultBulkDiscountRate,
	}
}

func NewBulkSameBookRuleWithParams(minBooks, discountRate int) *BulkSameBookRule {
	return &BulkSameBookRule{
		minBooks:     minBooks,
		discountRate: discountRate,
	}
}

func (r *BulkSameBookRule) Name() string {
	return "BulkSameBookRule"
}

func (r *BulkSameBookRule) CanApply(ruleCtx domain.RuleContext) bool {
	if len(ruleCtx.Books) == 0 {
		return false
	}

	bookGroups := groupBooksByID(ruleCtx.Books)
	for _, group := range bookGroups {
		if len(group) > r.minBooks {
			return true
		}
	}
	return false
}

func (r *BulkSameBookRule) Apply(ruleCtx domain.RuleContext) (domain.RuleResult, error) {
	if err := validateBooks(ruleCtx.Books); err != nil {
		return domain.RuleResult{}, err
	}

	bookGroups := groupBooksByID(ruleCtx.Books)

	appliedBooks := make([]domain.Book, 0, len(ruleCtx.Books))
	remainingBooks := make([]domain.Book, 0, len(ruleCtx.Books))
	var discountAmount int64

	bookIDs := getSortedKeys(bookGroups)

	for _, bookID := range bookIDs {
		group := bookGroups[bookID]
		if len(group) > r.minBooks {
			for i, book := range group {
				if i%2 == 1 {
					discountAmount += book.Price * int64(r.discountRate) / 100
				}
			}
			appliedBooks = append(appliedBooks, group...)
		} else {
			remainingBooks = append(remainingBooks, group...)
		}
	}

	return domain.RuleResult{
		AppliedBooks:   appliedBooks,
		RemainingBooks: remainingBooks,
		DiscountAmount: discountAmount,
		RuleName:       r.Name(),
	}, nil
}

func (r *BulkSameBookRule) BlocksOtherRules() bool {
	return false
}

func validateBooks(books []domain.Book) error {
	if books == nil {
		return errors.New("books cannot be nil")
	}

	for _, book := range books {
		if book.ID == "" {
			return errors.New("book ID cannot be empty")
		}
		if book.Price < 0 {
			return errors.New("book price cannot be negative")
		}
	}
	return nil
}

func groupBooksByID(books []domain.Book) map[string][]domain.Book {
	groups := make(map[string][]domain.Book)
	for _, book := range books {
		groups[book.ID] = append(groups[book.ID], book)
	}
	return groups
}

func getSortedKeys(groups map[string][]domain.Book) []string {
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
