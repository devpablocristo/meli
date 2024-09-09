package cardstransactions

import (
	"context"

	_ "embed"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
)

type mock struct{}

func (m mock) Reverse(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
	return &domain.ReversalOutput{
		ReverseID: "reverso_smart_auth_mlb_test_aotorres_FyLRqgiKVMw9Vb8y737ppvNBsYpvkf",
	}, nil
}

var _ ports.CardsTransactions = (*mock)(nil)

func NewMock() ports.CardsTransactions {
	return mock{}
}
