package ds

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

const (
	msgErrorSave = "error saving validation result"
)

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrorSave)
}
