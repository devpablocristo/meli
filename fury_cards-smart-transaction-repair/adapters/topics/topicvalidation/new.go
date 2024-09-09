package topicvalidation

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/publisher/v1"
)

type validationEvent struct {
	publisher publisher.Publisher
}

var _ ports.EventValidationResult = (*validationEvent)(nil)

func New(publisher publisher.Publisher) ports.EventValidationResult {
	return &validationEvent{
		publisher,
	}
}
