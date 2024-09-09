package domain

type ConfigParametersValidation struct {
	RulesSite map[string]*ConfigRulesSite
}

type ConfigRulesSite struct {
	StatusDetailAllowed        map[string]struct{}
	QtyReparationPerPeriodDays *ConfigQtyReparationPeriod
	MaxAmountReparation        *ConfigMaxAmountReparation
	EvaluationUser             *ConfigEvaluationUser
	QtyChargebackPerPeriodDays *ConfigQtyChargebackPeriod
	Blockeduser                *ConfigBlockeduser
	MerchantBlockedlist        []string
	SubTypeAllowed             map[string]struct{}
	UserLifetime               *ConfigUserLifetime
	MinTotalHistoricalAmount   *ConfigMinTotalHistoricalAmount
	MinTransactionsQty         *ConfigMinTransactionsQty
	Score                      *ConfigScore
}

type ConfigQtyReparationPeriod struct {
	TrxAgeRanges []ConfigQtyReparationPeriodTrxAgeRange
	Qty          int
	PeriodDays   int
	Score        *ConfigScoreQtyReparationPeriod
}

type ConfigScoreQtyReparationPeriod struct {
	MinScoreValue float64
	Qty           int
}

type ConfigQtyReparationPeriodTrxAgeRange struct {
	GteDays int
	LteDays *int
	Qty     int
	Score   *ConfigScoreQtyReparationPeriod
}

type ConfigMaxAmountReparation struct {
	TrxAgeRanges  []ConfigMaxAmountReparationTrxAgeRange
	MaxAmount     float64
	Currency      string
	DecimalDigits int
	Score         *ConfigScoreMaxAmountReparation
}

type ConfigMaxAmountReparationTrxAgeRange struct {
	GteDays              int
	LteDays              *int
	MaxAccumulatedAmount float64
	Currency             string
	DecimalDigits        int
	Score                *ConfigScoreMaxAmountReparationTrxAgeRange
}

type ConfigEvaluationUser struct {
	Tags []string
}

type ConfigQtyChargebackPeriod struct {
	Qty        int
	PeriodDays int
}

type ConfigUserLifetime struct {
	GteDays int
	Score   *ConfigScoreUserLifetime
}

type ConfigScoreUserLifetime struct {
	GteDays       int
	MinScoreValue float64
}

type ConfigMinTotalHistoricalAmount struct {
	MinAmount     float64
	Currency      string
	DecimalDigits int
}

type ConfigMinTransactionsQty struct {
	MinQty int
}

type ConfigBlockeduser struct{}

type ConfigScore struct {
	MinScoreValue float64
}

type ConfigScoreMaxAmountReparation struct {
	MinScoreValue        float64
	MaxAmountDecimalized float64
}

type ConfigScoreMaxAmountReparationTrxAgeRange struct {
	MinScoreValue                   float64
	MaxAccumulatedAmountDecimalized float64
}

func (c ConfigRulesSite) CheckApplyAnyScore() bool {
	return c.Score != nil ||
		(c.MaxAmountReparation != nil && c.MaxAmountReparation.Score != nil) ||
		(c.QtyReparationPerPeriodDays != nil && c.QtyReparationPerPeriodDays.Score != nil) ||
		(c.UserLifetime != nil && c.UserLifetime.Score != nil)
}
