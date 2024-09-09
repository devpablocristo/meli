package reparation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

const operationTypeAuthorization = "authorization"

func (r reparation) validateOrder(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	transactionRepair *storage.ReparationOutput,
	kycData *domain.SearchKycOutput,
	operationType string,
) error {

	err := r.validateInput(ctx, input)
	if err != nil {
		return err
	}

	err = r.validateIsRepaired(ctx, transactionRepair, input)
	if err != nil {
		return err
	}

	err = r.validateIdentificationExists(ctx, input, kycData)
	if err != nil {
		return err
	}

	err = r.validateIsAuthorization(ctx, operationType, input)
	if err != nil {
		return err
	}

	return nil
}

func (r reparation) validateInput(ctx context.Context, input domain.ReverseTransactionInput) error {
	var requiredFields string

	if strings.TrimSpace(input.PaymentID) == "" {
		requiredFields = fmt.Sprint(requiredFields, "PaymentID; ")
	}

	if input.UserID == 0 {
		requiredFields = fmt.Sprint(requiredFields, "UserID; ")
	}

	if requiredFields != "" {
		err := fmt.Errorf("required fields: %s", strings.TrimSpace(requiredFields))
		addMetricErr(ctx, err, "validate_order", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return buildErrorBadRequest(err, msgErrorRequest)
	}

	return nil
}

func (r reparation) validateIsRepaired(
	ctx context.Context,
	transactionRepair *storage.ReparationOutput,
	input domain.ReverseTransactionInput,
) error {

	if transactionRepair != nil {
		addMetricReverseAlready(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return buildErrorAlreadyRepaired(errors.New("reverse already done"), msgErrorRequest)
	}

	return nil
}

func (r reparation) validateIdentificationExists(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	kycData *domain.SearchKycOutput,
) error {

	if kycData.KycIdentificationID == "" || kycData.KycIdentificationID == "-" {
		addMetricUnablePersonIDEmpty(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return buildErrorUnableRequest(fmt.Errorf("kyc_identification from userID:%v is empty", input.UserID), msgErrorRequest)
	}

	return nil
}

func (r reparation) validateIsAuthorization(
	ctx context.Context,
	operationType string,
	input domain.ReverseTransactionInput,
) error {

	if operationType != operationTypeAuthorization {
		addMetricUnableTypeDiffAuthorization(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return buildErrorUnableRequest(fmt.Errorf("type of operation cannot be:%s", operationType), msgErrorRequest)
	}

	return nil
}
