package validation

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

func (v validationService) Save(
	ctx context.Context,
	inputResult storage.ValidationResult,
	wrapFields *field.WrappedFields,
) error {

	startTimer := wrapFields.Timers.Start(domain.TimerDsValidationSave)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerValidationSave(ctx, inputResult.SiteID, inputResult.FaqID, inputResult.Requester.ClientApplication, time.Since(start))
	}()

	return v.validationResultRepo.Save(ctx, inputResult.PaymentID, inputResult)
}
