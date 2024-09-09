package blockeduser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (b blocked) UnblockUserByReversalClearingInBulk(
	ctx context.Context,
	input domain.TransactionsNewsFeedBulkInput,
	wrapFields *field.WrappedFields,
) (*domain.TransactionsBulkOutput, error) {

	output := &domain.TransactionsBulkOutput{AuthorizationIDsWithErr: map[string]error{}}

	authorizationsIDs := fillAuthorizationIdsAndIncrMetricCheckBlockedUsers(ctx, input)

	blockedUsers, err := b.getBlockedUsers(ctx, input.SiteID, authorizationsIDs, wrapFields)
	if err != nil {
		addMetricErr(
			ctx,
			err,
			fmt.Sprint(actionUnblokedUserByReversalClearing, ".blockeduser_search_repo.get"),
			input.SiteID,
			domain.MetricTagSourceQReversalTransactions,
		)
		return nil, buildErrorBadGateway(err)
	}

	kycIDsUnblocked := []string{}
	for _, blockedUser := range blockedUsers {
		now := time.Now().UTC()
		unlock := true
		hasAddReversalClearing := false
		for i, cap := range blockedUser.CapturedRepairs {
			if transactionID, found := input.AuthorizationsNewsFeed[cap.AuthorizationID]; found && cap.ReversalClearing == nil {
				blockedUser.CapturedRepairs[i].ReversalClearing = &storage.ReversalClearing{
					ReverseID: string(transactionID),
					CreatedAt: now,
				}
				hasAddReversalClearing = true
			}
			if unlock && blockedUser.CapturedRepairs[i].ReversalClearing == nil {
				unlock = false
			}
		}

		if hasAddReversalClearing {
			if unlock {
				blockedUser.UnlockedAt = &now
				kycIDsUnblocked = append(kycIDsUnblocked, blockedUser.KycIdentificationID)
			}
			b.processUnblocked(ctx, blockedUser, input, output, wrapFields)
		}
	}

	if len(kycIDsUnblocked) > 0 {
		wrapFields.Fields.Add(
			fmt.Sprintf(`%s: %s`, domain.KeyFieldKycIDs, strings.Join(kycIDsUnblocked, ", ")),
			fmt.Sprintf(`unblocked total_%v`, len(kycIDsUnblocked)))
	}

	return output, nil
}

func (b blocked) getBlockedUsers(
	ctx context.Context,
	siteID string,
	authorizationIDs storage.AuthorizationIDsSearch,
	wrapFields *field.WrappedFields,
) ([]storage.BlockedUser, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerDsBlockedUserGetListByAuthorizationIDs)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTagTimerBlockedUserGetListByAuthorizationIDs(ctx, siteID, time.Since(start))
	}()

	return b.blockedUserSearchRepo.GetListByAuthorizationIDs(siteID, authorizationIDs)
}

func fillAuthorizationIdsAndIncrMetricCheckBlockedUsers(
	ctx context.Context,
	input domain.TransactionsNewsFeedBulkInput,
) storage.AuthorizationIDsSearch {
	addMetricBlockedUserCheckBlockedUserByReversalInBulk(ctx, input.SiteID)

	authorizationsIDs := storage.AuthorizationIDsSearch{}
	for k := range input.AuthorizationsNewsFeed {
		authorizationsIDs = append(authorizationsIDs, k)
		addMetricBlockedUserCheckBlockedUserByReversal(ctx, input.SiteID)
	}

	return authorizationsIDs
}

func (b blocked) processUnblocked(
	ctx context.Context,
	blockedUser storage.BlockedUser,
	input domain.TransactionsNewsFeedBulkInput,
	output *domain.TransactionsBulkOutput,
	wrapFields *field.WrappedFields,
) {
	err := b.save(ctx, blockedUser, blockedUser.KycIdentificationID, blockedUser.SiteID, wrapFields)
	if err != nil {
		addMetricErr(
			ctx,
			gnierrors.Cause(err),
			fmt.Sprint(actionUnblokedUserByReversalClearing, ".", actionBlockedlistSave),
			blockedUser.SiteID,
			domain.MetricTagSourceQReversalTransactions,
		)

		for _, cap := range blockedUser.CapturedRepairs {
			if _, found := input.AuthorizationsNewsFeed[cap.AuthorizationID]; found {
				output.AuthorizationIDsWithErr[cap.AuthorizationID] = err
			}
		}
		return
	}

	// Metrics only when user is unlocked.
	if blockedUser.UnlockedAt != nil {
		addMetricBlockedUserUnblockedByReversalClearing(ctx, blockedUser.SiteID)
	}
}
