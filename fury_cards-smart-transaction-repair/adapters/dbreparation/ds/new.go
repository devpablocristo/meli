package ds

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

type dsTransactionRepair struct {
	repo  dsclient.Service
	query dsclient.Query
}

var _ ports.ReparationSearchRepository = (*dsTransactionRepair)(nil)

func New(dsClient dsclient.Service, dsQuery dsclient.Query) ports.ReparationSearchRepository {
	return &dsTransactionRepair{
		repo:  dsClient,
		query: dsQuery,
	}
}
