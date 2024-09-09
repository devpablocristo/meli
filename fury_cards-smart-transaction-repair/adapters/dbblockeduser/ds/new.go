package ds

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

type dsBlockedUser struct {
	repo  dsclient.Service
	query dsclient.Query
}

var _ ports.BlockedUserSearchRepository = (*dsBlockedUser)(nil)

func New(dsClient dsclient.Service, dsQuery dsclient.Query) ports.BlockedUserSearchRepository {
	return &dsBlockedUser{
		repo:  dsClient,
		query: dsQuery,
	}
}
