package transactions

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
)

func addMetricIvalidPayload(ctx context.Context, siteID, source string) {
	metrics.Incr(
		ctx,
		domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserFail,
			"action", "consumer.decode_JSON",
			"detail", "consumer invalid payload",
			"source", source,
			"siteID", siteID,
		))
}

func addMetricTimerConsumerCaptures(ctx context.Context, siteID string, timing time.Duration) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerConsumerCaptures,
			"siteID", siteID,
		),
	)
}

func addMetricTimerConsumerReverses(ctx context.Context, siteID string, timing time.Duration) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerConsumerReverses,
			"siteID", siteID,
		),
	)
}
