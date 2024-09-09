package topicvalidation

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

type mock struct{}

func (m mock) Publish(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
	return nil
}

var _ ports.EventValidationResult = (*mock)(nil)

func NewMock() ports.EventValidationResult {
	return mock{}
}
