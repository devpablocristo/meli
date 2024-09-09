package blockeduser

import (
	"context"
	"fmt"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (b blocked) UnblockUser(ctx context.Context, kycIdentificationID string) error {
	blockedUser, err := b.blockedUserRepo.Get(ctx, kycIdentificationID)
	if err != nil {
		if gnierrors.Type(err) == domain.NotFound {
			return buildErrorNotFound(err)
		}
		addMetricErr(ctx, gnierrors.Cause(err), fmt.Sprint("unblock_user.", actionBlockedlistGet), "", "")
		return buildErrorBadGateway(err)
	}

	now := new(time.Time)
	*now = time.Now().UTC()

	blockedUser.UnlockedAt = now

	err = b.blockedUserRepo.Save(ctx, kycIdentificationID, *blockedUser)
	if err != nil {
		addMetricErr(ctx, gnierrors.Cause(err), fmt.Sprint("unblock_user.", actionBlockedlistSave), blockedUser.SiteID, "")
		return err
	}

	addMetricBlockedUserUnblocked(ctx, blockedUser.SiteID)
	return nil
}
