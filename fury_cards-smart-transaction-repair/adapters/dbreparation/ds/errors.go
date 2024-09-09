package ds

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

const (
	msgErrFail = "transaction repair document search failure"
)

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrFail)
}

func buildErrorNotFound(paymentID string, err error) error {
	return domain.BuildErrorNotFound(paymentID, err.Error(), msgErrFail)
}
