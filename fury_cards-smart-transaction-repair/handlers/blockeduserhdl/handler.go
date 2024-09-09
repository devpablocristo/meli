package blockeduserhdl

import (
	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

type BlockedUserHandler struct {
	service ports.BlockedUserService
	log     log.LogService
	errAPI  errAPI.ErrAPI
}

func NewHandler(service ports.BlockedUserService, log log.LogService) *BlockedUserHandler {
	return &BlockedUserHandler{
		service: service,
		log:     log,
		errAPI: errAPI.ErrAPI{
			Log: log,
		},
	}
}
