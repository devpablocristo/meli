package payment

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
)

func BuildRequestPaymentTransaction(transactionInput domain.TransactionInput) graphql.Request {
	return graphql.Request{
		Query: getPaymentTransaction,
		Params: map[string]interface{}{
			"id": transactionInput.PaymentID,
		},
	}
}
