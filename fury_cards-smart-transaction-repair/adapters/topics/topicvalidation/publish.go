package topicvalidation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/keygenerator"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

const filterWorksForDocumentTrue = "works_for_document:true"
const filterWorksForDocumentFalse = "works_for_document:false"

func (v *validationEvent) Publish(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
	var errPublish error

	switch {
	case input.IsEligible:
		validationResult := createMessageValidation(input)
		setApprovedData(&validationResult, input)
		err := v.publisher.Publish(ctx, validationResult, filterWorksForDocumentTrue)
		if err != nil {
			errPublish = treatError(err, validationResult, wrapFields)
		}
	default:
		currentReason := 0
		for rule, result := range input.Reason {
			validationResult := createMessageValidation(input)
			complementMessage(&validationResult, rule, result)

			currentReason++
			if currentReason == 1 {
				setReason(&validationResult, input)
				setApprovedData(&validationResult, input)
				err := v.publisher.Publish(ctx, validationResult, filterWorksForDocumentTrue)
				if err != nil {
					errPublish = treatError(err, validationResult, wrapFields)
				}
				continue
			}

			err := v.publisher.Publish(ctx, validationResult, filterWorksForDocumentFalse)
			if err != nil {
				errPublish = treatError(err, validationResult, wrapFields)
			}
		}
	}
	return errPublish
}

func createMessageValidation(input storage.ValidationResult) validationResultTopic {
	return validationResultTopic{
		ID:                  keygenerator.CreateKey(),
		PaymentID:           input.PaymentID,
		KycIdentificationID: input.KycIdentificationID,
		UserID:              input.UserID,
		IsEligible:          input.IsEligible,
		CreatedAt:           input.DateCreated,
		AuthorizationID:     input.AuthorizationID,
		Type:                input.Type,
		SiteID:              input.SiteID,
		FaqID:               input.FaqID,
		ClientApplication:   input.Requester.ClientApplication,
		ClientScope:         input.Requester.ClientScope,
		Amount:              input.TransactionEvaluated.Amount,
		Currency:            input.TransactionEvaluated.Currency,
		Score:               input.Score,
	}
}

func setReason(validation *validationResultTopic, input storage.ValidationResult) {
	validation.Reason = map[domain.RuleType]reasonDefault{}
	setResult(validation.Reason, input.Reason)
}

func setApprovedData(validation *validationResultTopic, input storage.ValidationResult) {
	validation.ApprovedData = map[domain.RuleType]reasonDefault{}
	setResult(validation.ApprovedData, input.ApprovedData)
}

func setResult(reason map[domain.RuleType]reasonDefault, result map[domain.RuleType]*domain.ReasonResult) {
	for k, v := range result {
		switch accepted := v.Accepted.(type) {
		case domain.ReasonResultPerPeriod:
			reason[k] = reasonDefault{
				Actual: v.Actual,
				Accepted: reasonResultPerPeriod{
					Qty:        accepted.Qty,
					PeriodDays: accepted.PeriodDays,
				},
			}
		default:
			reason[k] = reasonDefault(*v)
		}
	}
}

func complementMessage(validationResult *validationResultTopic, rule domain.RuleType, result *domain.ReasonResult) {
	validationResult.Rule = string(rule)
	validationResult.AllowedValue.toString(result.Accepted)
	validationResult.ValueDeclined.toString(result.Actual)
}

func treatError(err error, validationResult validationResultTopic, wrapFields *field.WrappedFields) error {
	body, _ := json.Marshal(validationResult)
	wrapFields.Fields.Add(fmt.Sprintf(`message_body_from_id_%s`, validationResult.ID), string(body))
	wrapFields.Fields.Add(fmt.Sprintf(`error_message_body_from_id_%s`, validationResult.ID), err.Error())

	return buildErrorBadGateway(err)
}
