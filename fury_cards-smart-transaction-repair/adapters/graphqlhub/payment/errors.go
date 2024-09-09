package payment

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	msgErrFailGetPaymentTransaction = "searchhub.get_payment_transaction.failure"
)

func buildErrorBadGateway(err error, msgErr string) error {
	return domain.BuildErrorBadGateway(err, msgErr)
}
