package dbblockeduser

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type BlockedUserDB struct {
	KycIdentificationID string              `json:"kyc_identification_id"`
	UnlockedAt          *time.Time          `json:"unlocked_at"`
	SiteID              string              `json:"site_id"`
	CreatedAt           time.Time           `json:"created_at"`
	LastBlockedAt       time.Time           `json:"last_blocked_at"`
	CountCaptures       int                 `json:"count_captures"`
	CapturedRepairs     []CapturedRepairsDB `json:"captured_repairs"`
}

type CapturedRepairsDB struct {
	AuthorizationID  string                `json:"authorization_id"`
	UserID           int64                 `json:"user_id"`
	PaymentID        string                `json:"payment_id"`
	TypeRepair       domain.TypeReparation `json:"type_repair"`
	RepairedAt       time.Time             `json:"repaired_at"`
	CreatedAt        time.Time             `json:"created_at"`
	ReversalClearing *ReversalClearing     `json:"reversal_clearing,omitempty"`
}

type ReversalClearing struct {
	ReverseID string    `json:"reverse_id"`
	CreatedAt time.Time `json:"created_at"`
}
