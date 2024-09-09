package blockeduser

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

func (b blocked) get(
	ctx context.Context,
	kycIdentificationID,
	siteID string,
	wrapFields *field.WrappedFields,
) (*storage.BlockedUser, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerKvsBlockelistGet)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerBlockedlistGet(ctx, siteID, time.Since(start))
	}()

	return b.blockedUserRepo.Get(ctx, kycIdentificationID)
}
