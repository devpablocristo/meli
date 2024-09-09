package topicvalidation

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type valueResult string

type validationResultTopic struct {
	ID                  string                            `json:"id"`
	PaymentID           string                            `json:"payment_id"`
	KycIdentificationID string                            `json:"kyc_identification_id"`
	UserID              int64                             `json:"user_id"`
	IsEligible          bool                              `json:"is_eligible"`
	CreatedAt           time.Time                         `json:"created_at"`
	AuthorizationID     string                            `json:"authorization_id"`
	Type                domain.TypeReparation             `json:"type_reparation"`
	SiteID              string                            `json:"site_id"`
	Rule                string                            `json:"rule"`
	FaqID               string                            `json:"faq_id,omitempty"`
	ValueDeclined       valueResult                       `json:"value_declined,omitempty"`
	AllowedValue        valueResult                       `json:"allowed_value,omitempty"`
	ClientApplication   string                            `json:"client_application"`
	ClientScope         string                            `json:"client_scope"`
	Amount              float64                           `json:"amount"`
	Currency            string                            `json:"currency"`
	Score               *float64                          `json:"score,omitempty"`
	ApprovedData        map[domain.RuleType]reasonDefault `json:"approved_data,omitempty"`
	Reason              map[domain.RuleType]reasonDefault `json:"reason,omitempty"`
}

type reasonDefault struct {
	Actual   interface{} `json:"actual,omitempty"`
	Accepted interface{} `json:"accepted,omitempty"`
}

type reasonResultPerPeriod struct {
	Qty        int `json:"qty"`
	PeriodDays int `json:"period_days"`
}

func (v *valueResult) toString(value interface{}) {
	b, _ := json.Marshal(&value)
	cast := string(b)
	cast = strings.ReplaceAll(cast, "\"", "")
	*v = valueResult(strings.ReplaceAll(cast, "null", ""))
}
