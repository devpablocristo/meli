package validationforcepublishhdl

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type msgValidationPublish struct {
	PaymentID             string                            `json:"payment_id" validate:"required"`
	KycIdentitificationID string                            `json:"kyc_identification_id" validate:"required"`
	UserID                int64                             `json:"user_id" validate:"required"`
	IsEligible            bool                              `json:"is_eligible"`
	CreatedAt             time.Time                         `json:"created_at" validate:"required"`
	AuthorizationID       string                            `json:"authorization_id" validate:"required"`
	Type                  domain.TypeReparation             `json:"type_reparation" validate:"required"`
	SiteID                string                            `json:"site_id" validate:"required"`
	FaqID                 string                            `json:"faq_id,omitempty"`
	ClientApplication     string                            `json:"client_application"`
	ClientScope           string                            `json:"client_scope"`
	Amount                float64                           `json:"amount"`
	Currency              string                            `json:"currency"`
	Reason                map[domain.RuleType]publicDefault `json:"reason,omitempty"`
}

type publicDefault struct {
	Actual   interface{} `json:"actual,omitempty"`
	Accepted interface{} `json:"accepted,omitempty"`
}

func (r publicDefault) getValue() *domain.ReasonResult {
	switch accepted := r.Accepted.(type) {
	case map[string]interface{}:
		return &domain.ReasonResult{
			Actual: r.Actual,
			Accepted: domain.ReasonResultPerPeriod{
				Qty:        int(accepted["qty"].(float64)),
				PeriodDays: int(accepted["period_days"].(float64)),
			},
		}
	default:
		return &domain.ReasonResult{
			Actual:   r.Actual,
			Accepted: r.Accepted,
		}
	}
}
