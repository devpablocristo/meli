package blockeduser

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (b blocked) BlockUserIfReparationExistsInBulk(
	ctx context.Context,
	input domain.TransactionsNewsFeedBulkInput,
	wrapFields *field.WrappedFields,
) (*domain.TransactionsBulkOutput, error) {

	output := &domain.TransactionsBulkOutput{AuthorizationIDsWithErr: map[string]error{}}

	authorizationsIDs := fillAuthorizationsAndIncrMetricCheckRepairs(ctx, input)

	transactionRepairs, err := b.getRepairs(ctx, input, wrapFields, authorizationsIDs)
	if err != nil {
		if gnierrors.Type(err) == domain.NotFound {
			return output, nil
		}
		action := "transaction_repair_search_repo.get_list_by_authorization_ids"
		addMetricErr(ctx, err, action, input.SiteID, domain.MetricTagSourceQCaptureTransactions)
		return nil, buildErrorBadGateway(err)
	}

	for _, repair := range transactionRepairs {
		err := b.block(ctx, repair, wrapFields)
		if err != nil {
			output.AuthorizationIDsWithErr[repair.AuthorizationID] = err
			continue
		}
	}

	return output, nil
}

func (b blocked) getRepairs(
	ctx context.Context,
	input domain.TransactionsNewsFeedBulkInput,
	wrapFields *field.WrappedFields,
	authorizationIDs []string,
) ([]storage.ReparationOutput, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerDsRepairGetList)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerRepairGetList(ctx, input.SiteID, time.Since(start))
	}()

	return b.reparationSearchRepo.GetListByAuthorizationIDs(authorizationIDs...)
}

func fillAuthorizationsAndIncrMetricCheckRepairs(
	ctx context.Context,
	input domain.TransactionsNewsFeedBulkInput,
) []string {
	addMetricBlockedUserCheckRepairByCaptureInBulk(ctx, input.SiteID)

	authorizationsIDs := []string{}
	for k := range input.AuthorizationsNewsFeed {
		authorizationsIDs = append(authorizationsIDs, k)
		addMetricBlockedUserCheckRepairByCapture(ctx, input.SiteID)
	}

	return authorizationsIDs
}
