package score

import (
	"net/url"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type Config struct {
	BaseURL url.URL
}

func NewConfig(configScoreApp domain.ConfigCardsDataModels) *Config {
	return &Config{
		BaseURL: configScoreApp.BaseURL,
	}
}
