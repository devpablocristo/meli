package domain

import "time"

type SearchKycInput struct {
	UserID string
}

type SearchKycOutput struct {
	KycIdentificationID string
	DateCreated         time.Time
}
