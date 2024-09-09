package transactions

import (
	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

type transactionType string

type TransactionsConsumer struct {
	service ports.BlockedUserService
	log     log.LogService
	errAPI  errAPI.ErrAPI
}

func NewConsumer(blockedUserService ports.BlockedUserService, log log.LogService) *TransactionsConsumer {
	return &TransactionsConsumer{
		service: blockedUserService,
		log:     log,
		errAPI: errAPI.ErrAPI{
			Log: log,
		},
	}
}
