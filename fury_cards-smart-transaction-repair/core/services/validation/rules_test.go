package validation

import (
	"context"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/stretchr/testify/assert"
)

func Test_rules_not_applicable(t *testing.T) {
	ctx := context.TODO()

	input := domain.ValidationInput{
		TransactionData: domain.TransactionData{},
		PaymentID:       "100",
		SiteID:          "MLB",
	}

	t.Run("max_qty_repair_per_period not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.QtyReparationPerPeriodDays = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertMaxQtyRepairPerPeriod(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleQtyReparationPerPeriodDays])
		assert.Nil(t, loadedData.RulesSite.QtyReparationPerPeriodDays)
		assert.Nil(t, evaluatedData.UnapprovedData, 0)
	})

	t.Run("user_is_blocked not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.Blockeduser = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertUserIsBlocked(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleBlockeduser])
		assert.Nil(t, loadedData.RulesSite.Blockeduser)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("max_repair_amount not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.MaxAmountReparation = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleMaxAmountReparation])
		assert.Nil(t, loadedData.RulesSite.MaxAmountReparation)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("allowed_status not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.StatusDetailAllowed = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertAllowedStatus(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleStatusDetailAllowed])
		assert.Nil(t, loadedData.RulesSite.StatusDetailAllowed)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("user_restriction not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.EvaluationUser = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertUserRestriction(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleRestrictions])
		assert.Nil(t, loadedData.RulesSite.EvaluationUser)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("max_qty_chargeback_per_period not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.QtyChargebackPerPeriodDays = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertMaxQtyChargebackPerPeriod(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleQtyChargebackPerPeriodDays])
		assert.Nil(t, loadedData.RulesSite.QtyChargebackPerPeriodDays)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("merchant_blockedlist not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.MerchantBlockedlist = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertMerchantBlockedlist(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleMerchantBlockedlist])
		assert.Nil(t, loadedData.RulesSite.MerchantBlockedlist)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("allowed_subtype not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
			PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
			SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		}
		loadedData.RulesSite.SubTypeAllowed = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertAllowedSubType(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleSubTypeAllowed])
		assert.Nil(t, loadedData.RulesSite.SubTypeAllowed)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("user_lifetime not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite: *mockParametersValidation(t).RulesSite["MLB"],
		}
		loadedData.RulesSite.UserLifetime = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleUserLifetime])
		assert.Nil(t, loadedData.RulesSite.UserLifetime)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("min_total_historical_amount not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite: *mockParametersValidation(t).RulesSite["MLB"],
		}
		loadedData.RulesSite.MinTotalHistoricalAmount = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertMinTotalHistoricalAmount(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleMinTotalHistoricalAmount])
		assert.Nil(t, loadedData.RulesSite.MinTotalHistoricalAmount)
		assert.Len(t, appliedRules, 0)
	})

	t.Run("min_transactions_qty not applicable", func(t *testing.T) {
		reason := domain.Reason{}
		loadedData := &loadedData{
			RulesSite: *mockParametersValidation(t).RulesSite["MLB"],
		}
		loadedData.RulesSite.MinTransactionsQty = nil
		appliedRules := ruleLog{}
		evaluatedData := &domain.EvaluatedData{}

		assertMinTransactionsQty(ctx, evaluatedData, input, loadedData, appliedRules)

		assert.Nil(t, reason[domain.RuleMinTransactionsQty])
		assert.Nil(t, loadedData.RulesSite.MinTransactionsQty)
		assert.Len(t, appliedRules, 0)
	})
}

