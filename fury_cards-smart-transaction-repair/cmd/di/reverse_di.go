package di

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/handlers/validationforcepublishhdl"

	configurationService "github.com/melisource/cards-smart-transaction-repair/adapters/configurations"
	dsblockeduser "github.com/melisource/cards-smart-transaction-repair/adapters/dbblockeduser/ds"
	dsreparation "github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation/ds"
	"github.com/melisource/cards-smart-transaction-repair/adapters/graphqlhub"
	"github.com/melisource/cards-smart-transaction-repair/adapters/score"
	"github.com/melisource/cards-smart-transaction-repair/adapters/topics/topicvalidation"

	"github.com/melisource/cards-smart-transaction-repair/adapters/cardstransactions"
	kvsblockeduser "github.com/melisource/cards-smart-transaction-repair/adapters/dbblockeduser/kvs"
	kvsreparation "github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation/kvs"
	dsValidation "github.com/melisource/cards-smart-transaction-repair/adapters/dbvalidationresult/ds"
	"github.com/melisource/cards-smart-transaction-repair/adapters/graphqlkyc"
	policyService "github.com/melisource/cards-smart-transaction-repair/adapters/policy"
	transactionscons "github.com/melisource/cards-smart-transaction-repair/consumers/transactions"
	validationresultcons "github.com/melisource/cards-smart-transaction-repair/consumers/validationresult"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/services/blockeduser"
	"github.com/melisource/cards-smart-transaction-repair/core/services/reparation"
	"github.com/melisource/cards-smart-transaction-repair/core/services/validation"
	"github.com/melisource/cards-smart-transaction-repair/handlers/blockeduserhdl"
	"github.com/melisource/cards-smart-transaction-repair/handlers/reversehdl"
	"github.com/melisource/fury_cards-go-toolkit/pkg/configclient/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/ioutil/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
	kvsclient "github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/policyclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/publisher/v1"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchhubclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchkycclient/v1"
)

const (
	applicationName = "cards-smart-transaction-repair"
)

type cmds struct {
	ReverseHandler                *reversehdl.ReverseHandler
	BlockedUserHandler            *blockeduserhdl.BlockedUserHandler
	TransactionsConsumer          *transactionscons.TransactionsConsumer
	ValidationResultConsumer      *validationresultcons.ValidationResultConsumer
	ValidationForcePublishHandler *validationforcepublishhdl.ValidationForcePublishHandler
}

func ConfigReverseDI() cmds {
	scope := strings.ToUpper(os.Getenv("SCOPE"))

	switch {
	case scope == "":
		return loadAppMock()
	default:
		return loadApp()
	}
}

