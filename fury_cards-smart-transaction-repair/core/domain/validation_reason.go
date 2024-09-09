package domain

type RuleType string

const (
	RuleQtyReparationPerPeriodDays RuleType = "qty_reparation_per_period"
	RuleMaxAmountReparation        RuleType = "max_amount_reparation"
	RuleStatusDetailAllowed        RuleType = "status_detail"
	RuleRestrictions               RuleType = "restrictions"
	RuleBlockeduser                RuleType = "blockeduser"
	RuleQtyChargebackPerPeriodDays RuleType = "qty_chargeback_per_period"
	RuleSubTypeAllowed             RuleType = "subtype"
	RuleMerchantBlockedlist        RuleType = "merchant_blockedlist"
	RuleUserLifetime               RuleType = "user_lifetime"
	RuleMinTotalHistoricalAmount   RuleType = "min_total_historical_amount"
	RuleMinTransactionsQty         RuleType = "min_transactions_qty"
	RuleScore                      RuleType = "score"
)

type EvaluatedData struct {
	UnapprovedData Reason
	ApprovedData   ApprovedData
}

type Reason map[RuleType]*ReasonResult
type ApprovedData map[RuleType]*ReasonResult

type ReasonResult struct {
	Actual   interface{}
	Accepted interface{}
}

type ReasonResultPerPeriod struct {
	Qty        int
	PeriodDays int
}
