package topicvalidation

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

const (
	msgErrFail = "validation_result_event_publication_failure"
)

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrFail)
}
