package domain

import (
	"net/url"
	"time"
)

type ConfigApplicationInit struct {
	Infra   ConfigInfra
	Options ConfigOptions
}

type ConfigInfra struct {
	URLTransactionsSearch string
	TransactionsApp       ConfigTransactionsApp
	KVSReparation         ConfigDatabaseKVS
	KVSBlockedlist        ConfigDatabaseKVS
	DSReparation          ConfigDatabaseDS
	DSBlockedUser         ConfigDatabaseDS
	DSValidation          ConfigDatabaseDS
	GraphqlKyc            ConfigGraphqlKyc
	Topics                ConfigTopics
	CardsDataModels       ConfigCardsDataModels
}

type ConfigOptions struct {
	TimeoutSecondsHTTP       time.Duration
	TimeoutHoursCacheDefault time.Duration
}

type ConfigDatabaseKVS struct {
	Name         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ReadRetries  int
	WriteRetries int
}

type ConfigDatabaseDS struct {
	Name                string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	Retries             int
	RetryDelay          time.Duration
	RetryAllowedMethods []string
}

type ConfigGraphqlKyc struct {
	URL          string
	PegasusToken string
	APISandbox   bool
}

type ConfigTopic struct {
	Name    string
	Timeout time.Duration
	Retries int
}

type ConfigTopics struct {
	ValidationResult ConfigTopic
}

type ConfigCardsDataModels struct {
	BaseURL url.URL
	Timeout time.Duration
}

type ConfigTransactionsApp struct {
	BetaScope              string
	BaseURL                url.URL
	BaseURLLegacy          url.URL
	ProvidersWithLegacyURL map[string]struct{}
}
