package validation

import (
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_runPrerequisitesForRuleValidation_Success(t *testing.T) {
	t.Run("local currency", func(t *testing.T) {
		rulesSite := domain.ConfigRulesSite{
			MaxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     2000,
				DecimalDigits: 0,
				Currency:      "MXN",
			},
		}

		transactionData := domain.TransactionData{
			Operation: domain.Operation{
				Settlement: domain.Settlement{Amount: 10, Currency: "USD"},
				Billing:    domain.Billing{Amount: 100, Currency: "MXN"},
			},
			SiteID:  "MLM",
			PayerID: 123,
		}
		transactionData.Operation.Billing.Currency = "MXN"

		input := domain.ValidationInput{
			SiteID:          "MLM",
			UserID:          123,
			TransactionData: transactionData,
		}

		err := runPrerequisitesForRuleValidation(input, rulesSite)
		assert.NoError(t, err)
	})

	t.Run("international currency", func(t *testing.T) {
		rulesSite := domain.ConfigRulesSite{
			MaxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     19,
				DecimalDigits: 0,
				Currency:      "USD",
			},
		}

		transactionData := domain.TransactionData{
			Operation: domain.Operation{
				Settlement: domain.Settlement{Amount: 10, Currency: "USD"},
				Billing:    domain.Billing{Amount: 100, Currency: "MXN"},
			},
			SiteID:  "MLA",
			PayerID: 123,
		}
		transactionData.Operation.Settlement.Currency = "USD"

		input := domain.ValidationInput{
			SiteID:          "MLA",
			UserID:          123,
			TransactionData: transactionData,
		}

		err := runPrerequisitesForRuleValidation(input, rulesSite)
		assert.NoError(t, err)
	})

}

func Test_runPrerequisitesForRuleValidation_SiteDifferent(t *testing.T) {
	rulesSite := domain.ConfigRulesSite{
		MaxAmountReparation: &domain.ConfigMaxAmountReparation{
			MaxAmount:     2000,
			DecimalDigits: 0,
			Currency:      "MXN",
		},
	}

	t.Run("should return different siteID error", func(t *testing.T) {
		transactionData := domain.TransactionData{
			Operation: domain.Operation{},
			SiteID:    "MLB",
			PayerID:   123,
		}

		input := domain.ValidationInput{
			SiteID:          "MLM",
			UserID:          123,
			TransactionData: transactionData,
		}

		err := runPrerequisitesForRuleValidation(input, rulesSite)

		errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
			errors.New(`received siteID(MLM) different from database(MLB)`), "validation failed", false)

		assert.Error(t, err)
		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	})
}

