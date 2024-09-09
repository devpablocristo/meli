package configurations

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const keyParametersValidation = "parameters-validation"

func (c *configService) LoadParametersValidation(ctx context.Context) (*domain.ConfigParametersValidation, error) {
	var parametersValidation parametersValidation

	err := c.serviceCache.GetJson(ctx, keyParametersValidation, &parametersValidation)
	if err != nil {
		return nil, treatError(err, keyParametersValidation)
	}

	return fillConfigParametersValidation(&parametersValidation), nil
}

func fillConfigParametersValidation(parametersValidation *parametersValidation) *domain.ConfigParametersValidation {
	return &domain.ConfigParametersValidation{
		RulesSite: fillConfigRulesSite(parametersValidation),
	}
}

func fillConfigRulesSite(parametersValidation *parametersValidation) map[string]*domain.ConfigRulesSite {
	rulesSite := map[string]*domain.ConfigRulesSite{}

	for site, rules := range parametersValidation.RulesSite {
		rulesSite[site] = &domain.ConfigRulesSite{
			StatusDetailAllowed:        rules.StatusDetailAllowed,
			QtyReparationPerPeriodDays: fillQtyReparationPerPeriodDays(rules.QtyReparationPerPeriodDays),
			MaxAmountReparation:        fillMaxAmountReparation(rules.MaxAmountReparation),
			EvaluationUser:             (*domain.ConfigEvaluationUser)(rules.EvaluationUser),
			QtyChargebackPerPeriodDays: (*domain.ConfigQtyChargebackPeriod)(rules.QtyChargebackPerPeriodDays),
			Blockeduser:                (*domain.ConfigBlockeduser)(rules.Blockeduser),
			MerchantBlockedlist:        rules.MerchantBlockedlist,
			SubTypeAllowed:             rules.SubTypeAllowed,
			UserLifetime:               fillUserLifetime(rules.UserLifetime),
			MinTotalHistoricalAmount:   (*domain.ConfigMinTotalHistoricalAmount)(rules.MinTotalHistoricalAmount),
			MinTransactionsQty:         (*domain.ConfigMinTransactionsQty)(rules.MinTransactionsQty),
			Score:                      (*domain.ConfigScore)(rules.Score),
		}
	}

	return rulesSite
}

func fillQtyReparationPerPeriodDays(qtyReparationPeriod *qtyReparationPeriod) *domain.ConfigQtyReparationPeriod {
	if qtyReparationPeriod == nil {
		return nil
	}

	ranges := []domain.ConfigQtyReparationPeriodTrxAgeRange{}
	for _, item := range qtyReparationPeriod.TrxAgeRanges {
		ranges = append(ranges, domain.ConfigQtyReparationPeriodTrxAgeRange{
			GteDays: item.GteDays,
			LteDays: item.LteDays,
			Qty:     item.Qty,
			Score:   (*domain.ConfigScoreQtyReparationPeriod)(item.Score),
		})
	}

	return &domain.ConfigQtyReparationPeriod{
		TrxAgeRanges: ranges,
		PeriodDays:   qtyReparationPeriod.PeriodDays,
		Qty:          qtyReparationPeriod.Qty,
		Score:        (*domain.ConfigScoreQtyReparationPeriod)(qtyReparationPeriod.Score),
	}
}

func fillMaxAmountReparation(maxAmount *maxAmountReparation) *domain.ConfigMaxAmountReparation {
	if maxAmount == nil {
		return nil
	}

	ranges := []domain.ConfigMaxAmountReparationTrxAgeRange{}
	for _, item := range maxAmount.TrxAgeRanges {
		ranges = append(ranges, domain.ConfigMaxAmountReparationTrxAgeRange{
			GteDays:              item.GteDays,
			LteDays:              item.LteDays,
			MaxAccumulatedAmount: item.MaxAccumulatedAmount,
			Currency:             item.Currency,
			DecimalDigits:        item.DecimalDigits,
			Score:                (*domain.ConfigScoreMaxAmountReparationTrxAgeRange)(item.ScoreMaxAmount),
		})
	}

	return &domain.ConfigMaxAmountReparation{
		TrxAgeRanges:  ranges,
		MaxAmount:     maxAmount.MaxAmount,
		Currency:      maxAmount.Currency,
		DecimalDigits: maxAmount.DecimalDigits,
		Score:         (*domain.ConfigScoreMaxAmountReparation)(maxAmount.Score),
	}
}

func fillUserLifetime(userLifetime *userLifetime) *domain.ConfigUserLifetime {
	if userLifetime == nil {
		return nil
	}

	return &domain.ConfigUserLifetime{
		GteDays: userLifetime.GteDays,
		Score:   (*domain.ConfigScoreUserLifetime)(userLifetime.Score),
	}
}
