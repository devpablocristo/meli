package configurations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/configclient/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_configService_LoadParametersValidation_Success(t *testing.T) {
	ctx := context.TODO()
	simulator := configclient.SimulateCache(profilesLocalDirectory, configclient.TTLDefault)

	mockService := New(simulator)
	mockService.ResetTTLDefaultCache(1 * time.Second)

	parametersValidation, err := mockService.LoadParametersValidation(ctx)

	parametersValidationExpected := mockConfigParametersValidationExpected(t)

	assert.NotNil(t, parametersValidation.RulesSite["MLM"].Blockeduser)
	assert.NotNil(t, parametersValidation.RulesSite["MLM"].EvaluationUser)
	assert.NotNil(t, parametersValidation.RulesSite["MLM"].MaxAmountReparation)
	assert.NotNil(t, parametersValidation.RulesSite["MLM"].QtyReparationPerPeriodDays)
	assert.NotNil(t, parametersValidation.RulesSite["MLM"].StatusDetailAllowed)
	assert.NotNil(t, parametersValidation.RulesSite["MLM"].UserLifetime)
	assert.Nil(t, parametersValidation.RulesSite["MLM"].QtyChargebackPerPeriodDays)
	assert.Nil(t, parametersValidation.RulesSite["MLM"].SubTypeAllowed)
	assert.Nil(t, parametersValidation.RulesSite["MLM"].MerchantBlockedlist)

	assert.NotNil(t, parametersValidation.RulesSite["MLA"].Blockeduser)
	assert.NotNil(t, parametersValidation.RulesSite["MLA"].EvaluationUser)
	assert.NotNil(t, parametersValidation.RulesSite["MLA"].MaxAmountReparation)
	assert.NotNil(t, parametersValidation.RulesSite["MLA"].QtyChargebackPerPeriodDays)
	assert.NotNil(t, parametersValidation.RulesSite["MLA"].QtyReparationPerPeriodDays)
	assert.NotNil(t, parametersValidation.RulesSite["MLA"].StatusDetailAllowed)
	assert.NotNil(t, parametersValidation.RulesSite["MLA"].SubTypeAllowed)
	assert.NotNil(t, parametersValidation.RulesSite["MLA"].MerchantBlockedlist)
	assert.Nil(t, parametersValidation.RulesSite["MLA"].UserLifetime)
	utilstest.AssertFieldsNoEmptyFromStruct(t, *parametersValidation.RulesSite["MLA"])

	assert.NotNil(t, parametersValidation.RulesSite["MLC"].QtyChargebackPerPeriodDays)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].Blockeduser)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].EvaluationUser)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].MaxAmountReparation)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].QtyReparationPerPeriodDays)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].StatusDetailAllowed)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].SubTypeAllowed)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].MerchantBlockedlist)
	assert.Nil(t, parametersValidation.RulesSite["MLC"].UserLifetime)

	assert.NoError(t, err)
	assert.Equal(t, parametersValidationExpected, parametersValidation)
}

func Test_configService_LoadParametersValidation_Error(t *testing.T) {
	ctx := context.TODO()
	simulator := configclient.SimulateCache("non-existent-location", configclient.TTLDefault)
	mockService := New(simulator)

	parametersValidation, err := mockService.LoadParametersValidation(ctx)

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("open non-existent-location: no such file or directory | ID: parameters-validation"),
		"error in configurations", false)

	assert.Nil(t, parametersValidation)
	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_configService_Error_LoadApplicationInit_NotFound(t *testing.T) {
	ctx := context.TODO()
	simulator := configclient.SimulateCache("testdataerr", configclient.TTLDefault)
	mockService := New(simulator)

	applicationInit, err := mockService.LoadApplicationInit(ctx)

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("open testdataerr: no such file or directory | ID: app-init"),
		"error in configurations", false)

	assert.Nil(t, applicationInit)
	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_configService_Error_LoadParametersValidation_NotFound(t *testing.T) {
	ctx := context.TODO()
	simulator := configclient.SimulateCache("testdataerr", configclient.TTLDefault)
	mockService := New(simulator)

	applicationInit, err := mockService.LoadParametersValidation(ctx)

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("open testdataerr: no such file or directory | ID: parameters-validation"),
		"error in configurations", false)

	assert.Nil(t, applicationInit)
	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
