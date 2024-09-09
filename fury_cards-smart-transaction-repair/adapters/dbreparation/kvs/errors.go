package kvs

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	msgErrFail = "transaction repair database failure"
)

func buildErrorNotFound(id string) error {
	return domain.BuildErrorNotFound(id, "transaction repair not found", msgErrFail)
}

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrFail)
}

func buildErrorBadRequest(err error) error {
	return domain.BuildErrorBadRequest(err, msgErrFail)
}
