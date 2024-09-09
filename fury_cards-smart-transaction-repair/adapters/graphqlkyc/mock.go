package graphqlkyc

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
)

type mock struct{}

func (m mock) GetUserKyc(
	ctx context.Context,
	input domain.SearchKycInput,
) (*domain.SearchKycOutput, error) {
	return &domain.SearchKycOutput{KycIdentificationID: "abc"}, nil
}

var _ ports.SearchKyc = (*mock)(nil)

func NewMock() ports.SearchKyc {
	return mock{}
}
