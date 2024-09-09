package blockeduser

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

func (b blocked) save(
	ctx context.Context,
	blockedUser storage.BlockedUser,
	personIdentification,
	siteID string,
	wrapFields *field.WrappedFields,
) error {
	startTimer := wrapFields.Timers.Start(domain.TimerKvsBlockelistSave)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerBlockedlistSave(ctx, siteID, time.Since(start))
	}()

	return b.blockedUserRepo.Save(ctx, personIdentification, blockedUser)
}