func Test_rules_applicable(t *testing.T) {
	ctx := context.TODO()

	input := domain.ValidationInput{
		TransactionData: domain.TransactionData{},
		PaymentID:       "100",
		SiteID:          "MLB",
	}

	loadedData := &loadedData{
		RulesSite:                      *mockParametersValidation(t).RulesSite["MLB"],
		PolicyOutput:                   &domain.PolicyOutput{IsAuthorized: true},
		SearchHubCountChargebackOutput: &domain.SearchHubCountChargebackOutput{},
		TransactionHistory:             &domain.TransactionHistoryOutput{},
		CalculatedScore:                new(float64),
	}

	appliedRules := ruleLog{}
	evaluatedData := &domain.EvaluatedData{
		UnapprovedData: domain.Reason{},
		ApprovedData:   domain.ApprovedData{},
	}

	assertMaxQtyRepairPerPeriod(ctx, evaluatedData, input, loadedData, appliedRules)
	assertAllowedStatus(ctx, evaluatedData, input, loadedData, appliedRules)
	assertUserIsBlocked(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
	assertUserRestriction(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMinTotalHistoricalAmount(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMaxQtyChargebackPerPeriod(ctx, evaluatedData, input, loadedData, appliedRules)
	assertAllowedSubType(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMerchantBlockedlist(ctx, evaluatedData, input, loadedData, appliedRules)
	assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)
	assertScore(ctx, evaluatedData, input, loadedData, appliedRules)

	assert.Contains(t, appliedRules["100|MLB"], domain.RuleQtyReparationPerPeriodDays+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleStatusDetailAllowed+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleBlockeduser+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleMaxAmountReparation+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleRestrictions+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleMinTotalHistoricalAmount+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleQtyChargebackPerPeriodDays+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleSubTypeAllowed+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleMerchantBlockedlist+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleUserLifetime+";")
	assert.Contains(t, appliedRules["100|MLB"], domain.RuleScore+";")
}

func Test_assertMaxRepairAmount_Accept(t *testing.T) {
	ctx := context.TODO()

	lteDays := new(int)
	*lteDays = 14
	trxAgeRanges := []domain.ConfigMaxAmountReparationTrxAgeRange{
		{
			GteDays:              7,
			LteDays:              lteDays,
			MaxAccumulatedAmount: 3000,
			DecimalDigits:        0,
			Currency:             "MXN",
			Score: &domain.ConfigScoreMaxAmountReparationTrxAgeRange{
				MinScoreValue:                   0.8,
				MaxAccumulatedAmountDecimalized: 4000,
			},
		},
		{
			GteDays:              15,
			MaxAccumulatedAmount: 5000,
			DecimalDigits:        0,
			Currency:             "MXN",
			Score: &domain.ConfigScoreMaxAmountReparationTrxAgeRange{
				MinScoreValue:                   0.9,
				MaxAccumulatedAmountDecimalized: 6000,
			},
		},
	}

	tests := []struct {
		name                  string
		configMaxAmount       float64
		configDecimalDigits   int
		accumulatedAmount     float64
		decimalDigits         int
		currency              string
		paymentAge            int
		trxAgeRanges          []domain.ConfigMaxAmountReparationTrxAgeRange
		scoreMaxAmountDefault domain.ConfigScoreMaxAmountReparation
		calculatedScore       float64
	}{
		{
			name:                "should accept - (max 2000 and 0 digits) (transaction 1999,99999 and 5 digits)",
			configMaxAmount:     2000,
			configDecimalDigits: 0,
			accumulatedAmount:   199999999,
			decimalDigits:       5,
			currency:            "USD",
		},
		{
			name:                "should accept, (max 2000,00 and 2 digits) (transaction 1999,99999 and 5 digits)",
			configMaxAmount:     200000,
			configDecimalDigits: 2,
			accumulatedAmount:   199999999,
			decimalDigits:       5,
			currency:            "MXN",
		},
		{
			name:                "should accept, (max 2000,00 and 2 digits) (transaction 1999,99 and 2 digits)",
			configMaxAmount:     200000,
			configDecimalDigits: 2,
			accumulatedAmount:   199999,
			decimalDigits:       2,
			currency:            "MXN",
		},
		{
			name:                "should accept - where maxAmount:2900 MXN and trxAge: 7 days and with Ranges",
			configMaxAmount:     2000,
			configDecimalDigits: 0,
			accumulatedAmount:   2900,
			decimalDigits:       2,
			currency:            "MXN",
			paymentAge:          7,
		},
		{
			name:                "should accept - where maxAmount:3500 MXN and trxAge: 15 days and with Ranges",
			configMaxAmount:     2000,
			configDecimalDigits: 0,
			accumulatedAmount:   3500,
			decimalDigits:       0,
			currency:            "MXN",
			trxAgeRanges:        trxAgeRanges,
			paymentAge:          15,
		},
		{
			name:                  "should accept - where maxAmount:3200 MXN - payment age 1: validated via score",
			configMaxAmount:       2000,
			configDecimalDigits:   0,
			accumulatedAmount:     3200,
			decimalDigits:         0,
			currency:              "MXN",
			paymentAge:            1,
			calculatedScore:       0.7,
			scoreMaxAmountDefault: domain.ConfigScoreMaxAmountReparation{MinScoreValue: 0.7, MaxAmountDecimalized: 3500.00},
		},
		{
			name:                  "should accept - where maxAmount:4000 MXN - payment age 8 - - with ranges: validated via score",
			configMaxAmount:       2000,
			configDecimalDigits:   0,
			accumulatedAmount:     4000,
			decimalDigits:         0,
			currency:              "MXN",
			paymentAge:            8,
			calculatedScore:       0.8,
			trxAgeRanges:          trxAgeRanges,
			scoreMaxAmountDefault: domain.ConfigScoreMaxAmountReparation{MinScoreValue: 0.7, MaxAmountDecimalized: 3500.00},
		},
		{
			name:                  "should accept - where maxAmount:4900 MXN - payment age 20 - with ranges: validated via score",
			configMaxAmount:       2000,
			configDecimalDigits:   0,
			accumulatedAmount:     4900,
			decimalDigits:         0,
			currency:              "MXN",
			paymentAge:            20,
			calculatedScore:       0.9,
			trxAgeRanges:          trxAgeRanges,
			scoreMaxAmountDefault: domain.ConfigScoreMaxAmountReparation{MinScoreValue: 0.7, MaxAmountDecimalized: 3500.00},
		},
		{
			name:                  "should accept - where maxAmount:12 USD - payment age 1: validated via score",
			configMaxAmount:       10,
			configDecimalDigits:   0,
			accumulatedAmount:     12,
			decimalDigits:         0,
			currency:              "USD",
			paymentAge:            1,
			calculatedScore:       0.8,
			scoreMaxAmountDefault: domain.ConfigScoreMaxAmountReparation{MinScoreValue: 0.7, MaxAmountDecimalized: 15.00},
		},
		{
			name:                "should accept - where maxAmount:22 USD - payment age 8 - with ranges: validated via score",
			configMaxAmount:     10,
			configDecimalDigits: 0,
			accumulatedAmount:   22,
			decimalDigits:       0,
			currency:            "USD",
			paymentAge:          8,
			calculatedScore:     0.9,
			trxAgeRanges: []domain.ConfigMaxAmountReparationTrxAgeRange{
				{
					GteDays:              8,
					MaxAccumulatedAmount: 20,
					DecimalDigits:        0,
					Currency:             "USD",
					Score: &domain.ConfigScoreMaxAmountReparationTrxAgeRange{
						MinScoreValue:                   0.9,
						MaxAccumulatedAmountDecimalized: 25,
					},
				},
			},
			scoreMaxAmountDefault: domain.ConfigScoreMaxAmountReparation{MinScoreValue: 0.7, MaxAmountDecimalized: 15.00},
		},
	}
	for _, tt := range tests {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite:       domain.ConfigRulesSite{},
			CalculatedScore: &tt.calculatedScore,
		}

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Billing:    domain.Billing{Amount: tt.accumulatedAmount, DecimalDigits: tt.decimalDigits},
				Settlement: domain.Settlement{Amount: tt.accumulatedAmount, DecimalDigits: tt.decimalDigits},
				PaymentAge: tt.paymentAge,
			},
		}

		loadedData.RulesSite.MaxAmountReparation = &domain.ConfigMaxAmountReparation{
			MaxAmount:     tt.configMaxAmount,
			DecimalDigits: tt.configDecimalDigits,
			Currency:      tt.currency,
			TrxAgeRanges:  tt.trxAgeRanges,
			Score:         &tt.scoreMaxAmountDefault,
		}

		input := domain.ValidationInput{
			TransactionData: *transactionData,
		}

		assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.Nil(t, evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation], tt.name)
		assert.Len(t, evaluatedData.ApprovedData, 1)
	}
}

