package reparation

import (
	"context"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

const (
	_daysToRepairRequest = 7
)

func addMetricErr(ctx context.Context, err error, action, siteID, faqID, clientApplication string) {
	incrMetricException(ctx, err, domain.MetricTagReverseFail, action, siteID, faqID, clientApplication)
}

func addMetricWarn(ctx context.Context, err error, action, siteID, faqID, clientApplication string) {
	incrMetricException(ctx, err, domain.MetricTagReverseWarning, action, siteID, faqID, clientApplication)
}

func incrMetricException(ctx context.Context, err error, typeTagReverse, action, siteID, faqID, clientApplication string) {
	errCause := gnierrors.Cause(err)
	if errCause != nil {
		err = errCause
	}

	tags := metrics.Tags(
		"type", typeTagReverse,
		"action", fmt.Sprint("reparation.srv.", action),
		"detail", err.Error(),
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricReverse, tags)
}

func addMetricReverseAlready(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagReverseAlready,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricReverse, tags)
}

func addMetricReverseNotEligible(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagReverseNotEligible,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricReverse, tags)
}

func addMetricReverseEligible(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagReverseEligible,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricReverse, tags)
}

func addMetricReverseSuccess(
	ctx context.Context,
	siteID string,
	daysToRepairRequest int,
	rangeOfDaysRepairRequest,
	faqID,
	clientApplication,
	commerceName,
	acquirerCode string,
) {
	tags := metrics.Tags(
		"type", domain.MetricTagReverseSuccess,
		"siteID", siteID,
		"days_to_repair_request", daysToRepairRequest,
		"range_of_days_repair_request", rangeOfDaysRepairRequest,
		"requester_client_app", clientApplication,
		"commerce_name", commerceName,
		"acquirer_code", acquirerCode,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricReverse, tags)
}

func addMetricReverseUnable(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagReverseUnable,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricReverse, tags)
}

func addMetricUnablePersonIDEmpty(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagUnablePersonIDEmpty,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricUnable, tags)

	addMetricReverseUnable(ctx, siteID, faqID, clientApplication)
}

func addMetricUnableTypeDiffAuthorization(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagUnableTypeDiffAuthorization,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricUnable, tags)

	addMetricReverseUnable(ctx, siteID, faqID, clientApplication)
}

func addMetricUnableUserNotAvailable(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagUnableUserNotAvailable,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricUnable, tags)

	addMetricReverseUnable(ctx, siteID, faqID, clientApplication)
}

func addMetricAuthorizationAlreadyReversed(ctx context.Context, siteID, faqID, clientApplication string) {
	tags := metrics.Tags(
		"type", domain.MetricTagAuthorizationAlreadyReversed,
		"siteID", siteID,
		"requester_client_app", clientApplication,
	)

	tags = appendFaqIDIfNotEmpty(tags, faqID)

	metrics.Incr(ctx, domain.MetricUnable, tags)

	addMetricReverseUnable(ctx, siteID, faqID, clientApplication)
}

func appendFaqIDIfNotEmpty(tags []string, faqID string) []string {
	if len(faqID) > 0 {
		tags = append(tags, fmt.Sprintf(`faqID:%s`, faqID))
	}

	return tags
}
