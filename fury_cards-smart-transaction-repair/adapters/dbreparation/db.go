package dbreparation

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type TransactionRepairDB struct {
	AuthorizationID     string                `json:"authorization_id"`
	KycIdentificationID string                `json:"kyc_identification_id"`
	UserID              int64                 `json:"user_id"`
	PaymentID           string                `json:"payment_id"`
	TransactionRepairID string                `json:"transaction_repair_id"`
	SiteID              string                `json:"site_id"`
	Type                domain.TypeReparation `json:"type"`
	CreatedAt           time.Time             `json:"created_at"`
	FaqID               string                `json:"faq_id,omitempty"`
	Requester           RequesterDB           `json:"requester"`
	MonetaryTransaction MonetaryTransactionDB `json:"monetary_transaction"`
}

type RequesterDB struct {
	ClientApplication string `json:"client_application"`
	ClientScope       string `json:"client_scope"`
}

type MonetaryTransactionDB struct {
	Billing    MonetaryValueDB `json:"billing"`
	Settlement MonetaryValueDB `json:"settlement"`
}

type MonetaryValueDB struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
