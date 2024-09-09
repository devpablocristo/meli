package graphqlhub

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchhubclient/v1"
)

type searchHub struct {
	service searchhubclient.Service
}

var _ ports.SearchHub = (*searchHub)(nil)

func New(service searchhubclient.Service) ports.SearchHub {
	return searchHub{
		service: service,
	}
}
