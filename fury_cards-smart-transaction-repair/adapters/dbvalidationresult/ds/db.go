package ds

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type validationResultDB struct {
	PaymentID            string                              `json:"payment_id"`
	KycIdentificationID  string                              `json:"kyc_identification_id"`
	UserID               int64                               `json:"user_id"`
	IsEligible           bool                                `json:"is_eligible"`
	CreatedAt            time.Time                           `json:"created_at"`
	AuthorizationID      string                              `json:"authorization_id"`
	Type                 domain.TypeReparation               `json:"type_reparation"`
	SiteID               string                              `json:"site_id"`
	FaqID                string                              `json:"faq_id,omitempty"`
	Requester            requesterDB                         `json:"requester"`
	TransactionEvaluated transactionEvaluatedDB              `json:"transaction_evaluated"`
	Score                *float64                            `json:"score,omitempty"`
	ApprovedData         map[domain.RuleType]reasonDefaultDB `json:"approved_data,omitempty"`
	Reason               map[domain.RuleType]reasonDefaultDB `json:"reason,omitempty"`
}

type reasonDefaultDB struct {
	Actual   interface{} `json:"actual,omitempty"`
	Accepted interface{} `json:"accepted,omitempty"`
}

type reasonResultPerPeriodDB struct {
	Qty        int `json:"qty"`
	PeriodDays int `json:"period_days"`
}

type requesterDB struct {
	ClientApplication string `json:"client_application"`
	ClientScope       string `json:"client_scope"`
}

type transactionEvaluatedDB struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