func Test_assertMaxRepairAmount_Not_Accept(t *testing.T) {
	ctx := context.TODO()

	t.Run("should NOT accept - (max 2000 and 0 digits) (transaction 2000,00001 and 5 digits)", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{},
		}

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Billing: domain.Billing{Amount: 200000001, DecimalDigits: 5},
			},
			SiteID: "MLM",
		}

		loadedData.RulesSite.MaxAmountReparation = &domain.ConfigMaxAmountReparation{
			MaxAmount:     2000,
			DecimalDigits: 0,
			Currency:      "MXN",
		}

		input := domain.ValidationInput{
			TransactionData: *transactionData,
		}

		assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation])
	})

	t.Run("should NOT accept, (max 2000 and 2 digits) (transaction 2000,00001 and 5 digits)", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{},
		}

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Settlement: domain.Settlement{Amount: 200000001, DecimalDigits: 5},
			},
			SiteID: "USD",
		}

		loadedData.RulesSite.MaxAmountReparation = &domain.ConfigMaxAmountReparation{
			MaxAmount:     200000,
			DecimalDigits: 2,
			Currency:      "USD",
		}

		input := domain.ValidationInput{
			TransactionData: *transactionData,
		}

		assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation])
		assert.Len(t, evaluatedData.ApprovedData, 0)
	})

	t.Run("should NOT accept, (max 2000 and 2 digits) (transaction 2000,01 and 2 digits)", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{},
		}

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Billing: domain.Billing{Amount: 200001, DecimalDigits: 2},
			},
			SiteID: "MLM",
		}

		loadedData.RulesSite.MaxAmountReparation = &domain.ConfigMaxAmountReparation{
			MaxAmount:     200000,
			DecimalDigits: 2,
			Currency:      "MXN",
		}

		input := domain.ValidationInput{
			TransactionData: *transactionData,
		}

		assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation])
	})

	t.Run("should NOT accept, max amount from score OK, score value no NOK", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		calculatedScore := 0.2

		loadedData := &loadedData{
			RulesSite:       domain.ConfigRulesSite{},
			CalculatedScore: &calculatedScore,
		}

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Billing:    domain.Billing{Amount: 3500, DecimalDigits: 0},
				Settlement: domain.Settlement{Amount: 3500, DecimalDigits: 0},
			},
		}

		siteCurrency := map[int]string{
			1: "MXN",
			2: "USD",
		}

		for _, v := range siteCurrency {
			loadedData.RulesSite.MaxAmountReparation = &domain.ConfigMaxAmountReparation{
				MaxAmount: 3000,
				Currency:  v,
				Score: &domain.ConfigScoreMaxAmountReparation{
					MinScoreValue:        0.7,
					MaxAmountDecimalized: 4000,
				},
			}

			input := domain.ValidationInput{
				TransactionData: *transactionData,
			}

			assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
			assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation])

			evaluatedData.UnapprovedData = domain.Reason{}
		}
	})

	t.Run("should NOT accept, score value OK, max amount from score NOK", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		calculatedScore := 0.8

		loadedData := &loadedData{
			RulesSite:       domain.ConfigRulesSite{},
			CalculatedScore: &calculatedScore,
		}

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Billing:    domain.Billing{Amount: 4001, DecimalDigits: 0},
				Settlement: domain.Settlement{Amount: 4001, DecimalDigits: 0},
			},
		}

		siteCurrency := map[int]string{
			1: "MXN",
			2: "USD",
		}

		for _, v := range siteCurrency {
			loadedData.RulesSite.MaxAmountReparation = &domain.ConfigMaxAmountReparation{
				MaxAmount: 3000,
				Currency:  v,
				Score: &domain.ConfigScoreMaxAmountReparation{
					MinScoreValue:        0.7,
					MaxAmountDecimalized: 4000,
				},
			}

			input := domain.ValidationInput{
				TransactionData: *transactionData,
			}

			assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
			assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation])

			evaluatedData.UnapprovedData = domain.Reason{}
		}
	})
}

