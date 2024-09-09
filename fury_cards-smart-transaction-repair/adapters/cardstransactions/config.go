package cardstransactions

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type Config struct {
	App domain.ConfigTransactionsApp
}

func NewConfig(configTransactionsApp domain.ConfigTransactionsApp) *Config {
	return &Config{
		App: configTransactionsApp,
	}
}
