package graphqlhub

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/adapters/graphqlhub/payment"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

func (s searchHub) GetPaymentTransaction(
	ctx context.Context,
	transactionInput domain.TransactionInput,
	wrapFields *field.WrappedFields,
) (domain.TransactionOutput, error) {
	return payment.GetPaymentTransaction(ctx, s.service, transactionInput, wrapFields)
}
