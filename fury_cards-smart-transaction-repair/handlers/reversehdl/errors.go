package reversehdl

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	msgError          = "invalid request"
	msgHeaderRequired = "%s header is required"
)

func (r ReverseHandler) buildErrorBadRequest(err error) error {
	return domain.BuildErrorBadRequest(err, msgError)
}
