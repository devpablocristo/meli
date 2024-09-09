package blockeduser

import (
	"context"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

// Description tag for metric.
type action string

const (
	actionBlockedlistSave                action = "blockedlist_repo.save"
	actionBlockedlistGet                 action = "blockedlist_repo.get"
	actionUnblokedUserByReversalClearing action = "unblock_user_by_reverse_clearing_in_bulk"
)

func addMetricErr(ctx context.Context, err error, action, siteID, source string) {
	errCause := gnierrors.Cause(err)
	if errCause != nil {
		err = errCause
	}

	metrics.Incr(
		ctx, domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserFail,
			"source", source,
			"action", action,
			"detail", err.Error(),
			"siteID", siteID,
		))
}

func addMetricErrFromReverseSrv(ctx context.Context, errCause error, action, siteID string) {
	metrics.Incr(
		ctx, domain.MetricReverse,
		metrics.Tags(
			"type", domain.MetricTagReverseFail,
			"action", action,
			"detail", errCause.Error(),
			"siteID", siteID,
		))
}

func addMetricBlockedUserCheckRepairByCapture(ctx context.Context, siteID string) {
	metrics.Incr(ctx, domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserCheckRepairByCapture,
			"siteID", siteID,
		))
}

func addMetricBlockedUserCheckRepairByCaptureInBulk(ctx context.Context, siteID string) {
	metrics.Incr(ctx, domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserCheckRepairByCaptureInBulk,
			"siteID", siteID,
		))
}

func addMetricBlockedUserBlocked(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagBlockedUserBlocked,
		"siteID", siteID,
	)
	tags = appendIfNotEmpty(tags, faqID, clientApplication)
	metrics.Incr(ctx, domain.MetricBlockedUser, tags)
}

func addMetricBlockedUserBlockedAgain(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagBlockedUserBlockedAgain,
		"siteID", siteID,
	)
	tags = appendIfNotEmpty(tags, faqID, clientApplication)
	metrics.Incr(ctx, domain.MetricBlockedUser, tags)
}

func addMetricBlockedUserUnblocked(ctx context.Context, siteID string) {
	metrics.Incr(
		ctx,
		domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserUnblocked,
			"siteID", siteID,
		))
}

func addMetricBlockedUserUnblockedByReversalClearing(ctx context.Context, siteID string) {
	metrics.Incr(
		ctx,
		domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserUnblockedByReversalClearing,
			"siteID", siteID,
		))
}

func addMetricBlockedUserCheckBlockedUserByReversal(ctx context.Context, siteID string) {
	metrics.Incr(ctx, domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserCheckBlockedUserByReversal,
			"siteID", siteID,
		))
}

func addMetricBlockedUserCheckBlockedUserByReversalInBulk(ctx context.Context, siteID string) {
	metrics.Incr(ctx, domain.MetricBlockedUser,
		metrics.Tags(
			"type", domain.MetricTagBlockedUserCheckBlockedUserByReversalInBulk,
			"siteID", siteID,
		))
}

func appendIfNotEmpty(tags []string, faqID, clientApplication string) []string {
	if len(faqID) > 0 {
		tags = append(tags, fmt.Sprintf(`faqID:%s`, faqID))
	}

	if len(clientApplication) > 0 {
		tags = append(tags, fmt.Sprintf(`requester_client_app:%s`, clientApplication))
	}

	return tags
}
