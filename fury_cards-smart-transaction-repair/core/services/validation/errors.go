package validation

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

// The error log is not needed at this point as it will be done in the reparation service

const (
	msgReasonErr                  = "customer not eligible for reversal"
	msgErr                        = "validation failed"
	msgErrParametersNotConfigured = "rules not found"
)

func buildErrorNotEligible(reason domain.Reason, creationDatetime time.Time) error {
	causeErr := &domain.ValidationCauseErrResponse{
		Reason:           msgReasonErr,
		CreationDatetime: creationDatetime,
		ReasonDetail:     reason,
	}

	return domain.BuildErrorNotEligible(causeErr, "validation result")
}

func buildErrorUnprocessableEntity(causeErr error) error {
	return domain.BuildErrorUnprocessableEntity(causeErr, msgErr)
}

func buildErrorBadGateway(causeErr error, msgErr string) error {
	return domain.BuildErrorBadGateway(causeErr, msgErr)
}
