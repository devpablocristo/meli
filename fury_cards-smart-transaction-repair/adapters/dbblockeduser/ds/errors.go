package ds

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

const (
	msgErrFail = "ds_blockeduser database failure"
)

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrFail)
}