func Test_assertMaxRepairAmount_WithRanges_Billing(t *testing.T) {
	ctx := context.TODO()

	currencyMXN := "MXN"
	lteDays := new(int)
	*lteDays = 14
	maxAmountReparation := &domain.ConfigMaxAmountReparation{
		MaxAmount:     1000,
		DecimalDigits: 0,
		Currency:      currencyMXN,
		TrxAgeRanges: []domain.ConfigMaxAmountReparationTrxAgeRange{
			{GteDays: 7, LteDays: lteDays, MaxAccumulatedAmount: 2000, DecimalDigits: 0, Currency: currencyMXN},
			{GteDays: 15, MaxAccumulatedAmount: 3000, DecimalDigits: 0, Currency: currencyMXN},
		},
	}

	tests := []struct {
		name                string
		amount              float64
		paymentAge          int
		maxAmountReparation *domain.ConfigMaxAmountReparation
		repairs             []storage.ReparationOutput
		notAccept           bool
	}{
		{
			name:                "should NOT accept - where amount:3001,00 MXN and trxAge: 15 days",
			paymentAge:          15,
			amount:              3001,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
		},
		{
			name:                "should accept - where amount:3000,00 MXN and trxAge: 15 days",
			paymentAge:          15,
			amount:              3000,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where amount:2001,00 MXN and trxAge: 10 days",
			paymentAge:          10,
			amount:              2001,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
		},
		{
			name:                "should accept - where amount:2000,00 MXN and trxAge: 10 days",
			paymentAge:          10,
			amount:              2000,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where amount:1001,00 MXN and trxAge: 6 days",
			paymentAge:          6,
			amount:              1001,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
		},
		{
			name:                "should accept - where amount:1000,00 MXN and trxAge: 6 days",
			paymentAge:          6,
			amount:              1000,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where amount:100,00 MXN and trxAge: 7 days and with 1 repair",
			paymentAge:          7,
			amount:              100,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Billing: storage.MonetaryValue{Amount: 1901, Currency: currencyMXN}}},
				},
			},
		},
		{
			name:                "should accept - where amount:99,00 MXN and trxAge: 7 days and with 1 repair",
			paymentAge:          7,
			amount:              99,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Billing: storage.MonetaryValue{Amount: 1901, Currency: currencyMXN}}},
				},
			},
		},
		{
			name:                "should NOT accept - where amount:101,00 MXN and trxAge: 7 days and with 2 repairs",
			paymentAge:          7,
			amount:              101,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Billing: storage.MonetaryValue{Amount: 1100, Currency: currencyMXN}}},
				},
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Billing: storage.MonetaryValue{Amount: 800, Currency: currencyMXN}}},
				},
			},
		},
		{
			name:                "should accept - where amount:100,00 MXN and trxAge: 7 days and with 2 repairs",
			paymentAge:          7,
			amount:              100,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Billing: storage.MonetaryValue{Amount: 1100, Currency: currencyMXN}}},
				},
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Billing: storage.MonetaryValue{Amount: 800, Currency: currencyMXN}}},
				},
			},
		},
		{
			name:       "should NOT accept - default (without ranges)",
			paymentAge: 15,
			amount:     1001,
			notAccept:  true,
			maxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     1000,
				DecimalDigits: 0,
				Currency:      currencyMXN,
			},
		},
		{
			name:       "should accept - default (without ranges)",
			paymentAge: 15,
			amount:     1000,
			notAccept:  false,
			maxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     1000,
				DecimalDigits: 0,
				Currency:      currencyMXN,
			},
		},
	}
	for _, tt := range tests {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{},
		}

		decimalDigits := 2
		amount := tt.amount * math.Pow10(decimalDigits)

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Billing:    domain.Billing{Amount: amount, DecimalDigits: decimalDigits},
				PaymentAge: tt.paymentAge,
			},
			SiteID: "MLM",
		}

		input := domain.ValidationInput{
			TransactionData: *transactionData,
		}

		loadedData.RulesSite.MaxAmountReparation = tt.maxAmountReparation
		loadedData.Repairs = tt.repairs

		assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		_, found := evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation]
		assert.Equal(t, found, tt.notAccept, tt.name)
	}
}

