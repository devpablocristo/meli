package graphqlhub

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/adapters/graphqlhub/transactionhistory"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

func (s searchHub) GetTransactionHistory(
	ctx context.Context,
	transactionHistoryInput domain.TransactionHistoryInput,
) (*domain.TransactionHistoryOutput, error) {

	request := transactionhistory.BuildRequestGetTransactionHistory(transactionHistoryInput)

	var data transactionhistory.History
	err := s.service.SetDataByQueryRequest(ctx, request, &data)
	if err != nil {
		return nil, buildErrorBadGateway(err, msgErrFailGetTransactionHistory)
	}

	return transactionhistory.FillOutput(data), nil
}