func Test_runPrerequisitesForRuleValidation_RuleMaxRepairAmount_CurrencyDifferent(t *testing.T) {
	t.Run("should return error empty billing", func(t *testing.T) {
		rulesSite := domain.ConfigRulesSite{
			MaxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     2000,
				DecimalDigits: 0,
				Currency:      "MXN",
			},
		}

		transactionData := domain.TransactionData{
			Operation: domain.Operation{
				Settlement: domain.Settlement{Amount: 10, Currency: "USD"},
			},
			SiteID:  "MLM",
			PayerID: 123,
		}

		transactionData.Operation.Billing.Currency = "USD"

		input := domain.ValidationInput{
			SiteID:          "MLM",
			UserID:          123,
			TransactionData: transactionData,
		}

		err := runPrerequisitesForRuleValidation(input, rulesSite)

		errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
			errors.New(`billing is empty in the transactional database`), "validation failed", false)

		assert.Error(t, err)
		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	})

	t.Run("should return error empty settlement", func(t *testing.T) {
		rulesSite := domain.ConfigRulesSite{
			MaxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     2000,
				DecimalDigits: 0,
				Currency:      "MXN",
			},
		}

		transactionData := domain.TransactionData{
			Operation: domain.Operation{
				Billing: domain.Billing{Amount: 10, Currency: "USD"},
			},
			SiteID:  "MLM",
			PayerID: 123,
		}

		transactionData.Operation.Billing.Currency = "USD"

		input := domain.ValidationInput{
			SiteID:          "MLM",
			UserID:          123,
			TransactionData: transactionData,
		}

		err := runPrerequisitesForRuleValidation(input, rulesSite)

		errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
			errors.New(`settlement is empty in the transactional database`), "validation failed", false)

		assert.Error(t, err)
		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	})
	t.Run("should return error for different local currency", func(t *testing.T) {
		rulesSite := domain.ConfigRulesSite{
			MaxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     2000,
				DecimalDigits: 0,
				Currency:      "MXN",
			},
		}

		transactionData := domain.TransactionData{
			Operation: domain.Operation{
				Settlement: domain.Settlement{Amount: 10, Currency: "USD"},
				Billing:    domain.Billing{Amount: 100, Currency: "MXN"},
			},
			SiteID:  "MLM",
			PayerID: 123,
		}

		transactionData.Operation.Billing.Currency = "USD"

		input := domain.ValidationInput{
			SiteID:          "MLM",
			UserID:          123,
			TransactionData: transactionData,
		}

		err := runPrerequisitesForRuleValidation(input, rulesSite)

		errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
			errors.New(`transaction billing currency(USD) different from the parameters(MXN)`), "validation failed", false)

		assert.Error(t, err)
		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	})

	t.Run("should return error for different international currency", func(t *testing.T) {
		rulesSite := domain.ConfigRulesSite{
			MaxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     16,
				DecimalDigits: 0,
				Currency:      "USD",
			},
		}

		transactionData := domain.TransactionData{
			Operation: domain.Operation{
				Settlement: domain.Settlement{Amount: 10, Currency: "USD"},
				Billing:    domain.Billing{Amount: 100, Currency: "MXN"},
			},
			SiteID:  "MLM",
			PayerID: 123,
		}

		transactionData.Operation.Settlement.Currency = "EUR"

		input := domain.ValidationInput{
			SiteID:          "MLM",
			UserID:          123,
			TransactionData: transactionData,
		}

		err := runPrerequisitesForRuleValidation(input, rulesSite)

		errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
			errors.New(`transaction settlement currency(EUR) different from the parameters(USD)`), "validation failed", false)

		assert.Error(t, err)
		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	})
}

func Test_runPrerequisitesForRuleValidation_Without_RuleMaxRepairAmount(t *testing.T) {
	rulesSite := domain.ConfigRulesSite{}

	input := domain.ValidationInput{}

	err := runPrerequisitesForRuleValidation(input, rulesSite)
	assert.NoError(t, err)
}

func Test_runPrerequisitesForRuleValidation_RuleAllowedStatus_StatusDetailEmpty(t *testing.T) {
	rulesSite := domain.ConfigRulesSite{
		StatusDetailAllowed: map[string]struct{}{
			"pending_capture": {},
		},
	}

	input := domain.ValidationInput{}

	err := runPrerequisitesForRuleValidation(input, rulesSite)

	errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
		errors.New(`status detail transaction is empty`), "validation failed", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_runPrerequisitesForRuleValidation_RuleAllowedSubType_SubTypeEmpty(t *testing.T) {
	rulesSite := domain.ConfigRulesSite{
		SubTypeAllowed: map[string]struct{}{
			"purchase": {},
		},
	}

	input := domain.ValidationInput{}

	err := runPrerequisitesForRuleValidation(input, rulesSite)

	errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
		errors.New(`subtype transaction is empty`), "validation failed", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_runPrerequisitesForRuleValidation_RuleMerchantBlockedlist_CardAcceptorNameEmpty(t *testing.T) {
	rulesSite := domain.ConfigRulesSite{
		MerchantBlockedlist: []string{
			"123componentes",
		},
	}

	input := domain.ValidationInput{}

	err := runPrerequisitesForRuleValidation(input, rulesSite)

	errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
		errors.New(`merchant name is empty`), "validation failed", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_runPrerequisitesForRuleValidation_RuleUserLifetime_DatecreatedNil(t *testing.T) {
	rulesSite := domain.ConfigRulesSite{
		UserLifetime: &domain.ConfigUserLifetime{
			GteDays: 10,
		},
	}

	input := domain.ValidationInput{}

	err := runPrerequisitesForRuleValidation(input, rulesSite)

	errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
		errors.New(`date created user is empty`), "validation failed", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
