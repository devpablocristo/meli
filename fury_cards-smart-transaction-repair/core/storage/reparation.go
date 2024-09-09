package storage

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type BaseReparation struct {
	AuthorizationID     string
	KycIdentificationID string
	PaymentID           string
	UserID              int64
	TransactionRepairID string
	SiteID              string
	Type                domain.TypeReparation
	CreatedAt           time.Time
	FaqID               string
	MonetaryTransaction MonetaryTransaction
	Requester           domain.Requester
}

type ReparationInput struct {
	BaseReparation
}

type ReparationOutput struct {
	BaseReparation
}

type MonetaryTransaction struct {
	Billing    MonetaryValue
	Settlement MonetaryValue
}

type MonetaryValue struct {
	Amount   float64
	Currency string
}
