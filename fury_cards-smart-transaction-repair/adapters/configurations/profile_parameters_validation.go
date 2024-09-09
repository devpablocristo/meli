package configurations

type parametersValidation struct {
	RulesSite map[string]rulesSite `json:"rules_site"`
}

type rulesSite struct {
	StatusDetailAllowed        map[string]struct{}       `json:"status_detail_allowed"`
	QtyReparationPerPeriodDays *qtyReparationPeriod      `json:"qty_reparation_per_period_days"`
	MaxAmountReparation        *maxAmountReparation      `json:"max_amount_reparation"`
	EvaluationUser             *evaluationUser           `json:"evaluation_user"`
	QtyChargebackPerPeriodDays *qtyChargebackPeriod      `json:"qty_chargeback_per_period_days"`
	Blockeduser                *blockeduser              `json:"blockeduser"`
	MerchantBlockedlist        []string                  `json:"merchant_blockedlist"`
	SubTypeAllowed             map[string]struct{}       `json:"subtype_allowed"`
	UserLifetime               *userLifetime             `json:"user_lifetime"`
	MinTotalHistoricalAmount   *minTotalHistoricalAmount `json:"min_total_historical_amount"`
	MinTransactionsQty         *minTransactionsQty       `json:"min_transactions_qty"`
	Score                      *score                    `json:"score"`
}

type qtyReparationPeriod struct {
	TrxAgeRanges []qtyReparationPeriodTrxAgeRange `json:"trx_age_range"`
	Qty          int                              `json:"qty"`
	PeriodDays   int                              `json:"period_days"`
	Score        *scoreQtyReparationPeriod        `json:"score"`
}

type qtyReparationPeriodTrxAgeRange struct {
	GteDays int                       `json:"gte_days"`
	LteDays *int                      `json:"lte_days"`
	Qty     int                       `json:"qty"`
	Score   *scoreQtyReparationPeriod `json:"score"`
}

type scoreQtyReparationPeriod struct {
	MinScoreValue float64 `json:"min_score_value"`
	Qty           int     `json:"qty"`
}

type maxAmountReparation struct {
	TrxAgeRanges  []maxAmountReparationTrxAgeRange `json:"trx_age_range"`
	MaxAmount     float64                          `json:"max_amount"`
	Currency      string                           `json:"currency"`
	DecimalDigits int                              `json:"decimal_digits"`
	Score         *scoreMaxAmountReparation        `json:"score"`
}

type maxAmountReparationTrxAgeRange struct {
	GteDays              int                                  `json:"gte_days"`
	LteDays              *int                                 `json:"lte_days"`
	MaxAccumulatedAmount float64                              `json:"max_accumulated_amount"`
	Currency             string                               `json:"currency"`
	DecimalDigits        int                                  `json:"decimal_digits"`
	ScoreMaxAmount       *scoreMaxAmountReparationTrxAgeRange `json:"score"`
}

type scoreMaxAmountReparation struct {
	MinScoreValue        float64 `json:"min_score_value"`
	MaxAmountDecimalized float64 `json:"max_amount"`
}

type scoreMaxAmountReparationTrxAgeRange struct {
	MinScoreValue                   float64 `json:"min_score_value"`
	MaxAccumulatedAmountDecimalized float64 `json:"max_accumulated_amount"`
}

type evaluationUser struct {
	Tags []string `json:"tags"`
}

type qtyChargebackPeriod struct {
	Qty        int `json:"qty"`
	PeriodDays int `json:"period_days"`
}

type userLifetime struct {
	GteDays int                `json:"gte_days"`
	Score   *scoreUserLifetime `json:"score"`
}

type scoreUserLifetime struct {
	GteDays       int     `json:"gte_days"`
	MinScoreValue float64 `json:"min_score_value"`
}

type minTotalHistoricalAmount struct {
	MinAmount     float64 `json:"min_amount"`
	Currency      string  `json:"currency"`
	DecimalDigits int     `json:"decimal_digits"`
}

type minTransactionsQty struct {
	MinQty int `json:"min_qty"`
}

type blockeduser struct{}

type score struct {
	MinScoreValue float64 `json:"min_score_value"`
}