func Test_assertMaxRepairAmount_WithRanges_Settlement(t *testing.T) {
	ctx := context.TODO()

	currencyUSD := "USD"
	lteDays := new(int)
	*lteDays = 14
	maxAmountReparation := &domain.ConfigMaxAmountReparation{
		MaxAmount:     100,
		DecimalDigits: 0,
		Currency:      currencyUSD,
		TrxAgeRanges: []domain.ConfigMaxAmountReparationTrxAgeRange{
			{GteDays: 7, LteDays: lteDays, MaxAccumulatedAmount: 200, DecimalDigits: 0, Currency: currencyUSD},
			{GteDays: 15, MaxAccumulatedAmount: 300, DecimalDigits: 0, Currency: currencyUSD},
		},
	}

	tests := []struct {
		name                string
		amount              float64
		paymentAge          int
		maxAmountReparation *domain.ConfigMaxAmountReparation
		repairs             []storage.ReparationOutput
		notAccept           bool
	}{
		{
			name:                "should NOT accept - where amount:301,00 USD and trxAge: 15 days",
			paymentAge:          15,
			amount:              301,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
		},
		{
			name:                "should accept - where amount:300,00 USD and trxAge: 15 days",
			paymentAge:          15,
			amount:              300,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where amount:201,00 USD and trxAge: 10 days",
			paymentAge:          10,
			amount:              201,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
		},
		{
			name:                "should accept - where amount:200,00 USD and trxAge: 10 days",
			paymentAge:          10,
			amount:              200,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where amount:101,00 USD and trxAge: 6 days",
			paymentAge:          6,
			amount:              101,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
		},
		{
			name:                "should accept - where amount:100,00 USD and trxAge: 6 days",
			paymentAge:          6,
			amount:              100,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where amount:10,00 USD and trxAge: 7 days and with 1 repair",
			paymentAge:          7,
			amount:              10,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Settlement: storage.MonetaryValue{Amount: 191, Currency: currencyUSD}}},
				},
			},
		},
		{
			name:                "should accept - where amount:9,00 USD and trxAge: 7 days and with 1 repair",
			paymentAge:          7,
			amount:              9,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Settlement: storage.MonetaryValue{Amount: 191, Currency: currencyUSD}}},
				},
			},
		},
		{
			name:                "should NOT accept - where amount:10,00 USD and trxAge: 7 days and with 2 repairs",
			paymentAge:          7,
			amount:              10,
			maxAmountReparation: maxAmountReparation,
			notAccept:           true,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Settlement: storage.MonetaryValue{Amount: 100, Currency: currencyUSD}}},
				},
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Settlement: storage.MonetaryValue{Amount: 91, Currency: currencyUSD}}},
				},
			},
		},
		{
			name:                "should accept - where amount:9,00 USD and trxAge: 7 days and with 2 repairs",
			paymentAge:          7,
			amount:              9,
			maxAmountReparation: maxAmountReparation,
			notAccept:           false,
			repairs: []storage.ReparationOutput{
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Settlement: storage.MonetaryValue{Amount: 100, Currency: currencyUSD}}},
				},
				{
					BaseReparation: storage.BaseReparation{
						MonetaryTransaction: storage.MonetaryTransaction{Settlement: storage.MonetaryValue{Amount: 91, Currency: currencyUSD}}},
				},
			},
		},
		{
			name:       "should NOT accept - default (without ranges)",
			paymentAge: 15,
			amount:     101,
			notAccept:  true,
			maxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     100,
				DecimalDigits: 0,
				Currency:      currencyUSD,
			},
		},
		{
			name:       "should accept - default (without ranges)",
			paymentAge: 15,
			amount:     100,
			notAccept:  false,
			maxAmountReparation: &domain.ConfigMaxAmountReparation{
				MaxAmount:     100,
				DecimalDigits: 0,
				Currency:      currencyUSD,
			},
		},
	}
	for _, tt := range tests {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{},
		}

		decimalDigits := 2
		amount := tt.amount * math.Pow10(decimalDigits)

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				Settlement: domain.Settlement{Amount: amount, DecimalDigits: decimalDigits},
				PaymentAge: tt.paymentAge,
			},
			SiteID: "MLA",
		}

		input := domain.ValidationInput{
			TransactionData: *transactionData,
		}

		loadedData.RulesSite.MaxAmountReparation = tt.maxAmountReparation
		loadedData.Repairs = tt.repairs

		assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		_, found := evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation]
		assert.Equal(t, found, tt.notAccept, tt.name)
	}
}

