package configurations

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const keyAppInit = "app-init"

func (c *configService) LoadApplicationInit(ctx context.Context) (*domain.ConfigApplicationInit, error) {
	var applicationInit applicationInit

	err := c.serviceCache.GetJson(ctx, keyAppInit, &applicationInit)
	if err != nil {
		return nil, treatError(err, keyAppInit)
	}

	return fillConfigApplicationInit(&applicationInit), nil
}

func fillConfigApplicationInit(applicationInit *applicationInit) *domain.ConfigApplicationInit {
	configApplicationInit := &domain.ConfigApplicationInit{
		Infra: domain.ConfigInfra{
			URLTransactionsSearch: applicationInit.Infra.URLTransactionsSearch,
			TransactionsApp:       toConfigTransactionsApp(applicationInit.Infra.TransactionsApp),
			KVSReparation:         domain.ConfigDatabaseKVS(applicationInit.Infra.KVSReparation),
			KVSBlockedlist:        domain.ConfigDatabaseKVS(applicationInit.Infra.KVSBlockedlist),
			DSReparation:          domain.ConfigDatabaseDS(applicationInit.Infra.DSReparation),
			DSBlockedUser:         domain.ConfigDatabaseDS(applicationInit.Infra.DSBlockedUser),
			DSValidation:          domain.ConfigDatabaseDS(applicationInit.Infra.DSValidation),
			GraphqlKyc:            domain.ConfigGraphqlKyc(applicationInit.Infra.GraphqlKyc),
			Topics: domain.ConfigTopics{
				ValidationResult: domain.ConfigTopic(applicationInit.Infra.Topics.ValidationResult),
			},
			CardsDataModels: toConfigCardsDataModels(applicationInit.Infra.CardsDataModels),
		},
		Options: domain.ConfigOptions(applicationInit.Options),
	}

	setTimeout(&configApplicationInit.Infra)
	setTimeoutOptions(&configApplicationInit.Options)
	return configApplicationInit
}

func setTimeout(infra *domain.ConfigInfra) {
	infra.DSReparation.ReadTimeout *= time.Millisecond
	infra.DSReparation.WriteTimeout *= time.Millisecond

	infra.DSBlockedUser.ReadTimeout *= time.Millisecond
	infra.DSBlockedUser.WriteTimeout *= time.Millisecond

	infra.DSValidation.ReadTimeout *= time.Millisecond
	infra.DSValidation.WriteTimeout *= time.Millisecond

	infra.KVSBlockedlist.ReadTimeout *= time.Millisecond
	infra.KVSBlockedlist.WriteTimeout *= time.Millisecond

	infra.KVSReparation.ReadTimeout *= time.Millisecond
	infra.KVSReparation.WriteTimeout *= time.Millisecond

	infra.Topics.ValidationResult.Timeout *= time.Millisecond

	infra.CardsDataModels.Timeout *= time.Millisecond
}

func setTimeoutOptions(options *domain.ConfigOptions) {
	options.TimeoutSecondsHTTP *= time.Second
	options.TimeoutHoursCacheDefault *= time.Hour
}

func treatError(err error, keyProfile string) error {
	if strings.Contains(err.Error(), msgErrNotFoundProfile) {
		return buildErrorNotFound(keyProfile, err)
	}

	return buildErrorBadGateway(err)
}

func toConfigTransactionsApp(transactionsApp transactionsApp) domain.ConfigTransactionsApp {
	baseURLParsed, err := url.Parse(transactionsApp.BaseURL)
	if err != nil {
		panic(err)
	}

	baseURLLegacyParsed, err := url.Parse(transactionsApp.BaseURLLegacy)
	if err != nil {
		panic(err)
	}

	config := domain.ConfigTransactionsApp{
		BetaScope:              transactionsApp.BetaScope,
		BaseURL:                *baseURLParsed,
		BaseURLLegacy:          *baseURLLegacyParsed,
		ProvidersWithLegacyURL: transactionsApp.ProvidersWithLegacyURL,
	}

	return config
}

func toConfigCardsDataModels(cardsDataModels cardsDataModels) domain.ConfigCardsDataModels {
	baseURLParsed, err := url.Parse(cardsDataModels.BaseURL)
	if err != nil {
		panic(err)
	}

	config := domain.ConfigCardsDataModels{
		BaseURL: *baseURLParsed,
		Timeout: cardsDataModels.Timeout,
	}

	return config
}
