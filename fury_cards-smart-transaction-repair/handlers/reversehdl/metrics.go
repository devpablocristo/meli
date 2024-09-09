package reversehdl

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
)

func addMetricUnauthorized(ctx context.Context, siteID, clientApplication string) {
	metrics.Incr(ctx, domain.MetricReverse,
		metrics.Tags(
			"type", domain.MetricTagReverseFail,
			"action", "reverse.check_fury_request",
			"detail", "unauthorized manual request",
			"siteID", siteID,
			"requester_client_app", clientApplication,
		))
}

func addMetricIvalidPayload(ctx context.Context, clientApplication string) {
	metrics.Incr(ctx, domain.MetricReverse,
		metrics.Tags(
			"type", domain.MetricTagReverseFail,
			"action", "reverse.validate_request",
			"detail", "reverse invalid payload",
			"requester_client_app", clientApplication,
		))
}

func addMetricHeaderRequired(ctx context.Context, siteID, clientApplication string) {
	metrics.Incr(ctx, domain.MetricReverse,
		metrics.Tags(
			"type", domain.MetricTagReverseFail,
			"action", "reverse.validate_header",
			"detail", "header is required",
			"siteID", siteID,
			"requester_client_app", clientApplication,
		))
}

func addMetricTimerHandler(ctx context.Context, timing time.Duration, siteID, clientApplication string) {
	metrics.Timing(ctx, domain.MetricTimer, timing,
		metrics.Tags(
			"type", domain.MetricTagTimerReverseHandler,
			"siteID", siteID,
			"requester_client_app", clientApplication,
		),
	)
}
