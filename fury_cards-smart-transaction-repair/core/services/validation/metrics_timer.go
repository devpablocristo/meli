package validation

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
)

func addMetricTimerLoadParametersValidation(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerConfigurationsLoadParametersValidation,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerGetListByUserIDAndCreationPeriod(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerDsRepairGetListByUserIDAndCreationPeriod,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerValidationSave(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerDsValidationSave,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerPolicyAgentEvaluate(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerPolicyAgentEvaluate,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerSearchHubCountChargebackFromUser(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerSearchHubCountChargebackFromUser,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerSearchHubGetTransactionHistory(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerSearchHubGetTransactionHistory,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerGetScore(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerGetScore,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}

func addMetricTimerEventValidationResultPublish(ctx context.Context, siteID, faqID, clientApplication string, timing time.Duration) {
	tags := metrics.Tags(
		"type", domain.MetricTagTimerEventValidationResultPublish,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Timing(ctx, domain.MetricTimer, timing, tags)
}
