package score

import (
	"context"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
)

type mock struct{}

func (m mock) GetScore(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
	return &domain.ScoreOutput{
		Result: "1",
		Score:  1.987745545,
	}, nil
}

var _ ports.Score = (*mock)(nil)

func NewMock() ports.Score {
	return mock{}
}
