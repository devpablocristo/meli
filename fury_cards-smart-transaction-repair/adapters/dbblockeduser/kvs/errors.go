package kvs

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

const (
	msgErrFail = "blockeduser database failure"
)

func buildErrorNotFound(id string) error {
	return domain.BuildErrorNotFound(id, "user not found", msgErrFail)
}

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrFail)
}
