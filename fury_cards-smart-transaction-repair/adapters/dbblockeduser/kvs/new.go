package kvs

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	kvs "github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
)

type blockedUserRepository struct {
	repo kvs.Service
}

var _ ports.BlockedUserRepository = (*blockedUserRepository)(nil)

func New(ksvClient kvs.Service) ports.BlockedUserRepository {
	return &blockedUserRepository{
		repo: ksvClient,
	}
}
