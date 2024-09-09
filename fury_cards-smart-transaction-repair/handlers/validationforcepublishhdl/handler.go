package validationforcepublishhdl

import (
	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

type ValidationForcePublishHandler struct {
	service ports.ValidationService
	log     log.LogService
	errAPI  errAPI.ErrAPI
}

func NewHandler(service ports.ValidationService, log log.LogService) *ValidationForcePublishHandler {
	return &ValidationForcePublishHandler{
		service: service,
		log:     log,
		errAPI: errAPI.ErrAPI{
			Log: log,
		},
	}
}
