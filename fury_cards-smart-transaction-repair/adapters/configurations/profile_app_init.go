package configurations

import (
	"time"
)

type applicationInit struct {
	Infra   infra   `json:"infra"`
	Options options `json:"options"`
}

type infra struct {
	URLTransactionsSearch string          `json:"url_transactions_search"`
	TransactionsApp       transactionsApp `json:"transactions_app"`
	KVSReparation         databaseKVS     `json:"kvs_reparation"`
	KVSBlockedlist        databaseKVS     `json:"kvs_blockedlist"`
	DSReparation          databaseDS      `json:"ds_reparation"`
	DSValidation          databaseDS      `json:"ds_validation"`
	DSBlockedUser         databaseDS      `json:"ds_blockeduser"`
	GraphqlKyc            graphqlKyc      `json:"graphql_kyc"`
	Topics                topics          `json:"topics"`
	CardsDataModels       cardsDataModels `json:"cards_data_models"`
}

type options struct {
	TimeoutSecondsHTTP       time.Duration `json:"timeout_seconds_http"`
	TimeoutHoursCacheDefault time.Duration `json:"timeout_hours_cache_default"`
}

type databaseKVS struct {
	Name         string        `json:"name"`
	ReadTimeout  time.Duration `json:"read_timeout_milliseconds"`
	WriteTimeout time.Duration `json:"write_timeout_milliseconds"`
	ReadRetries  int           `json:"read_retries"`
	WriteRetries int           `json:"write_retries"`
}

type databaseDS struct {
	Name                string        `json:"name"`
	ReadTimeout         time.Duration `json:"read_timeout_milliseconds"`
	WriteTimeout        time.Duration `json:"write_timeout_milliseconds"`
	Retries             int           `json:"retries"`
	RetryDelay          time.Duration `json:"retry_delay_milliseconds"`
	RetryAllowedMethods []string      `json:"retry_allowed_methods"`
}

type graphqlKyc struct {
	URL          string `json:"url"`
	PegasusToken string `json:"x-pegasus-token"`
	APISandbox   bool   `json:"x-api-sandbox"`
}

type topic struct {
	Name    string        `json:"name"`
	Timeout time.Duration `json:"timeout_milliseconds"`
	Retries int           `json:"retries"`
}

type topics struct {
	ValidationResult topic `json:"validation_result"`
}

type cardsDataModels struct {
	BaseURL string        `json:"base_url"`
	Timeout time.Duration `json:"timeout_milliseconds"`
}

type transactionsApp struct {
	BetaScope              string              `json:"beta_scope"`
	BaseURL                string              `json:"base_url"`
	BaseURLLegacy          string              `json:"base_url_legacy"`
	ProvidersWithLegacyURL map[string]struct{} `json:"providers_with_legacy_url"`
}
