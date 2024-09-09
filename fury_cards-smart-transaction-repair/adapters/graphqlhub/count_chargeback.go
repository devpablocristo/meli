package graphqlhub

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

func (s searchHub) CountChargebackFromUser(
	ctx context.Context,
	input domain.SearchHubCountChargebackInput,
) (*domain.SearchHubCountChargebackOutput, error) {

	output, err := s.service.CountChargebacksInTheLastDays(ctx, input.UserID, input.LastDays)
	if err != nil {
		return nil, buildErrorBadGateway(err, msgErrFailCountChargeback)
	}

	return (*domain.SearchHubCountChargebackOutput)(output), nil
}
