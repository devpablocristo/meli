package validation

import (
	"fmt"
	"math"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type ruleLog map[string][]domain.RuleType

func sumAmountFromIssuerAccounts(issuerAccounts []domain.HistoryIssuerAccounts, currencyLocal bool) (amount float64) {
	for _, account := range issuerAccounts {
		for _, transaction := range account.Transactions {
			if currencyLocal {
				amount += transaction.Billing.Amount / math.Pow10(transaction.Billing.DecimalDigits)
				continue
			}
			amount += transaction.Settlement.Amount / math.Pow10(transaction.Settlement.DecimalDigits)
		}
	}
	return
}

func setLogAppliedRules(appliedRules ruleLog, paymentID, siteID string, rule domain.RuleType) {
	key := fmt.Sprintf(`%s|%s`, paymentID, siteID)
	appliedRules[key] = append(appliedRules[key], domain.RuleType(fmt.Sprint(rule, ";")))
}

func round(value float64, precision int) float64 {
	factor := math.Pow10(precision)
	round := math.Round(value * factor)
	return round / math.Pow10(precision)
}

func checkScore(calculatedScore, ruleScore float64) bool {
	return calculatedScore >= ruleScore
}
