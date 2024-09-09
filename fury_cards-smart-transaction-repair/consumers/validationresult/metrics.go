package validationresult

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
)

func addMetricIvalidPayload(ctx context.Context) {
	metrics.Incr(
		ctx,
		domain.MetricEventValidationResult,
		metrics.Tags(
			"type", domain.MetricTagEventFail,
			"action", "consumer_decode_JSON",
			"detail", "consumer_invalid_payload",
		))
}

func addMetricTimerConsumerValidationResult(ctx context.Context, siteID string, timing time.Duration) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerConsumerValidationResult,
			"siteID", siteID,
		),
	)
}

func addMetricValidationResultReceived(ctx context.Context, siteID string) {
	metrics.Incr(
		ctx,
		domain.MetricEventValidationResult,
		metrics.Tags(
			"type", domain.MetricTagEventReceived,
			"siteID", siteID,
		))
}

func addMetricValidationResultDone(ctx context.Context, siteID string) {
	metrics.Incr(
		ctx,
		domain.MetricEventValidationResult,
		metrics.Tags(
			"type", domain.MetricTagEventDone,
			"siteID", siteID,
		))
}
