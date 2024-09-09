package graphqlkyc

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchkycclient/v1"
)

type searchKyc struct {
	service searchkycclient.Service
}

var _ ports.SearchKyc = (*searchKyc)(nil)

func New(service searchkycclient.Service) ports.SearchKyc {
	return searchKyc{
		service: service,
	}
}
