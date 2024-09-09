package ds

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

func (d dsValidationResult) Save(ctx context.Context, paymentID string, input storage.ValidationResult) error {
	validationResult := d.fillValidationResultDB(input)

	if err := d.repo.SaveDocumentWithContext(ctx, validationResult.PaymentID, validationResult); err != nil {
		return buildErrorBadGateway(err)
	}
	return nil
}

func (d dsValidationResult) fillValidationResultDB(input storage.ValidationResult) validationResultDB {
	result := validationResultDB{
		PaymentID:            input.PaymentID,
		KycIdentificationID:  input.KycIdentificationID,
		UserID:               input.UserID,
		IsEligible:           input.IsEligible,
		CreatedAt:            input.DateCreated,
		AuthorizationID:      input.AuthorizationID,
		Type:                 input.Type,
		SiteID:               input.SiteID,
		FaqID:                input.FaqID,
		Requester:            requesterDB(input.Requester),
		TransactionEvaluated: transactionEvaluatedDB(input.TransactionEvaluated),
		Score:                input.Score,
	}

	result.ApprovedData = map[domain.RuleType]reasonDefaultDB{}
	setResult(result.ApprovedData, input.ApprovedData)

	if result.IsEligible {
		return result
	}

	result.Reason = map[domain.RuleType]reasonDefaultDB{}
	setResult(result.Reason, input.Reason)

	return result
}

func setResult(reasonDB map[domain.RuleType]reasonDefaultDB, reason map[domain.RuleType]*domain.ReasonResult) {
	for k, v := range reason {
		switch accepted := v.Accepted.(type) {
		case domain.ReasonResultPerPeriod:
			reasonDB[k] = reasonDefaultDB{
				Actual: v.Actual,
				Accepted: reasonResultPerPeriodDB{
					Qty:        accepted.Qty,
					PeriodDays: accepted.PeriodDays,
				},
			}
		default:
			reasonDB[k] = reasonDefaultDB(*v)
		}
	}
}