func Test_assertUserLifetime(t *testing.T) {
	ctx := context.TODO()

	t.Run("should accept", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				UserLifetime: &domain.ConfigUserLifetime{
					GteDays: 10,
				},
			},
		}

		input := domain.ValidationInput{
			UserKycData: domain.SearchKycOutput{
				DateCreated: time.Date(2023, 3, 7, 0, 0, 0, 0, time.UTC),
			},
		}
		assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.Nil(t, evaluatedData.UnapprovedData[domain.RuleUserLifetime])
	})

	t.Run("should NOT accept", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				UserLifetime: &domain.ConfigUserLifetime{
					GteDays: 10,
				},
			},
		}

		input := domain.ValidationInput{
			UserKycData: domain.SearchKycOutput{
				DateCreated: time.Now().UTC().AddDate(0, 0, loadedData.RulesSite.UserLifetime.GteDays-1),
			},
		}
		assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleUserLifetime])
	})

	t.Run("should accept BECAUSE score", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}
		minScore := 0.3
		calculatedScore := 0.35

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				UserLifetime: &domain.ConfigUserLifetime{
					GteDays: 10,
					Score: &domain.ConfigScoreUserLifetime{
						MinScoreValue: minScore,
						GteDays:       -100,
					},
				},
			},
			CalculatedScore: &calculatedScore,
		}

		input := domain.ValidationInput{
			UserKycData: domain.SearchKycOutput{
				DateCreated: time.Now().UTC().AddDate(0, 0, loadedData.RulesSite.UserLifetime.GteDays-1),
			},
		}
		assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.Nil(t, evaluatedData.UnapprovedData[domain.RuleUserLifetime])
	})

	t.Run("should NOT accept, even with a good score", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}
		minScore := 0.3
		calculatedScore := 0.35

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				UserLifetime: &domain.ConfigUserLifetime{
					GteDays: 10,
					Score: &domain.ConfigScoreUserLifetime{
						MinScoreValue: minScore,
						GteDays:       5,
					},
				},
			},
			CalculatedScore: &calculatedScore,
		}

		input := domain.ValidationInput{
			UserKycData: domain.SearchKycOutput{
				DateCreated: time.Now().UTC().AddDate(0, 0, loadedData.RulesSite.UserLifetime.GteDays-1),
			},
		}
		assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleUserLifetime])
	})

	t.Run("should NOT accept low score", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}
		minScore := 0.3
		calculatedScore := 0.29

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				UserLifetime: &domain.ConfigUserLifetime{
					GteDays: 10,
					Score: &domain.ConfigScoreUserLifetime{
						MinScoreValue: minScore,
					},
				},
			},
			CalculatedScore: &calculatedScore,
		}

		input := domain.ValidationInput{
			UserKycData: domain.SearchKycOutput{
				DateCreated: time.Now().UTC().AddDate(0, 0, loadedData.RulesSite.UserLifetime.GteDays-1),
			},
		}
		assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleUserLifetime])
	})
}

func Test_assertMinTotalHistoricalAmount(t *testing.T) {
	ctx := context.TODO()

	t.Run("should accept - Settlement", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				MinTotalHistoricalAmount: &domain.ConfigMinTotalHistoricalAmount{MinAmount: 20, Currency: "USD", DecimalDigits: 0},
			},
			TransactionHistory: &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{
						MatchTotal: 2,
						Transactions: []domain.HistoryOperation{
							{Settlement: domain.Settlement{Amount: 10, Currency: "USD", DecimalDigits: 0}},
							{Settlement: domain.Settlement{Amount: 8, Currency: "USD", DecimalDigits: 0}},
						},
					},
					{
						MatchTotal: 1,
						Transactions: []domain.HistoryOperation{
							{Settlement: domain.Settlement{Amount: 11, Currency: "USD", DecimalDigits: 0}},
						},
					},
				},
			},
		}

		input := domain.ValidationInput{}

		assertMinTotalHistoricalAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.Nil(t, evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount])
	})

	t.Run("should accept - Billing", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				MinTotalHistoricalAmount: &domain.ConfigMinTotalHistoricalAmount{MinAmount: 100, Currency: "MXN", DecimalDigits: 0},
			},
			TransactionHistory: &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{
						MatchTotal: 4,
						Transactions: []domain.HistoryOperation{
							{Billing: domain.Billing{Amount: 10, Currency: "MXN", DecimalDigits: 0}},
							{Billing: domain.Billing{Amount: 20, Currency: "MXN", DecimalDigits: 0}},
							{Billing: domain.Billing{Amount: 60, Currency: "MXN", DecimalDigits: 0}},
							{Billing: domain.Billing{Amount: 5, Currency: "MXN", DecimalDigits: 0}},
						},
					},
					{
						MatchTotal: 2,
						Transactions: []domain.HistoryOperation{
							{Billing: domain.Billing{Amount: 2, Currency: "MXN", DecimalDigits: 0}},
							{Billing: domain.Billing{Amount: 3, Currency: "MXN", DecimalDigits: 0}},
						},
					},
				},
			},
		}

		input := domain.ValidationInput{}

		assertMinTotalHistoricalAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.Nil(t, evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount])
	})

	t.Run("should NOT accept - Settlement", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				MinTotalHistoricalAmount: &domain.ConfigMinTotalHistoricalAmount{MinAmount: 100, Currency: "USD", DecimalDigits: 0},
			},
			TransactionHistory: &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{
						MatchTotal: 2,
						Transactions: []domain.HistoryOperation{
							{Settlement: domain.Settlement{Amount: 80, Currency: "USD", DecimalDigits: 0}},
							{Settlement: domain.Settlement{Amount: 10, Currency: "USD", DecimalDigits: 0}},
						},
					},
					{
						MatchTotal: 2,
						Transactions: []domain.HistoryOperation{
							{Settlement: domain.Settlement{Amount: 5, Currency: "USD", DecimalDigits: 0}},
							{Settlement: domain.Settlement{Amount: 4, Currency: "USD", DecimalDigits: 0}},
						},
					},
				},
			},
		}

		input := domain.ValidationInput{}

		assertMinTotalHistoricalAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount])
		assert.Equal(t, float64(99), evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount].Actual)
		assert.Equal(t, float64(100), evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount].Accepted)
	})

	t.Run("should NOT accept - Billing", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				MinTotalHistoricalAmount: &domain.ConfigMinTotalHistoricalAmount{MinAmount: 20, Currency: "MXN", DecimalDigits: 0},
			},
			TransactionHistory: &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{
						MatchTotal: 3,
						Transactions: []domain.HistoryOperation{
							{Billing: domain.Billing{Amount: 2, Currency: "MXN", DecimalDigits: 0}},
							{Billing: domain.Billing{Amount: 2, Currency: "MXN", DecimalDigits: 0}},
							{Billing: domain.Billing{Amount: 2, Currency: "MXN", DecimalDigits: 0}},
						},
					},
					{
						MatchTotal: 2,
						Transactions: []domain.HistoryOperation{
							{Billing: domain.Billing{Amount: 2, Currency: "MXN", DecimalDigits: 0}},
							{Billing: domain.Billing{Amount: 2, Currency: "MXN", DecimalDigits: 0}},
						},
					},
				},
			},
		}

		input := domain.ValidationInput{}

		assertMinTotalHistoricalAmount(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount])
		assert.Equal(t, float64(10), evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount].Actual)
		assert.Equal(t, float64(20), evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount].Accepted)
	})
}

