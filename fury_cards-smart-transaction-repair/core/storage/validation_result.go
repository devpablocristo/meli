package storage

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type ValidationResult struct {
	PaymentID            string
	KycIdentificationID  string
	UserID               int64
	IsEligible           bool
	DateCreated          time.Time
	AuthorizationID      string
	Type                 domain.TypeReparation
	SiteID               string
	FaqID                string
	Requester            domain.Requester
	TransactionEvaluated TransactionEvaluated
	Score                *float64
	ApprovedData         domain.ApprovedData
	Reason               domain.Reason
}

type TransactionEvaluated struct {
	Amount   float64
	Currency string
}
