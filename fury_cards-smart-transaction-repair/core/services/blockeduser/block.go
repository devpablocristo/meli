package blockeduser

import (
	"context"
	"fmt"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (b blocked) block(
	ctx context.Context,
	transactionRepair storage.ReparationOutput,
	wrapFields *field.WrappedFields,
) error {

	blockedUser, err := b.get(ctx, transactionRepair.KycIdentificationID, transactionRepair.SiteID, wrapFields)
	if err != nil {
		if gnierrors.Type(err) == domain.NotFound {
			return b.create(ctx, transactionRepair, wrapFields)
		}

		action := fmt.Sprint("block.", actionBlockedlistGet)
		wrapFields.Fields.Add(string(domain.KeyFieldAuthorizationID), fmt.Sprint(transactionRepair.AuthorizationID, ":", action))
		addMetricErr(ctx, err, action, transactionRepair.SiteID, domain.MetricTagSourceQCaptureTransactions)

		return buildErrorBadGateway(err)
	}

	for _, cr := range blockedUser.CapturedRepairs {
		if cr.AuthorizationID == transactionRepair.AuthorizationID {
			return nil
		}
	}

	return b.append(ctx, *blockedUser, transactionRepair, wrapFields)
}

func (b blocked) create(
	ctx context.Context,
	transactionRepair storage.ReparationOutput,
	wrapFields *field.WrappedFields,
) error {

	now := time.Now().UTC()

	blockedUser := storage.BlockedUser{
		KycIdentificationID: transactionRepair.KycIdentificationID,
		UnlockedAt:          nil,
		SiteID:              transactionRepair.SiteID,
		CreatedAt:           now,
		LastBlockedAt:       now,
		CountCaptures:       1,
		CapturedRepairs: []storage.CapturedRepairs{
			fillCapturedRepair(transactionRepair, now),
		},
	}

	err := b.save(ctx, blockedUser, transactionRepair.KycIdentificationID, transactionRepair.SiteID, wrapFields)
	if err != nil {
		action := fmt.Sprint("create.", actionBlockedlistSave)
		wrapFields.Fields.Add(string(domain.KeyFieldAuthorizationID), fmt.Sprint(transactionRepair.AuthorizationID, ":", action))

		addMetricErr(
			ctx,
			gnierrors.Cause(err),
			action,
			transactionRepair.SiteID,
			domain.MetricTagSourceQCaptureTransactions,
		)
		return err
	}

	wrapFields.Fields.Add(fmt.Sprint(domain.KeyFieldUserID, ":", transactionRepair.UserID), "block")
	addMetricBlockedUserBlocked(ctx, transactionRepair.SiteID, transactionRepair.FaqID, transactionRepair.Requester.ClientApplication)
	return nil
}

func (b blocked) append(
	ctx context.Context,
	blockedUser storage.BlockedUser,
	transactionRepair storage.ReparationOutput,
	wrapFields *field.WrappedFields,
) error {
	now := time.Now().UTC()

	blockedUser.UnlockedAt = nil
	blockedUser.LastBlockedAt = now
	blockedUser.CapturedRepairs = append(blockedUser.CapturedRepairs, fillCapturedRepair(transactionRepair, now))
	blockedUser.CountCaptures = len(blockedUser.CapturedRepairs)

	err := b.save(ctx, blockedUser, transactionRepair.KycIdentificationID, transactionRepair.SiteID, wrapFields)
	if err != nil {
		action := fmt.Sprint("append.", actionBlockedlistSave)
		wrapFields.Fields.Add(string(domain.KeyFieldAuthorizationID), fmt.Sprint(transactionRepair.AuthorizationID, ":", action))

		addMetricErr(
			ctx,
			gnierrors.Cause(err),
			action,
			transactionRepair.SiteID,
			domain.MetricTagSourceQCaptureTransactions,
		)
		return err
	}

	wrapFields.Fields.Add(fmt.Sprint(domain.KeyFieldUserID, ":", transactionRepair.UserID), "block_again")
	addMetricBlockedUserBlockedAgain(ctx, transactionRepair.SiteID, transactionRepair.FaqID, transactionRepair.Requester.ClientApplication)
	return nil
}

func fillCapturedRepair(transactionRepair storage.ReparationOutput, now time.Time) storage.CapturedRepairs {
	return storage.CapturedRepairs{
		AuthorizationID: transactionRepair.AuthorizationID,
		UserID:          transactionRepair.UserID,
		PaymentID:       transactionRepair.PaymentID,
		TypeRepair:      transactionRepair.Type,
		RepairedAt:      transactionRepair.CreatedAt,
		CreatedAt:       now,
	}
}
