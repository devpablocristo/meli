package ds

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_dsRepository_Save_Success_Eligible(t *testing.T) {

	input := mockValidationResultEligible(t)

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultEligible)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "eligible")

			return nil
		},
	}

	ds := New(mockDsClientService)

	err := ds.Save(context.TODO(), "50844833020", input)

	assert.NoError(t, err)
}

func Test_dsRepository_Save_Error_Unknown(t *testing.T) {
	input := mockValidationResultEligible(t)

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			return errors.New("some unknown error")
		},
	}

	ds := New(mockDsClientService)

	err := ds.Save(context.TODO(), "50844833020", input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`some unknown error`), "error saving validation result", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_dsRepository_Save_Success_NotEligible_Full(t *testing.T) {
	input := mockValidationResultNotEligibleFull(t)

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleFull)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "not_eligible_full")

			return nil
		},
	}

	ds := New(mockDsClientService)

	utilstest.AssertFieldsNoEmptyFromStruct(t, input, "ApprovedData")
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)

	v := reflect.ValueOf(domain.ConfigRulesSite{})
	assert.Equal(t, v.NumField(), len(input.Reason), "new test cases are missing")
}

func Test_dsRepository_Save_Success_NotEligible_One_By_One(t *testing.T) {
	all := map[string]func(t *testing.T){
		"StatusDetailAllowed":        dsRepository_Save_Success_NotEligible_StatusNotAllowed,
		"QtyReparationPerPeriodDays": dsRepository_Save_Success_NotEligible_QtyRepairPeriod,
		"MaxAmountReparation":        dsRepository_Save_Success_NotEligible_MaxAmount,
		"EvaluationUser":             dsRepository_Save_Success_NotEligible_Restrictions,
		"QtyChargebackPerPeriodDays": dsRepository_Save_Success_NotEligible_QtyChargebackPeriod,
		"Blockeduser":                dsRepository_Save_Success_NotEligible_Blocked,
		"MerchantBlockedlist":        dsRepository_Save_Success_NotEligible_MerchantBlockedlist,
		"SubTypeAllowed":             dsRepository_Save_Success_NotEligible_SubTypeAllowed,
		"UserLifetime":               dsRepository_Save_Success_NotEligible_UserLifetime,
		"MinTotalHistoricalAmount":   dsRepository_Save_Success_NotEligible_MinTotalHistoricalAmount,
		"MinTransactionsQty":         dsRepository_Save_Success_NotEligible_MinTransactionsQty,
		"Score":                      dsRepository_Save_Success_NotEligible_Score,
	}

	v := reflect.ValueOf(domain.ConfigRulesSite{})
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		function, found := all[name]
		assert.True(t, found, "not_found_test_"+name)
		if found {
			function(t)
		}
	}
}

func dsRepository_Save_Success_NotEligible_StatusNotAllowed(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleStatusDetailAllowed: inputFull.Reason[domain.RuleStatusDetailAllowed]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleStatusNotAllowed)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "status_not_allowed")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_MaxAmount(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleMaxAmountReparation: inputFull.Reason[domain.RuleMaxAmountReparation]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleMaxAmount)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "max_eligible")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

type RuleQtyPerPeriodMaxAllowed struct {
	Qty        int
	PeriodDays int
}

func dsRepository_Save_Success_NotEligible_QtyRepairPeriod(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleQtyReparationPerPeriodDays: inputFull.Reason[domain.RuleQtyReparationPerPeriodDays]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleQtyRepairPeriod)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "qty_repair_period")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_Blocked(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleBlockeduser: inputFull.Reason[domain.RuleBlockeduser]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleBlocked)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "blocked")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_Restrictions(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleRestrictions: inputFull.Reason[domain.RuleRestrictions]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleRestrictions)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "restrictions")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_QtyChargebackPeriod(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleQtyChargebackPerPeriodDays: inputFull.Reason[domain.RuleQtyChargebackPerPeriodDays]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleQtyChargebackPeriod)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "qty_chargeback_period")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_MerchantBlockedlist(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleMerchantBlockedlist: inputFull.Reason[domain.RuleMerchantBlockedlist]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleRuleMerchantBlockedlist)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "merchant_blockedlist")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_SubTypeAllowed(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleSubTypeAllowed: inputFull.Reason[domain.RuleSubTypeAllowed]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleSubTypeNotAllowed)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "subtype_allowed")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_UserLifetime(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleUserLifetime: inputFull.Reason[domain.RuleUserLifetime]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleUserLifetime)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "user_lifetime")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_MinTotalHistoricalAmount(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleMinTotalHistoricalAmount: inputFull.Reason[domain.RuleMinTotalHistoricalAmount]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleMinTotalHistoricalAmount)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "min_total_historical_amount")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_MinTransactionsQty(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleMinTransactionsQty: inputFull.Reason[domain.RuleMinTransactionsQty]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleMinTransactionsQty)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "min_transactions_qty")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}

func dsRepository_Save_Success_NotEligible_Score(t *testing.T) {
	inputFull := mockValidationResultNotEligibleFull(t)
	input := inputFull
	input.Reason = domain.Reason{domain.RuleScore: inputFull.Reason[domain.RuleScore]}

	mockDsClientService := mockDsClientService{
		SaveDocumentWithContextStub: func(ctx context.Context, key string, value interface{}) error {
			valueExpected := string(_resultNotEligibleScore)

			bValue, err := json.Marshal(value)
			assert.NoError(t, err)

			assert.Equal(t, "50844833020", key)
			assert.JSONEq(t, valueExpected, string(bValue), "score")

			return nil
		},
	}

	ds := New(mockDsClientService)
	err := ds.Save(context.TODO(), "50844833020", input)
	assert.NoError(t, err)
}
