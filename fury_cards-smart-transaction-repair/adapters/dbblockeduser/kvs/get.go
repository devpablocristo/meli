package kvs

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbblockeduser"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	kvsclient "github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
)

func (b blockedUserRepository) Get(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
	var blockedUserDB dbblockeduser.BlockedUserDB

	err := b.repo.GetValue(ctx, kycIdentificationID, &blockedUserDB)
	if err != nil {
		if err == kvsclient.ErrNotFound {
			return nil, buildErrorNotFound(kycIdentificationID)
		}
		return nil, buildErrorBadGateway(err)
	}

	return dbblockeduser.ToBlockedUserOutput(blockedUserDB), nil
}
