package blockeduser

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
)

type blocked struct {
	reparationSearchRepo  ports.ReparationSearchRepository
	blockedUserRepo       ports.BlockedUserRepository
	blockedUserSearchRepo ports.BlockedUserSearchRepository
}

var _ ports.BlockedUserService = (*blocked)(nil)

func New(
	reparationSearchRepo ports.ReparationSearchRepository,
	blockedUserRepo ports.BlockedUserRepository,
	blockedUserSearchRepo ports.BlockedUserSearchRepository,
) ports.BlockedUserService {

	return blocked{
		reparationSearchRepo:  reparationSearchRepo,
		blockedUserRepo:       blockedUserRepo,
		blockedUserSearchRepo: blockedUserSearchRepo,
	}
}
