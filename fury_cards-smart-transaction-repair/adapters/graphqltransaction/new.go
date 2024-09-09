package graphqltransaction

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	ioutilabs "github.com/melisource/fury_cards-go-toolkit/pkg/ioutil/v1"
	json "github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
)

type transactionSearch struct {
	config     *Config
	httpClient rest.HTTPClient
	ioutil     ioutilabs.IOutil
	json       json.JSON
}

var _ ports.SearchTransaction = (*transactionSearch)(nil)

func New(
	config *Config,
	httpClient rest.HTTPClient,
	ioutil ioutilabs.IOutil,
	json json.JSON) ports.SearchTransaction {

	return &transactionSearch{
		config:     config,
		httpClient: httpClient,
		ioutil:     ioutil,
		json:       json,
	}
}
