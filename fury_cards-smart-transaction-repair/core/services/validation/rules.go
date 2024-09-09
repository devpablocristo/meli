package validation

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	// for international transactions, the currency used is USD.
	internationalCurrency                        = "USD"
	precisionRoundAmountAfterValidationMaxRepair = 8
	precisionRoundScore                          = 4
)

func assertMaxQtyRepairPerPeriod(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {
	// Checks if this rule applies for the siteID.
	configQtyReparationPeriod := loadedData.RulesSite.QtyReparationPerPeriodDays
	if configQtyReparationPeriod == nil {
		return
	}

	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleQtyReparationPerPeriodDays)

	paymentAge := input.TransactionData.Operation.PaymentAge
	for _, item := range configQtyReparationPeriod.TrxAgeRanges {
		if paymentAge >= item.GteDays && (item.LteDays == nil || paymentAge <= *item.LteDays) {
			configQtyReparationPeriod = &domain.ConfigQtyReparationPeriod{
				Qty:        item.Qty,
				PeriodDays: configQtyReparationPeriod.PeriodDays,
				Score:      item.Score,
			}
			break
		}
	}

	resultEvaluatedData := &domain.ReasonResult{
		Actual:   len(loadedData.Repairs),
		Accepted: domain.ReasonResultPerPeriod{Qty: configQtyReparationPeriod.Qty, PeriodDays: configQtyReparationPeriod.PeriodDays},
	}

	switch {
	case len(loadedData.Repairs) < configQtyReparationPeriod.Qty:
		evaluatedData.ApprovedData[domain.RuleQtyReparationPerPeriodDays] = resultEvaluatedData
		return

	case configQtyReparationPeriod.Score != nil && checkScore(*loadedData.CalculatedScore, configQtyReparationPeriod.Score.MinScoreValue):
		resultEvaluatedData.Accepted =
			domain.ReasonResultPerPeriod{Qty: configQtyReparationPeriod.Score.Qty, PeriodDays: configQtyReparationPeriod.PeriodDays}

		if len(loadedData.Repairs) < configQtyReparationPeriod.Score.Qty {
			evaluatedData.ApprovedData[domain.RuleQtyReparationPerPeriodDays] = resultEvaluatedData
			return
		}
	}

	evaluatedData.UnapprovedData[domain.RuleQtyReparationPerPeriodDays] = resultEvaluatedData
	addMetricRuleMaxQtyRepairPerPeriod(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertUserIsBlocked(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	if loadedData.RulesSite.Blockeduser == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleBlockeduser)

	resultEvaluatedData := &domain.ReasonResult{
		Actual: loadedData.UserIsBlocked,
	}

	if !loadedData.UserIsBlocked {
		evaluatedData.ApprovedData[domain.RuleBlockeduser] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleBlockeduser] = resultEvaluatedData

	addMetricRuleUserIsBlocked(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertMaxRepairAmount(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configMaxAmountReparation := loadedData.RulesSite.MaxAmountReparation
	if configMaxAmountReparation == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleMaxAmountReparation)

	var accumulatedAmountSettlement float64
	var accumulatedAmountBilling float64
	{
		paymentAge := input.TransactionData.Operation.PaymentAge
		for _, item := range configMaxAmountReparation.TrxAgeRanges {
			if paymentAge >= item.GteDays && (item.LteDays == nil || paymentAge <= *item.LteDays) {
				configMaxAmountReparation = &domain.ConfigMaxAmountReparation{
					MaxAmount:     item.MaxAccumulatedAmount,
					Currency:      item.Currency,
					DecimalDigits: item.DecimalDigits,
				}

				if item.Score != nil {
					configMaxAmountReparation.Score = &domain.ConfigScoreMaxAmountReparation{
						MinScoreValue:        item.Score.MinScoreValue,
						MaxAmountDecimalized: item.Score.MaxAccumulatedAmountDecimalized,
					}
				}

				for _, repair := range loadedData.Repairs {
					accumulatedAmountSettlement += repair.MonetaryTransaction.Settlement.Amount
					accumulatedAmountBilling += repair.MonetaryTransaction.Billing.Amount
				}
				break
			}
		}
	}

	maxAmount := configMaxAmountReparation.MaxAmount
	if configMaxAmountReparation.DecimalDigits > 0 {
		maxAmount /= math.Pow10(configMaxAmountReparation.DecimalDigits)
	}

	var amount float64
	{
		switch configMaxAmountReparation.Currency {
		case internationalCurrency:
			amount = input.TransactionData.Operation.Settlement.Amount / math.Pow10(input.TransactionData.Operation.Settlement.DecimalDigits)
			amount += accumulatedAmountSettlement
		default:
			amount = input.TransactionData.Operation.Billing.Amount / math.Pow10(input.TransactionData.Operation.Billing.DecimalDigits)
			amount += accumulatedAmountBilling
		}
	}

	resultEvaluatedData := &domain.ReasonResult{
		Actual:   round(amount, precisionRoundAmountAfterValidationMaxRepair),
		Accepted: maxAmount,
	}

	switch {
	case amount <= maxAmount:
		evaluatedData.ApprovedData[domain.RuleMaxAmountReparation] = resultEvaluatedData
		return

	case configMaxAmountReparation.Score != nil &&
		checkScore(*loadedData.CalculatedScore, configMaxAmountReparation.Score.MinScoreValue):

		resultEvaluatedData.Accepted = configMaxAmountReparation.Score.MaxAmountDecimalized
		if amount <= configMaxAmountReparation.Score.MaxAmountDecimalized {
			evaluatedData.ApprovedData[domain.RuleMaxAmountReparation] = resultEvaluatedData
			return
		}
	}

	evaluatedData.UnapprovedData[domain.RuleMaxAmountReparation] = resultEvaluatedData

	addMetricRuleMaxAmount(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertAllowedStatus(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configStatusDetailAllowed := loadedData.RulesSite.StatusDetailAllowed
	if configStatusDetailAllowed == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleStatusDetailAllowed)

	statusAllowed := []string{}
	for status := range loadedData.RulesSite.StatusDetailAllowed {
		statusAllowed = append(statusAllowed, status)
	}

	resultEvaluatedData := &domain.ReasonResult{
		Actual:   input.TransactionData.StatusDetail,
		Accepted: statusAllowed,
	}

	if _, found := configStatusDetailAllowed[input.TransactionData.StatusDetail]; found {
		evaluatedData.ApprovedData[domain.RuleStatusDetailAllowed] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleStatusDetailAllowed] = resultEvaluatedData

	addMetricRuleStatusNotAllowed(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertUserRestriction(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {
	// Checks if this rule applies for the siteID.
	if loadedData.RulesSite.EvaluationUser == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleRestrictions)

	resultEvaluatedData := &domain.ReasonResult{
		Actual: loadedData.PolicyOutput.RestrictsFailed,
	}

	if loadedData.PolicyOutput.IsAuthorized {
		evaluatedData.ApprovedData[domain.RuleRestrictions] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleRestrictions] = resultEvaluatedData

	addMetricRuleUserWithRestriction(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertMaxQtyChargebackPerPeriod(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configQtyChargebackPerPeriodDays := loadedData.RulesSite.QtyChargebackPerPeriodDays
	if configQtyChargebackPerPeriodDays == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleQtyChargebackPerPeriodDays)

	resultEvaluatedData := &domain.ReasonResult{
		Actual: loadedData.SearchHubCountChargebackOutput.Total,
		Accepted: domain.ReasonResultPerPeriod{
			Qty:        configQtyChargebackPerPeriodDays.Qty,
			PeriodDays: configQtyChargebackPerPeriodDays.PeriodDays,
		},
	}

	if loadedData.SearchHubCountChargebackOutput.Total <= configQtyChargebackPerPeriodDays.Qty {
		evaluatedData.ApprovedData[domain.RuleQtyChargebackPerPeriodDays] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleQtyChargebackPerPeriodDays] = resultEvaluatedData

	addMetricRuleMaxQtyChargebackPerPeriod(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertAllowedSubType(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configSubTypeAllowed := loadedData.RulesSite.SubTypeAllowed
	if configSubTypeAllowed == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleSubTypeAllowed)

	subTypeAllowed := []string{}
	for subtype := range loadedData.RulesSite.SubTypeAllowed {
		subTypeAllowed = append(subTypeAllowed, subtype)
	}

	resultEvaluatedData := &domain.ReasonResult{
		Actual:   input.TransactionData.Operation.SubType,
		Accepted: subTypeAllowed,
	}

	if _, found := configSubTypeAllowed[input.TransactionData.Operation.SubType]; found {
		evaluatedData.ApprovedData[domain.RuleSubTypeAllowed] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleSubTypeAllowed] = resultEvaluatedData

	addMetricRuleSubTypeNotAllowed(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertMerchantBlockedlist(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	merchantBlockedlist := loadedData.RulesSite.MerchantBlockedlist
	if merchantBlockedlist == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleMerchantBlockedlist)

	resultEvaluatedData := &domain.ReasonResult{
		Actual: input.TransactionData.Operation.Authorization.CardAcceptor.Name,
	}

	merchantAccepted := true
	for _, merchant := range merchantBlockedlist {
		if strings.Contains(input.TransactionData.Operation.Authorization.CardAcceptor.Name, merchant) {
			merchantAccepted = false
			break
		}
	}

	if merchantAccepted {
		evaluatedData.ApprovedData[domain.RuleMerchantBlockedlist] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleMerchantBlockedlist] = resultEvaluatedData

	addMetricRuleMerchantBlockedlist(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertUserLifetime(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configUserLifetime := loadedData.RulesSite.UserLifetime
	if configUserLifetime == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleUserLifetime)

	resultEvaluatedData := &domain.ReasonResult{
		Actual:   input.UserKycData.DateCreated,
		Accepted: configUserLifetime.GteDays,
	}

	minDateCreated := time.Now().UTC().AddDate(0, 0, -configUserLifetime.GteDays)

	switch {
	case input.UserKycData.DateCreated.Before(minDateCreated):
		evaluatedData.ApprovedData[domain.RuleUserLifetime] = resultEvaluatedData
		return

	case configUserLifetime.Score != nil && checkScore(*loadedData.CalculatedScore, configUserLifetime.Score.MinScoreValue):
		resultEvaluatedData.Accepted = configUserLifetime.Score.GteDays
		minDateCreated = time.Now().UTC().AddDate(0, 0, -configUserLifetime.Score.GteDays)
		if input.UserKycData.DateCreated.Before(minDateCreated) {
			evaluatedData.ApprovedData[domain.RuleUserLifetime] = resultEvaluatedData
			return
		}
	}

	evaluatedData.UnapprovedData[domain.RuleUserLifetime] = resultEvaluatedData
	addMetricRuleUserLifetime(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertMinTotalHistoricalAmount(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configMinTotalHistoricalAmount := loadedData.RulesSite.MinTotalHistoricalAmount
	if configMinTotalHistoricalAmount == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleMinTotalHistoricalAmount)

	minAmount := configMinTotalHistoricalAmount.MinAmount
	if configMinTotalHistoricalAmount.DecimalDigits > 0 {
		minAmount /= math.Pow10(configMinTotalHistoricalAmount.DecimalDigits)
	}

	var amount float64
	{
		switch configMinTotalHistoricalAmount.Currency {
		case internationalCurrency:
			amount = sumAmountFromIssuerAccounts(loadedData.TransactionHistory.IssuerAccounts, false)
		default:
			amount = sumAmountFromIssuerAccounts(loadedData.TransactionHistory.IssuerAccounts, true)
		}
	}

	resultEvaluatedData := &domain.ReasonResult{
		Actual:   amount,
		Accepted: minAmount,
	}

	if amount >= minAmount {
		evaluatedData.ApprovedData[domain.RuleMinTotalHistoricalAmount] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleMinTotalHistoricalAmount] = resultEvaluatedData
	addMetricRuleMinTotalHistoricalAmount(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertMinTransactionsQty(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configMinTransactionsQty := loadedData.RulesSite.MinTransactionsQty
	if configMinTransactionsQty == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleMinTransactionsQty)

	qty := 0
	for _, account := range loadedData.TransactionHistory.IssuerAccounts {
		qty += account.MatchTotal
		if qty >= configMinTransactionsQty.MinQty {
			evaluatedData.ApprovedData[domain.RuleMinTransactionsQty] = &domain.ReasonResult{
				Actual:   qty,
				Accepted: configMinTransactionsQty.MinQty,
			}
			return
		}
	}

	evaluatedData.UnapprovedData[domain.RuleMinTransactionsQty] = &domain.ReasonResult{
		Actual:   qty,
		Accepted: configMinTransactionsQty.MinQty,
	}

	addMetricRuleMinTransactionsQty(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}

func assertScore(
	ctx context.Context,
	evaluatedData *domain.EvaluatedData,
	input domain.ValidationInput,
	loadedData *loadedData,
	appliedRules ruleLog,
) {

	// Checks if this rule applies for the siteID.
	configScore := loadedData.RulesSite.Score
	if configScore == nil {
		return
	}
	setLogAppliedRules(appliedRules, input.PaymentID, input.SiteID, domain.RuleScore)

	resultEvaluatedData := &domain.ReasonResult{
		Actual:   *loadedData.CalculatedScore,
		Accepted: configScore.MinScoreValue,
	}

	if checkScore(*loadedData.CalculatedScore, configScore.MinScoreValue) {
		evaluatedData.ApprovedData[domain.RuleScore] = resultEvaluatedData
		return
	}

	evaluatedData.UnapprovedData[domain.RuleScore] = resultEvaluatedData

	addMetricRuleScore(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
}
