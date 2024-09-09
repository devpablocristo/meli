package reparation

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (r reparation) Reverse(ctx context.Context, input domain.ReverseTransactionInput, wrapFields *field.WrappedFields) error {
	loadedData, err := r.loadData(ctx, input, wrapFields)
	if err != nil {
		return err
	}

	err = r.validateOrder(ctx, input, loadedData.TransactionRepair, loadedData.UserKycData, loadedData.TransactionData.Type)
	if err != nil {
		return err
	}

	err = r.executeValidation(ctx, input, *loadedData.TransactionData, *loadedData.UserKycData, wrapFields)
	if err != nil {
		return err
	}

	reversalOutput, err := r.reverse(ctx, loadedData.TransactionData, input, wrapFields)
	if err != nil {
		return err
	}

	r.addTransactionReparation(
		ctx,
		reversalOutput.ReverseID,
		input,
		loadedData.UserKycData,
		loadedData.TransactionData.Operation,
		wrapFields,
	)

	buildMetricReverseSuccess(
		ctx, loadedData.TransactionData, loadedData.TransactionData.SiteID, input.FaqID, input.Requester.ClientApplication)

	return nil
}

func (r reparation) executeValidation(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	transactionData domain.TransactionData,
	userKycData domain.SearchKycOutput,
	wrapFields *field.WrappedFields,
) error {

	validationInput := domain.ValidationInput{
		PaymentID:       input.PaymentID,
		UserID:          input.UserID,
		SiteID:          input.SiteID,
		FaqID:           input.FaqID,
		TransactionData: transactionData,
		Type:            domain.TypeReverse,
		UserKycData:     userKycData,
		Requester:       input.Requester,
	}

	err := r.validation.ExecuteValidation(ctx, validationInput, wrapFields)
	if err != nil {
		if gnierrors.Type(err) == domain.NotEligible {
			addMetricReverseNotEligible(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
		}
		return buildErrorFromValidationService(err)
	}

	addMetricReverseEligible(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
	return nil
}

func (r reparation) reverse(
	ctx context.Context,
	transactionData *domain.TransactionData,
	input domain.ReverseTransactionInput,
	wrapFields *field.WrappedFields) (*domain.ReversalOutput, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerCardsTransactionReverse)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerCardsTransactionReverse(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	reversalInput := domain.ReversalInput{
		TransactionData:      transactionData,
		HeaderValueXClientID: input.HeaderValueXClientID,
	}

	reversalOutput, err := r.cardsTransactions.Reverse(ctx, reversalInput)
	if err != nil {
		if gnierrors.Type(err) == domain.AuthorizationAlreadyReversed {
			addMetricAuthorizationAlreadyReversed(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
			decrementMetricReverseEligible(ctx, input.SiteID, input.FaqID)
			return nil, buildErrorAuthorizationAlreadyReversed(err, msgErrorRequest)
		}

		addMetricErr(ctx, err, "reverse", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, reBuildError(ctx, r.log, domain.BadGateway.String(), err)
	}

	return reversalOutput, nil
}

func (r reparation) addTransactionReparation(
	ctx context.Context,
	reverseID string,
	input domain.ReverseTransactionInput,
	userKyc *domain.SearchKycOutput,
	operation domain.Operation,
	wrapFields *field.WrappedFields,
) {
	startTimer := wrapFields.Timers.Start(domain.TimerKvsRepairSave)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerRepairSave(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	storageReparationInput := r.fillStorageReparationInput(reverseID, input, userKyc, operation)

	err := r.reparationRepo.Save(ctx, storageReparationInput)
	if err != nil {
		addMetricWarn(ctx, err, "add_transaction_reparation", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		r.logErrorAndNotStop(ctx, fmt.Sprintf("could not insert transaction repair, by paymentID: %s", input.PaymentID), err)
	}
}

func (r reparation) logErrorAndNotStop(ctx context.Context, msgLog string, err error) {
	r.log.Warn(ctx, msgLog, r.log.Err(err))
}

func (r reparation) fillStorageReparationInput(
	reverseID string,
	input domain.ReverseTransactionInput,
	userKyc *domain.SearchKycOutput,
	operation domain.Operation,
) storage.ReparationInput {
	return storage.ReparationInput{
		BaseReparation: storage.BaseReparation{
			AuthorizationID:     operation.Authorization.ID,
			PaymentID:           input.PaymentID,
			KycIdentificationID: userKyc.KycIdentificationID,
			UserID:              input.UserID,
			TransactionRepairID: reverseID,
			SiteID:              input.SiteID,
			Type:                domain.TypeReverse,
			CreatedAt:           time.Now().UTC(),
			FaqID:               input.FaqID,
			Requester:           input.Requester,
			MonetaryTransaction: storage.MonetaryTransaction{
				Billing: storage.MonetaryValue{
					Amount:   operation.Billing.Amount / math.Pow10(operation.Billing.DecimalDigits),
					Currency: operation.Billing.Currency,
				},
				Settlement: storage.MonetaryValue{
					Amount:   operation.Settlement.Amount / math.Pow10(operation.Settlement.DecimalDigits),
					Currency: operation.Settlement.Currency,
				},
			},
		},
	}
}

func buildMetricReverseSuccess(
	ctx context.Context,
	transactionData *domain.TransactionData,
	siteID, faqID, clientApplication string,
) {
	daysToRepairRequest := transactionData.Operation.PaymentAge

	rangeOfDaysRepairRequest := "until_7_days"

	if daysToRepairRequest > _daysToRepairRequest {
		rangeOfDaysRepairRequest = "after_7_days"
	}

	addMetricReverseSuccess(
		ctx,
		siteID,
		daysToRepairRequest,
		rangeOfDaysRepairRequest,
		faqID,
		clientApplication,
		transactionData.Operation.Authorization.CardAcceptor.Name,
		transactionData.Operation.AcquirerCode)
}
