package blockeduser

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
)

func addMetricTimerBlockedlistGet(ctx context.Context, siteID string, timing time.Duration) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerKvsBlockelistGet,
			"siteID", siteID,
		),
	)
}

func addMetricTimerBlockedlistSave(ctx context.Context, siteID string, timing time.Duration) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerKvsBlockelistSave,
			"siteID", siteID,
		),
	)
}

func addMetricTimerRepairGetList(ctx context.Context, siteID string, timing time.Duration) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerDsRepairGetList,
			"siteID", siteID,
		),
	)
}

func addMetricTagTimerBlockedUserGetListByAuthorizationIDs(ctx context.Context, siteID string, timing time.Duration) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerDsBlockedUserGetListByAuthorizationIDs,
			"siteID", siteID,
		),
	)
}
