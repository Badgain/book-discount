package service

import (
	"time"

	"github.com/Badgain/book-discount/internal/domain"
)

type mockTimeProvider struct {
	currentTime time.Time
}

func (m *mockTimeProvider) Now() time.Time {
	return m.currentTime
}

var _ domain.TimeProvider = (*mockTimeProvider)(nil)
