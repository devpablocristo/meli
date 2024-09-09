package graphqlhub

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	msgErrFailCountChargeback        = "search-hub count_chargeback failure"
	msgErrFailGetTransactionAndCards = "search-hub get_transaction_and_cards failure"
	msgErrFailGetTransactionHistory  = "search-hub get_transaction_history failure"
)

func buildErrorBadGateway(err error, msgErr string) error {
	return domain.BuildErrorBadGateway(err, msgErr)
}
