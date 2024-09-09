package validationresult

import (
	"context"
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_go-core/pkg/web"
)

func (v ValidationResultConsumer) getRequest(
	ctx context.Context,
	r *http.Request,
	wrapFields *field.WrappedFields,
) (*msgValidationResult, error) {

	var request msgTopicRequest

	if err := web.DecodeJSON(r, &request); err != nil {
		addMetricIvalidPayload(ctx)
		return nil, v.buildErrorBadRequestAndLog(ctx, err, wrapFields)
	}

	return &request.Body, nil
}

func (v ValidationResultConsumer) buildValidationResultInput(
	request msgValidationResult,
	wrapFields *field.WrappedFields,
) storage.ValidationResult {

	input := storage.ValidationResult{
		PaymentID:           request.PaymentID,
		AuthorizationID:     request.AuthorizationID,
		KycIdentificationID: request.KycIdentificationID,
		UserID:              request.UserID,
		SiteID:              request.SiteID,
		IsEligible:          request.IsEligible,
		DateCreated:         request.CreatedAt,
		Type:                request.Type,
		FaqID:               request.FaqID,
		Requester: domain.Requester{
			ClientApplication: request.ClientApplication,
			ClientScope:       request.ClientScope,
		},
		TransactionEvaluated: storage.TransactionEvaluated{
			Amount:   request.Amount,
			Currency: request.Currency,
		},
		Score:        request.Score,
		ApprovedData: fillResult(request.ApprovedData),
		Reason:       fillResult(request.Reason),
	}

	setFieldsInputParameters(wrapFields, request)

	return input
}

func fillResult(requestReason map[domain.RuleType]reasonDefault) map[domain.RuleType]*domain.ReasonResult {
	reason := domain.Reason{}

	for k, v := range requestReason {
		reason[k] = v.getValue()
	}

	return reason
}