func Test_assertMinTransactionsQty(t *testing.T) {
	ctx := context.TODO()

	t.Run("should accept", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				MinTransactionsQty: &domain.ConfigMinTransactionsQty{MinQty: 18},
			},
			TransactionHistory: &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{MatchTotal: 10},
					{MatchTotal: 8},
				},
			},
		}

		input := domain.ValidationInput{}

		assertMinTransactionsQty(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.Nil(t, evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount])
	})

	t.Run("should NOT accept", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				MinTransactionsQty: &domain.ConfigMinTransactionsQty{MinQty: 5},
			},
			TransactionHistory: &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{MatchTotal: 1},
					{MatchTotal: 3},
				},
			},
		}

		input := domain.ValidationInput{}

		assertMinTransactionsQty(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleMinTransactionsQty])
		assert.Equal(t, 4, evaluatedData.UnapprovedData[domain.RuleMinTransactionsQty].Actual)
		assert.Equal(t, 5, evaluatedData.UnapprovedData[domain.RuleMinTransactionsQty].Accepted)
	})
}

func Test_ExecuteMetrics(t *testing.T) {
	ctx := context.TODO()
	all := map[string]func(ctx context.Context, siteID, faqID, clientApplication string){
		"StatusDetailAllowed":        addMetricRuleStatusNotAllowed,
		"QtyReparationPerPeriodDays": addMetricRuleMaxQtyRepairPerPeriod,
		"MaxAmountReparation":        addMetricRuleMaxAmount,
		"EvaluationUser":             addMetricRuleUserWithRestriction,
		"QtyChargebackPerPeriodDays": addMetricRuleMaxQtyChargebackPerPeriod,
		"Blockeduser":                addMetricRuleUserIsBlocked,
		"MerchantBlockedlist":        addMetricRuleMerchantBlockedlist,
		"SubTypeAllowed":             addMetricRuleSubTypeNotAllowed,
		"UserLifetime":               addMetricRuleUserLifetime,
		"MinTotalHistoricalAmount":   addMetricRuleMinTotalHistoricalAmount,
		"MinTransactionsQty":         addMetricRuleMinTransactionsQty,
		"Score":                      addMetricRuleScore,
	}

	v := reflect.ValueOf(domain.ConfigRulesSite{})
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		function, found := all[name]
		assert.True(t, found, "not_found_metric_to: "+name)
		if found {
			function(ctx, "site_id", "faq_id", "app_test")
		}
	}
}

