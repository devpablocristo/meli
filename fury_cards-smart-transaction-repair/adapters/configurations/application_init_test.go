package configurations

import (
	"context"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/configclient/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

const (
	profilesLocalDirectory = "testdata"
)

func Test_configService_LoadApplicationInit_Success(t *testing.T) {
	ctx := context.TODO()
	simulator := configclient.SimulateCache(profilesLocalDirectory, configclient.TTLDefault)
	mockService := New(simulator)

	applicationInit, err := mockService.LoadApplicationInit(ctx)

	applicationInitExpected := mockConfigApplicationInitExpected(t)

	assert.NoError(t, err)
	assert.Equal(t, applicationInitExpected, applicationInit)

	utilstest.AssertFieldsNoEmptyFromStruct(t, applicationInit.Infra)
	utilstest.AssertFieldsNoEmptyFromStruct(t, applicationInit.Options)
}

func Test_configService_LoadApplicationInit_Error(t *testing.T) {
	ctx := context.TODO()
	simulator := configclient.SimulateCache("non-existent-location", configclient.TTLDefault)
	mockService := New(simulator)

	applicationInit, err := mockService.LoadApplicationInit(ctx)

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("open non-existent-location: no such file or directory | ID: app-init"),
		"error in configurations", false)

	assert.Nil(t, applicationInit)
	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_treatError(t *testing.T) {
	err := treatError(errors.New("error unknown"), "profile")

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("error unknown"),
		"error in configurations", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_toConfigBaseURL_Errors(t *testing.T) {
	t.Run("transactions_app", func(t *testing.T) {
		baseURL := transactionsApp{
			BaseURL: "%",
		}
		assert.Panics(t, func() {
			toConfigTransactionsApp(baseURL)
		})

		baseURL = transactionsApp{
			BaseURL:       "https://www.",
			BaseURLLegacy: "%",
		}
		assert.Panics(t, func() {
			toConfigTransactionsApp(baseURL)
		})
	})

	t.Run("cardsDataModels", func(t *testing.T) {
		baseURL := cardsDataModels{
			BaseURL: "%",
		}
		assert.Panics(t, func() {
			toConfigCardsDataModels(baseURL)
		})
	})
}
