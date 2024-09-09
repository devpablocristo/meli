package reversehdl

import (
	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

type ReverseHandler struct {
	service ports.ReparationService
	log     log.LogService
	errAPI  errAPI.ErrAPI
}

func NewHandler(service ports.ReparationService, log log.LogService) *ReverseHandler {
	return &ReverseHandler{
		service: service,
		log:     log,
		errAPI: errAPI.ErrAPI{
			Log: log,
		},
	}
}
