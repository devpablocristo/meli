package reparation

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

const (
	countGoroutine = 3
)

type loadedData struct {
	TransactionRepair *storage.ReparationOutput
	TransactionData   *domain.TransactionData
	UserKycData       *domain.SearchKycOutput
}

func (r reparation) loadData(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	wrapFields *field.WrappedFields,
) (*loadedData, error) {

	var errTransaction, errSearchKyc error

	loadedData := &loadedData{}
	wg := sync.WaitGroup{}
	wg.Add(countGoroutine)

	go func() {
		defer wg.Done()
		loadedData.TransactionRepair = r.getTransactionRepair(ctx, input, wrapFields)
	}()

	go func() {
		defer wg.Done()
		loadedData.UserKycData, errSearchKyc = r.getUserKyc(ctx, input, wrapFields)
	}()

	go func() {
		defer wg.Done()
		loadedData.TransactionData, errTransaction = r.getTransaction(ctx, input, wrapFields)
	}()

	wg.Wait()

	if errTransaction != nil {
		return nil, errTransaction
	}

	if errSearchKyc != nil {
		return nil, errSearchKyc
	}

	return loadedData, nil
}

func (r reparation) getTransactionRepair(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	wrapFields *field.WrappedFields,
) *storage.ReparationOutput {

	startTimer := wrapFields.Timers.Start(domain.TimerDsRepairGetByPaymentIDAndUserID)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerGetByPaymentIDAndUserID(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	reparation, err := r.reparationSearchRepo.GetByPaymentIDAndUserID(input.PaymentID, input.UserID)
	if err != nil {
		if gnierrors.Type(err) != domain.NotFound {
			addMetricErr(ctx, err, "get_transaction_repair", input.SiteID, input.FaqID, input.Requester.ClientApplication)
			r.logErrorAndNotStop(ctx, fmt.Sprintf("unrecovered transaction repair, by paymentID: %s", input.PaymentID), err)
		}
	}
	return reparation
}

func (r reparation) getUserKyc(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	wrapFields *field.WrappedFields,
) (*domain.SearchKycOutput, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerSearchKycGetUserKyc)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerGetUserKyc(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	kycInput := domain.SearchKycInput{
		UserID: strconv.Itoa(int(input.UserID)),
	}

	kycOutput, err := r.searchKyc.GetUserKyc(ctx, kycInput)
	if err != nil {
		if gnierrors.Type(err) == domain.UnavailableForLegalReasons {
			addMetricUnableUserNotAvailable(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
			return nil, buildErrorUnableRequest(err, msgErrorRequest)
		}

		addMetricErr(ctx, err, "get_user_kyc", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, reBuildError(ctx, r.log, domain.BadGateway.String(), err)
	}

	return kycOutput, nil
}

func (r reparation) getTransaction(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	wrapFields *field.WrappedFields,
) (*domain.TransactionData, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerGraphqlGetPaymentTransactionByPaymentID)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerGetPaymentTransactionByPaymentID(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	transactionInput := domain.TransactionInput{
		PaymentID: input.PaymentID,
	}

	transactionOutput, err := r.searchHub.GetPaymentTransaction(ctx, transactionInput, wrapFields)
	if err != nil {
		addMetricErr(ctx, err, "get_transaction", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, reBuildError(ctx, r.log, domain.BadGateway.String(), err)
	}

	return transactionOutput.TransactionData, nil
}
