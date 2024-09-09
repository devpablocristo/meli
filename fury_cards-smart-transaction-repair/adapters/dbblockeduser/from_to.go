package dbblockeduser

import "github.com/melisource/cards-smart-transaction-repair/core/storage"

func ToBlockedUserOutput(blockedUserDB BlockedUserDB) *storage.BlockedUser {
	return &storage.BlockedUser{
		KycIdentificationID: blockedUserDB.KycIdentificationID,
		UnlockedAt:          blockedUserDB.UnlockedAt,
		SiteID:              blockedUserDB.SiteID,
		CreatedAt:           blockedUserDB.CreatedAt,
		LastBlockedAt:       blockedUserDB.LastBlockedAt,
		CountCaptures:       blockedUserDB.CountCaptures,
		CapturedRepairs:     toCapturedRepairs(blockedUserDB.CapturedRepairs),
	}
}

func toCapturedRepairs(capturedRepairs []CapturedRepairsDB) []storage.CapturedRepairs {
	captRepairs := []storage.CapturedRepairs{}

	for _, item := range capturedRepairs {
		captRepairs = append(captRepairs, storage.CapturedRepairs{
			AuthorizationID:  item.AuthorizationID,
			UserID:           item.UserID,
			PaymentID:        item.PaymentID,
			TypeRepair:       item.TypeRepair,
			RepairedAt:       item.RepairedAt,
			CreatedAt:        item.CreatedAt,
			ReversalClearing: (*storage.ReversalClearing)(item.ReversalClearing),
		})
	}

	return captRepairs
}
