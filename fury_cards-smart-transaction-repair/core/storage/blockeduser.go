package storage

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type BlockedUser struct {
	KycIdentificationID string
	UnlockedAt          *time.Time
	SiteID              string
	CreatedAt           time.Time
	LastBlockedAt       time.Time
	CountCaptures       int
	CapturedRepairs     []CapturedRepairs
}

type CapturedRepairs struct {
	AuthorizationID  string
	UserID           int64
	PaymentID        string
	TypeRepair       domain.TypeReparation
	RepairedAt       time.Time
	CreatedAt        time.Time
	ReversalClearing *ReversalClearing
}

type ReversalClearing struct {
	ReverseID string
	CreatedAt time.Time
}

type AuthorizationIDsSearch []interface{}
