package paymentcards

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
)

func BuildRequestTransactionAndCards(transactionCardsInput domain.TransactionWalletInput) graphql.Request {
	return graphql.Request{
		Query: queryTransactionAndCards,
		Params: map[string]interface{}{
			"user_id":    transactionCardsInput.UserID,
			"payment_id": transactionCardsInput.PaymentID,
		},
	}
}
