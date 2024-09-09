package validation

import (
	"context"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func addMetricErr(ctx context.Context, err error, action, siteID, faqID, clientApplication string) {
	incrMetricException(ctx, err, domain.MetricTagReverseFail, action, siteID, faqID, clientApplication)
}

func addMetricWarn(ctx context.Context, err error, action, siteID, faqID, clientApplication string) {
	incrMetricException(ctx, err, domain.MetricTagReverseWarning, action, siteID, faqID, clientApplication)
}

func incrMetricException(ctx context.Context, err error, typeTagReverseExp, action, siteID, faqID, clientApplication string) {
	errCause := gnierrors.Cause(err)
	if errCause != nil {
		err = errCause
	}

	tags := metrics.Tags(
		"type", typeTagReverseExp,
		"action", fmt.Sprint("validation.srv.", action),
		"detail", err.Error(),
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricReverse, tags)
}

func addMetricRule(ctx context.Context, metricTag domain.MetricTagRule, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", string(metricTag),
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricRule, tags)
}

func addMetricRuleMaxAmount(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagRuleMaxAmount, siteID, faqID, clientApplication)
}

func addMetricRuleMaxQtyRepairPerPeriod(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagRuleMaxQtyRepairPerPeriod, siteID, faqID, clientApplication)
}

func addMetricRuleUserIsBlocked(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagRuleUserIsBlocked, siteID, faqID, clientApplication)
}

func addMetricRuleStatusNotAllowed(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagRuleStatusNotAllowed, siteID, faqID, clientApplication)
}

func addMetricRuleUserWithRestriction(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagRuleUserWithRestriction, siteID, faqID, clientApplication)
}

func addMetricRuleMaxQtyChargebackPerPeriod(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagRuleMaxQtyChargebackPerPeriod, siteID, faqID, clientApplication)
}

func addMetricRuleSubTypeNotAllowed(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagRuleSubTypeNotAllowed, siteID, faqID, clientApplication)
}

func addMetricRuleMerchantBlockedlist(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagMerchantBlockedlist, siteID, faqID, clientApplication)
}

func addMetricRuleUserLifetime(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagUserLifetime, siteID, faqID, clientApplication)
}

func addMetricRuleMinTotalHistoricalAmount(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagMinTotalHistoricalAmount, siteID, faqID, clientApplication)
}

func addMetricRuleMinTransactionsQty(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagMinTransactionsQty, siteID, faqID, clientApplication)
}

func addMetricRuleScore(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetricRule(ctx, domain.MetricTagScore, siteID, faqID, clientApplication)
}

func addMetricInvalidPrerequisitesForRuleValidation(ctx context.Context, siteID, faqID, clientApplication string) {
	addMetric(ctx, domain.MetricUnable, domain.MetricTagInvalidPrerequisitesForRuleValidation, siteID, faqID, clientApplication)
	addMetric(ctx, domain.MetricReverse, domain.MetricTagReverseUnable, siteID, faqID, clientApplication)
}

func addMetric(ctx context.Context, metricName, typeTag, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", typeTag,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)
	tags = appendFaqIDIfNotEmpty(tags, faqID)
	metrics.Incr(ctx, metricName, tags)
}

func appendFaqIDIfNotEmpty(tags []string, faqID string) []string {
	if len(faqID) > 0 {
		tags = append(tags, fmt.Sprintf(`faqID:%s`, faqID))
	}

	return tags
}
