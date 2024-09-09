package policy

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
)

type mock struct{}

func (m mock) EvaluateWithUser(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
	return &domain.PolicyOutput{
		IsAuthorized:    true,
		RestrictsFailed: nil,
	}, nil
}

func NewMock() ports.Policy {
	return mock{}
}
