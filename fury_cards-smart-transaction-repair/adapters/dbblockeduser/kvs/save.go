package kvs

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbblockeduser"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

func (b blockedUserRepository) Save(ctx context.Context, kycIdentificationID string, blockedUser storage.BlockedUser) error {
	err := b.repo.Put(ctx, kycIdentificationID, fillBlockedUserDB(blockedUser))
	if err != nil {
		return buildErrorBadGateway(err)
	}

	return nil
}

func fillBlockedUserDB(blockedUser storage.BlockedUser) dbblockeduser.BlockedUserDB {
	return dbblockeduser.BlockedUserDB{
		KycIdentificationID: blockedUser.KycIdentificationID,
		UnlockedAt:          blockedUser.UnlockedAt,
		SiteID:              blockedUser.SiteID,
		CreatedAt:           blockedUser.CreatedAt,
		LastBlockedAt:       blockedUser.LastBlockedAt,
		CountCaptures:       blockedUser.CountCaptures,
		CapturedRepairs:     toCapturedRepairsDB(blockedUser.CapturedRepairs),
	}
}

func toCapturedRepairsDB(capturedRepairs []storage.CapturedRepairs) []dbblockeduser.CapturedRepairsDB {
	captRepairsDB := []dbblockeduser.CapturedRepairsDB{}

	for _, item := range capturedRepairs {
		captRepairsDB = append(captRepairsDB, dbblockeduser.CapturedRepairsDB{
			AuthorizationID:  item.AuthorizationID,
			UserID:           item.UserID,
			PaymentID:        item.PaymentID,
			TypeRepair:       item.TypeRepair,
			RepairedAt:       item.RepairedAt,
			CreatedAt:        item.CreatedAt,
			ReversalClearing: (*dbblockeduser.ReversalClearing)(item.ReversalClearing),
		})
	}

	return captRepairsDB
}
