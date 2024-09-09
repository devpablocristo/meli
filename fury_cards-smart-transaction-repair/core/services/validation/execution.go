package validation

import (
	"context"
	"math"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

func (v validationService) ExecuteValidation(
	ctx context.Context,
	input domain.ValidationInput,
	wrapFields *field.WrappedFields,
) error {
	loadedData, err := v.loadData(ctx, input, wrapFields)
	if err != nil {
		return err
	}

	evaluatedData, err := v.validateEligibilityForReparation(ctx, input, loadedData, wrapFields)
	if err != nil {
		return err
	}

	transactionEvaluated := fillTransactionEvaluated(input, loadedData.RulesSite.MaxAmountReparation)

	validationResult := fillValidationResult(input.PaymentID, evaluatedData, input, transactionEvaluated, loadedData.CalculatedScore)

	v.PublishEvent(ctx, validationResult, wrapFields)

	if !validationResult.IsEligible {
		return buildErrorNotEligible(validationResult.Reason, input.TransactionData.Operation.CreationDatetime)
	}

	return nil
}

func (v validationService) validateEligibilityForReparation(
	ctx context.Context,
	input domain.ValidationInput,
	loadedData *loadedData,
	wrapFields *field.WrappedFields,
) (*domain.EvaluatedData, error) {
	err := runPrerequisitesForRuleValidation(input, loadedData.RulesSite)
	if err != nil {
		addMetricInvalidPrerequisitesForRuleValidation(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, err
	}

	return v.validateRulesAndGetReasons(ctx, input, loadedData, wrapFields), nil
}

func (v validationService) validateRulesAndGetReasons(
	ctx context.Context,
	input domain.ValidationInput,
	loadedData *loadedData,
	wrapFields *field.WrappedFields,
) *domain.EvaluatedData {

	appliedRules := ruleLog{}
	evaluatedData := &domain.EvaluatedData{
		UnapprovedData: domain.Reason{},
		ApprovedData:   domain.ApprovedData{},
	}

	assertMaxRepairAmount(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMaxQtyRepairPerPeriod(ctx, evaluatedData, input, loadedData, appliedRules)
	assertAllowedStatus(ctx, evaluatedData, input, loadedData, appliedRules)
	assertUserIsBlocked(ctx, evaluatedData, input, loadedData, appliedRules)
	assertUserRestriction(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMaxQtyChargebackPerPeriod(ctx, evaluatedData, input, loadedData, appliedRules)
	assertAllowedSubType(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMerchantBlockedlist(ctx, evaluatedData, input, loadedData, appliedRules)
	assertUserLifetime(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMinTotalHistoricalAmount(ctx, evaluatedData, input, loadedData, appliedRules)
	assertMinTransactionsQty(ctx, evaluatedData, input, loadedData, appliedRules)
	assertScore(ctx, evaluatedData, input, loadedData, appliedRules)

	wrapFields.Fields.Add(string(domain.KeyFieldRulesApplied), appliedRules)

	if len(evaluatedData.UnapprovedData) == 0 {
		evaluatedData.UnapprovedData = nil
	}

	if len(evaluatedData.ApprovedData) == 0 {
		evaluatedData.ApprovedData = nil
	}

	return evaluatedData
}

func fillValidationResult(
	paymentID string,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	transactionEvaluated storage.TransactionEvaluated,
	score *float64,
) *storage.ValidationResult {
	return &storage.ValidationResult{
		PaymentID:            paymentID,
		IsEligible:           evaluatedData.UnapprovedData == nil,
		DateCreated:          time.Now().UTC(),
		AuthorizationID:      input.TransactionData.Operation.Authorization.ID,
		Type:                 input.Type,
		KycIdentificationID:  input.UserKycData.KycIdentificationID,
		UserID:               input.UserID,
		SiteID:               input.TransactionData.SiteID,
		FaqID:                input.FaqID,
		Requester:            input.Requester,
		TransactionEvaluated: transactionEvaluated,
		Score:                score,
		ApprovedData:         evaluatedData.ApprovedData,
		Reason:               evaluatedData.UnapprovedData,
	}
}

func fillTransactionEvaluated(
	input domain.ValidationInput,
	configMaxAmountReparation *domain.ConfigMaxAmountReparation,
) storage.TransactionEvaluated {
	var transactionEvaluated storage.TransactionEvaluated

	if configMaxAmountReparation == nil {
		return transactionEvaluated
	}

	switch configMaxAmountReparation.Currency {
	case internationalCurrency:
		decimalizedAmount :=
			input.TransactionData.Operation.Settlement.Amount / math.Pow10(input.TransactionData.Operation.Settlement.DecimalDigits)

		transactionEvaluated = storage.TransactionEvaluated{
			Amount:   decimalizedAmount,
			Currency: input.TransactionData.Operation.Settlement.Currency,
		}

	default:
		decimalizedAmount :=
			input.TransactionData.Operation.Billing.Amount / math.Pow10(input.TransactionData.Operation.Billing.DecimalDigits)

		transactionEvaluated = storage.TransactionEvaluated{
			Amount:   decimalizedAmount,
			Currency: input.TransactionData.Operation.Billing.Currency,
		}
	}

	return transactionEvaluated
}