func loadApp() cmds {
	logService := log.New(log.WithStacktraceOnError(false))

	var configurations ports.Configurations
	var cfgApplicationInit *domain.ConfigApplicationInit
	{
		ctx := context.TODO()
		var err error
		configurations = configurationService.New(configclient.NewConfigServiceCache(configclient.TTLDefault))

		cfgApplicationInit, err = configurations.LoadApplicationInit(ctx)
		if err != nil {
			panic(err)
		}
		configurations.ResetTTLDefaultCache(cfgApplicationInit.Options.TimeoutHoursCacheDefault)

		cfgParametersValidation, err := configurations.LoadParametersValidation(ctx)
		if err != nil {
			panic(err)
		}

		logService.Info(ctx, "configurations",
			logService.Any("initial_settings", cfgApplicationInit),
			logService.Any("parameters_validation", cfgParametersValidation))
	}

	var reparationRepo ports.ReparationRepository
	{
		configKVSRepair := kvsclient.NewConfig(cfgApplicationInit.Infra.KVSReparation.Name).
			WithReadRetries(cfgApplicationInit.Infra.KVSReparation.ReadRetries).
			WithReadTimeout(cfgApplicationInit.Infra.KVSReparation.ReadTimeout).
			WithWriteRetries(cfgApplicationInit.Infra.KVSReparation.WriteRetries).
			WithWriteTimeout(cfgApplicationInit.Infra.KVSReparation.WriteTimeout)

		kvsCliReparation, err := kvsclient.New(configKVSRepair)
		if err != nil {
			panic(err)
		}
		reparationRepo = kvsreparation.New(kvsCliReparation, logService)
	}

	var blockedUserRepo ports.BlockedUserRepository
	{
		configKVSBlockedlist := kvsclient.NewConfig(cfgApplicationInit.Infra.KVSBlockedlist.Name).
			WithReadRetries(cfgApplicationInit.Infra.KVSBlockedlist.ReadRetries).
			WithReadTimeout(cfgApplicationInit.Infra.KVSBlockedlist.ReadTimeout).
			WithWriteRetries(cfgApplicationInit.Infra.KVSBlockedlist.WriteRetries).
			WithWriteTimeout(cfgApplicationInit.Infra.KVSBlockedlist.WriteTimeout)

		kvsCliBlockedlist, err := kvsclient.New(configKVSBlockedlist)
		if err != nil {
			panic(err)
		}
		blockedUserRepo = kvsblockeduser.New(kvsCliBlockedlist)
	}

	var validationResultRepo ports.ValidationResultRepository
	{
		configDSValidation := dsclient.NewConfig(cfgApplicationInit.Infra.DSValidation.Name).
			WithReadTimeout(cfgApplicationInit.Infra.DSValidation.ReadTimeout).
			WithWriteTimeout(cfgApplicationInit.Infra.DSValidation.WriteTimeout).
			WithRetry(dsclient.NewRetryStrategy(cfgApplicationInit.Infra.DSValidation.Retries))

		validationResultRepo = dsValidation.New(dsclient.New(configDSValidation))
	}

	var reparationSearchRepo ports.ReparationSearchRepository
	{
		configDSRepair := dsclient.NewConfig(cfgApplicationInit.Infra.DSReparation.Name).
			WithReadTimeout(cfgApplicationInit.Infra.DSReparation.ReadTimeout).
			WithWriteTimeout(cfgApplicationInit.Infra.DSReparation.WriteTimeout).
			WithRetry(dsclient.NewRetryStrategy(cfgApplicationInit.Infra.DSReparation.Retries))

		reparationSearchRepo = dsreparation.New(dsclient.New(configDSRepair), dsclient.NewQuery())
	}

	var policy ports.Policy
	{
		policyCli, err := policyclient.New(applicationName, os.Getenv("SCOPE"))
		if err != nil {
			panic(err)
		}
		policy = policyService.New(policyCli)
	}

	var searchHub ports.SearchHub
	{
		baseURLTransactionsSearch, err := url.Parse(cfgApplicationInit.Infra.URLTransactionsSearch)
		if err != nil {
			panic(err)
		}
		configSearchHubCli := searchhubclient.NewConfigTracebility(baseURLTransactionsSearch, cfgApplicationInit.Options.TimeoutSecondsHTTP)
		searchHubCli := searchhubclient.New(configSearchHubCli)
		searchHub = graphqlhub.New(searchHubCli)
	}

	var blockedUserSearchRepository ports.BlockedUserSearchRepository
	{
		configDSBlocked := dsclient.NewConfig(cfgApplicationInit.Infra.DSBlockedUser.Name).
			WithReadTimeout(cfgApplicationInit.Infra.DSBlockedUser.ReadTimeout).
			WithWriteTimeout(cfgApplicationInit.Infra.DSBlockedUser.WriteTimeout).
			WithRetry(dsclient.NewRetryStrategy(cfgApplicationInit.Infra.DSBlockedUser.Retries))

		blockedUserSearchRepository = dsblockeduser.New(dsclient.New(configDSBlocked), dsclient.NewQuery())
	}

	var eventValidationResult ports.EventValidationResult
	{
		config := publisher.NewConfig(cfgApplicationInit.Infra.Topics.ValidationResult.Name).
			WithRetries(cfgApplicationInit.Infra.Topics.ValidationResult.Retries).
			WithTimeout(cfgApplicationInit.Infra.Topics.ValidationResult.Timeout)

		publisher, err := publisher.New(config)
		if err != nil {
			panic(err)
		}

		eventValidationResult = topicvalidation.New(publisher)
	}

	var scoreService ports.Score
	{
		httpCustom := rest.NewHTTPWithTracebilityAndTimeout(cfgApplicationInit.Infra.CardsDataModels.Timeout)
		config := score.NewConfig(cfgApplicationInit.Infra.CardsDataModels)
		scoreService = score.New(config, httpCustom, ioutil.New(), json.NewJSON(), logService)
	}

	blockedUserService := blockeduser.New(reparationSearchRepo, blockedUserRepo, blockedUserSearchRepository)

	validationService := validation.New(
		logService,
		configurations,
		validationResultRepo,
		reparationSearchRepo,
		blockedUserService,
		policy,
		searchHub,
		eventValidationResult,
		scoreService,
	)

	reparationService := newReparationService(
		cfgApplicationInit,
		logService,
		validationService,
		reparationSearchRepo,
		reparationRepo,
		searchHub,
	)

	reverseHdl := reversehdl.NewHandler(reparationService, logService)

	blockedUserHdl := blockeduserhdl.NewHandler(blockedUserService, logService)

	transactionsConsumer := transactionscons.NewConsumer(blockedUserService, logService)

	validationResultConsumer := validationresultcons.NewConsumer(validationService, logService)

	validationforcepublishhdl := validationforcepublishhdl.NewHandler(validationService, logService)

	return cmds{
		ReverseHandler:                reverseHdl,
		BlockedUserHandler:            blockedUserHdl,
		TransactionsConsumer:          transactionsConsumer,
		ValidationResultConsumer:      validationResultConsumer,
		ValidationForcePublishHandler: validationforcepublishhdl,
	}
}

