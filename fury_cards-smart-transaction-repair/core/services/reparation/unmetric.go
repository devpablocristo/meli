package reparation

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
)

// Decrement when:
// - reversal is declined for known reason from transactions.
func decrementMetricReverseEligible(ctx context.Context, siteID, faqID string) {
	tags := metrics.Tags(
		"type", domain.MetricTagReverseEligible,
		"siteID", siteID,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Decr(ctx, domain.MetricReverse, tags)
}
