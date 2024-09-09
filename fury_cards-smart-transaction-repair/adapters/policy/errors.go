package policy

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

const (
	msgErrFail = "transaction repair policy validated failure"
)

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrFail)
}