func loadAppMock() cmds {
	reparationRepo := kvsreparation.NewMock()
	validationResultRepo := dsValidation.NewMock()
	cardsTransactions := cardstransactions.NewMock()
	configurationsServiceApp := configurationService.NewMock()
	blockedUserRepo := kvsblockeduser.NewMock()
	transactionRepairSearchRepo := dsreparation.NewMock()
	policyMock := policyService.NewMock()
	searchHubAppMock := graphqlhub.NewMock()
	blockedUserSearchRepo := dsblockeduser.NewMock()
	eventValidationResult := topicvalidation.NewMock()
	scoreService := score.NewMock()

	blockedUserService := blockeduser.New(transactionRepairSearchRepo, blockedUserRepo, blockedUserSearchRepo)
	validationService := validation.New(
		log.New(),
		configurationsServiceApp,
		validationResultRepo,
		transactionRepairSearchRepo,
		blockedUserService,
		policyMock,
		searchHubAppMock,
		eventValidationResult,
		scoreService,
	)

	reverseService := reparation.New(
		log.New(),
		validationService,
		reparationRepo,
		searchHubAppMock,
		cardsTransactions,
		dsreparation.NewMock(),
		graphqlkyc.NewMock(),
	)

	reverseHdl := reversehdl.NewHandler(reverseService, log.New())
	blockedUserHdl := blockeduserhdl.NewHandler(blockedUserService, log.New())
	transactionsConsumer := transactionscons.NewConsumer(blockedUserService, log.New())
	validationResultConsumer := validationresultcons.NewConsumer(validationService, log.New())
	validationforcepublishhdl := validationforcepublishhdl.NewHandler(validationService, log.New())

	return cmds{
		ReverseHandler:                reverseHdl,
		BlockedUserHandler:            blockedUserHdl,
		TransactionsConsumer:          transactionsConsumer,
		ValidationResultConsumer:      validationResultConsumer,
		ValidationForcePublishHandler: validationforcepublishhdl,
	}
}

func newReparationService(
	cfgApplicationInit *domain.ConfigApplicationInit,
	logService log.LogService,
	validationService ports.ValidationService,
	transactionRepairSearchRepo ports.ReparationSearchRepository,
	reparationRepo ports.ReparationRepository,
	searchHub ports.SearchHub,
) ports.ReparationService {
	httpCustom := rest.NewHTTPWithTracebilityAndTimeout(cfgApplicationInit.Options.TimeoutSecondsHTTP)

	configReverse := cardstransactions.NewConfig(cfgApplicationInit.Infra.TransactionsApp)

	transactionApp := cardstransactions.New(configReverse, httpCustom, ioutil.New(), json.NewJSON(), logService)

	var searchKyc ports.SearchKyc
	{
		baseURLSearchKyc, err := url.Parse(cfgApplicationInit.Infra.GraphqlKyc.URL)
		if err != nil {
			panic(err)
		}
		searchKycAppConfig := searchkycclient.NewConfigTracebility(
			baseURLSearchKyc,
			cfgApplicationInit.Options.TimeoutSecondsHTTP,
			cfgApplicationInit.Infra.GraphqlKyc.PegasusToken,
			cfgApplicationInit.Infra.GraphqlKyc.APISandbox,
		)
		searchKycAppClient := searchkycclient.New(searchKycAppConfig)
		searchKyc = graphqlkyc.New(searchKycAppClient)
	}

	return reparation.New(
		logService,
		validationService,
		reparationRepo,
		searchHub,
		transactionApp,
		transactionRepairSearchRepo,
		searchKyc,
	)
}
