package graphqlhub

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/adapters/graphqlhub/paymentcards"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

func (s searchHub) GetTransactionAndCards(
	ctx context.Context,
	transactionCardsInput domain.TransactionWalletInput,
) (*domain.TransactionWalletOutput, error) {

	request := paymentcards.BuildRequestTransactionAndCards(transactionCardsInput)

	var data paymentcards.TransactionWalletSearch

	err := s.service.SetDataByQueryRequest(ctx, request, &data)
	if err != nil {
		return nil, buildErrorBadGateway(err, msgErrFailGetTransactionAndCards)
	}

	return paymentcards.FillOutput(data), nil
}
