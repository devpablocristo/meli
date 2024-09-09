package ds

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

type mock struct {
}

func (m mock) Save(ctx context.Context, paymentID string, input storage.ValidationResult) error {
	return nil
}

func NewMock() ports.ValidationResultRepository {
	return &mock{}
}
