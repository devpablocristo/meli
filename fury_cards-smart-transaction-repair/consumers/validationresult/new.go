package validationresult

import (
	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

const consumerName = "consumer_validation_result"

type ValidationResultConsumer struct {
	service ports.ValidationService
	log     log.LogService
	errAPI  errAPI.ErrAPI
}

func NewConsumer(validationService ports.ValidationService, log log.LogService) *ValidationResultConsumer {
	return &ValidationResultConsumer{
		service: validationService,
		log:     log,
		errAPI: errAPI.ErrAPI{
			Log: log,
		},
	}
}
