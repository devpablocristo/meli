package blockeduser

import (
	"context"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (b blocked) CheckUserIsBlocked(
	ctx context.Context,
	input domain.BlockedUserInput,
	wrapFields *field.WrappedFields,
) (bool, error) {

	blockedUser, err := b.get(ctx, input.KycIdentificationID, input.SiteID, wrapFields)
	if err != nil {
		if gnierrors.Type(err) == domain.NotFound {
			return false, nil
		}
		addMetricErrFromReverseSrv(ctx, gnierrors.Cause(err), fmt.Sprint("check_user.", actionBlockedlistGet), input.SiteID)
		return false, buildErrorBadGateway(err)
	}

	if blockedUser.UnlockedAt != nil {
		return false, nil
	}

	return true, nil
}