func Test_assertMaxQtyRepairPerPeriod_WithRanges(t *testing.T) {
	ctx := context.TODO()

	lteDays := new(int)
	*lteDays = 14
	minScore05 := 0.5
	minScore06 := 0.6
	minScore07 := 0.7
	qtyReparationPeriod := &domain.ConfigQtyReparationPeriod{
		Qty:        1,
		PeriodDays: 30,
		TrxAgeRanges: []domain.ConfigQtyReparationPeriodTrxAgeRange{
			{GteDays: 7, LteDays: lteDays, Qty: 2, Score: &domain.ConfigScoreQtyReparationPeriod{MinScoreValue: minScore06, Qty: 3}},
			{GteDays: 15, Qty: 3, Score: &domain.ConfigScoreQtyReparationPeriod{MinScoreValue: minScore07, Qty: 4}},
		},
		Score: &domain.ConfigScoreQtyReparationPeriod{MinScoreValue: minScore05},
	}

	tests := []struct {
		name                string
		paymentAge          int
		qtyReparationPeriod *domain.ConfigQtyReparationPeriod
		repairs             []storage.ReparationOutput
		notAccept           bool
		calculatedScore     float64
	}{
		{
			name:                "should NOT accept - where trxAge: 15 days, with 3 repairs",
			paymentAge:          15,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 3),
			notAccept:           true,
		},
		{
			name:                "should accept - where trxAge: 15 days, with 2 repairs",
			paymentAge:          15,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 2),
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where trxAge: 10 days, with 2 repairs",
			paymentAge:          10,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 2),
			notAccept:           true,
		},
		{
			name:                "should accept - where trxAge: 10 days, with 1 repair",
			paymentAge:          10,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 1),
			notAccept:           false,
		},
		{
			name:                "should NOT accept - where trxAge: 6 days, with 1 repair",
			paymentAge:          6,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 1),
			notAccept:           true,
		},
		{
			name:                "should NOT accept - where trxAge: 6 days, without repair",
			paymentAge:          6,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 0),
			notAccept:           false,
		},
		{
			name:       "should NOT accept - default (without ranges)",
			paymentAge: 15,
			repairs:    make([]storage.ReparationOutput, 1),
			notAccept:  true,
			qtyReparationPeriod: &domain.ConfigQtyReparationPeriod{
				Qty:        1,
				PeriodDays: 30,
			},
		},
		{
			name:       "should accept - default (without ranges)",
			paymentAge: 15,
			notAccept:  false,
			repairs:    make([]storage.ReparationOutput, 0),
			qtyReparationPeriod: &domain.ConfigQtyReparationPeriod{
				Qty:        1,
				PeriodDays: 30,
			},
		},
		{
			name:       "should accept because of the score  - default (without ranges)",
			paymentAge: 15,
			repairs:    make([]storage.ReparationOutput, 1),
			notAccept:  false,
			qtyReparationPeriod: &domain.ConfigQtyReparationPeriod{
				Qty:        1,
				PeriodDays: 30,
				Score: &domain.ConfigScoreQtyReparationPeriod{
					MinScoreValue: minScore05,
					Qty:           2,
				},
			},
			calculatedScore: 0.6,
		},
		{
			name:                "should accept because score - where trxAge: 10 days, with 2 repairs",
			paymentAge:          10,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 2),
			notAccept:           false,
			calculatedScore:     0.6,
		},
		{
			name:                "should accept because score - where trxAge: 20 days, with 3 repairs",
			paymentAge:          20,
			qtyReparationPeriod: qtyReparationPeriod,
			repairs:             make([]storage.ReparationOutput, 3),
			notAccept:           false,
			calculatedScore:     0.7,
		},
	}
	for _, tt := range tests {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite:       domain.ConfigRulesSite{},
			CalculatedScore: &tt.calculatedScore,
		}

		transactionData := &domain.TransactionData{
			Operation: domain.Operation{
				PaymentAge: tt.paymentAge,
			},
			SiteID: "MLM",
		}

		input := domain.ValidationInput{
			TransactionData: *transactionData,
		}

		loadedData.RulesSite.QtyReparationPerPeriodDays = tt.qtyReparationPeriod
		loadedData.Repairs = tt.repairs

		assertMaxQtyRepairPerPeriod(ctx, evaluatedData, input, loadedData, appliedRules)
		_, found := evaluatedData.UnapprovedData[domain.RuleQtyReparationPerPeriodDays]
		assert.Equal(t, found, tt.notAccept, tt.name)
	}
}

func Test_assertScore(t *testing.T) {
	ctx := context.Background()

	t.Run("should accept", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				Score: &domain.ConfigScore{
					MinScoreValue: 0.7,
				},
			},
		}

		calculatedScore := 0.7000
		if loadedData.RulesSite.CheckApplyAnyScore() {
			loadedData.CalculatedScore = &calculatedScore
		}

		input := domain.ValidationInput{}

		assertScore(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.Nil(t, evaluatedData.UnapprovedData[domain.RuleScore])
	})

	t.Run("should NOT accept", func(t *testing.T) {
		evaluatedData := &domain.EvaluatedData{
			UnapprovedData: domain.Reason{},
			ApprovedData:   domain.ApprovedData{},
		}
		appliedRules := ruleLog{}

		loadedData := &loadedData{
			RulesSite: domain.ConfigRulesSite{
				Score: &domain.ConfigScore{
					MinScoreValue: 0.7,
				},
			},
		}

		calculatedScore := 0.6999
		if loadedData.RulesSite.CheckApplyAnyScore() {
			loadedData.CalculatedScore = &calculatedScore
		}

		input := domain.ValidationInput{}

		assertScore(ctx, evaluatedData, input, loadedData, appliedRules)
		assert.NotNil(t, evaluatedData.UnapprovedData[domain.RuleScore])
	})
}
