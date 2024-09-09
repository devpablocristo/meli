package cardstransactions

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	ioutilabs "github.com/melisource/fury_cards-go-toolkit/pkg/ioutil/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

type cardsTransactions struct {
	config     *Config
	httpClient rest.HTTPClient
	ioutil     ioutilabs.IOutil
	json       json.JSON
	log        log.LogService
}

var _ ports.CardsTransactions = (*cardsTransactions)(nil)

func New(
	config *Config,
	httpClient rest.HTTPClient,
	ioutil ioutilabs.IOutil,
	json json.JSON,
	log log.LogService) ports.CardsTransactions {

	return &cardsTransactions{
		config:     config,
		httpClient: httpClient,
		ioutil:     ioutil,
		json:       json,
		log:        log,
	}
}
