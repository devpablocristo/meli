package reparation

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
)

func addMetricTimerGetByPaymentIDAndUserID(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerDsRepairGetByPaymentIDAndUserID,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerGetPaymentTransactionByPaymentID(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerGraphqlGetPaymentTransactionByPaymentID,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerCardsTransactionReverse(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerCardsTransactionReverse,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerRepairSave(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerKvsRepairSave,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerGetUserKyc(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerSearchKycGetUserKyc,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}
