package kvs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

const logRequiredNoStop = "empty SiteID when saving repair transaction for authorization: %s"

func (k kvsRepository) Save(ctx context.Context, transactionRepair storage.ReparationInput) error {
	if err := k.validateSave(ctx, transactionRepair); err != nil {
		return err
	}

	transacRepair := domainToTransactionRepairDB(transactionRepair)

	// Why Put? The "get" method is done at the beginning of the reverse process.
	if err := k.repo.Put(ctx, transacRepair.AuthorizationID, transacRepair); err != nil {
		return buildErrorBadGateway(err)
	}
	return nil
}

func (k kvsRepository) validateSave(ctx context.Context, transactionRepair storage.ReparationInput) error {
	var requiredFields string
	if transactionRepair.AuthorizationID == "" {
		requiredFields = fmt.Sprint(requiredFields, "AuthorizationID; ")
	}
	if transactionRepair.PaymentID == "" {
		requiredFields = fmt.Sprint(requiredFields, "PaymentID; ")
	}
	if strings.TrimSpace(transactionRepair.KycIdentificationID) == "" || transactionRepair.KycIdentificationID == "-" {
		requiredFields = fmt.Sprint(requiredFields, "KycIdentificationID; ")
	}
	if transactionRepair.UserID == 0 {
		requiredFields = fmt.Sprint(requiredFields, "UserID; ")
	}
	if transactionRepair.SiteID == "" {
		k.log.Warn(ctx, fmt.Sprintf(logRequiredNoStop, transactionRepair.AuthorizationID))
	}
	if transactionRepair.Type == "" {
		k.log.Warn(ctx, fmt.Sprintf(logRequiredNoStop, transactionRepair.AuthorizationID))
	}
	if requiredFields != "" {
		msg := fmt.Sprint("required fields: ", strings.TrimSpace(requiredFields))
		return buildErrorBadRequest(errors.New(msg))
	}
	return nil
}

func domainToTransactionRepairDB(transactionRepair storage.ReparationInput) *dbreparation.TransactionRepairDB {
	return &dbreparation.TransactionRepairDB{
		AuthorizationID:     transactionRepair.AuthorizationID,
		PaymentID:           transactionRepair.PaymentID,
		KycIdentificationID: transactionRepair.KycIdentificationID,
		UserID:              transactionRepair.UserID,
		TransactionRepairID: transactionRepair.TransactionRepairID,
		SiteID:              transactionRepair.SiteID,
		Type:                transactionRepair.Type,
		CreatedAt:           transactionRepair.CreatedAt,
		FaqID:               transactionRepair.FaqID,
		Requester:           dbreparation.RequesterDB(transactionRepair.Requester),
		MonetaryTransaction: toMonetaryTransactionDB(transactionRepair.MonetaryTransaction),
	}
}

func toMonetaryTransactionDB(monetaryReparation storage.MonetaryTransaction) dbreparation.MonetaryTransactionDB {
	return dbreparation.MonetaryTransactionDB{
		Billing:    dbreparation.MonetaryValueDB(monetaryReparation.Billing),
		Settlement: dbreparation.MonetaryValueDB(monetaryReparation.Settlement),
	}
}
