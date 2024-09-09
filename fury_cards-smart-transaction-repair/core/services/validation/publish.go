package validation

import (
	"context"
	"fmt"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

func (v validationService) PublishEvent(
	ctx context.Context,
	validationResult *storage.ValidationResult,
	wrapFields *field.WrappedFields,
) {

	startTimer := wrapFields.Timers.Start(domain.TimerEventValidationResultPublish)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerEventValidationResultPublish(
			ctx, validationResult.SiteID, validationResult.FaqID, validationResult.Requester.ClientApplication, time.Since(start))
	}()

	err := v.eventValidationResult.Publish(ctx, *validationResult, wrapFields)
	if err != nil {
		addMetricWarn(
			ctx,
			err,
			"unpublished_event_validation_result",
			validationResult.SiteID,
			validationResult.FaqID,
			validationResult.Requester.ClientApplication,
		)
		v.logErrorAndNotStop(
			ctx,
			fmt.Sprintf("unpublished_event_validation_result_by_paymentID: %s", validationResult.PaymentID),
			err,
			wrapFields,
		)
	}
}

func (v validationService) logErrorAndNotStop(ctx context.Context, msgLog string, err error, wrapFields *field.WrappedFields) {
	logFields := []log.Field{v.log.Err(err)}
	logFields = append(logFields, wrapFields.Fields.ToLogField(log.InfoLevel))
	logFields = append(logFields, wrapFields.Timers.ToLogField())

	v.log.Warnln(ctx, msgLog, logFields...)
}
