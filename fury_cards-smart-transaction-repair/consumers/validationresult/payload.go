package validationresult

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type msgTopicRequest struct {
	Body msgValidationResult `json:"msg"`
}

type msgValidationResult struct {
	PaymentID           string                            `json:"payment_id" validate:"required"`
	KycIdentificationID string                            `json:"kyc_identification_id" validate:"required"`
	UserID              int64                             `json:"user_id" validate:"required"`
	IsEligible          bool                              `json:"is_eligible"`
	CreatedAt           time.Time                         `json:"created_at" validate:"required"`
	AuthorizationID     string                            `json:"authorization_id" validate:"required"`
	Type                domain.TypeReparation             `json:"type_reparation" validate:"required"`
	SiteID              string                            `json:"site_id" validate:"required"`
	FaqID               string                            `json:"faq_id,omitempty"`
	ClientApplication   string                            `json:"client_application"`
	ClientScope         string                            `json:"client_scope"`
	Amount              float64                           `json:"amount"`
	Currency            string                            `json:"currency"`
	Score               *float64                          `json:"score"`
	ApprovedData        map[domain.RuleType]reasonDefault `json:"approved_data,omitempty"`
	Reason              map[domain.RuleType]reasonDefault `json:"reason,omitempty"`
}

type reasonDefault struct {
	Actual   interface{} `json:"actual,omitempty"`
	Accepted interface{} `json:"accepted,omitempty"`
}

type simpleMessageResponse struct {
	Message string `json:"message"`
}

func (r reasonDefault) getValue() *domain.ReasonResult {
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
