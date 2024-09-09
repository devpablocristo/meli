package policy

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

func (e evaluate) EvaluateWithUser(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
	res, err := e.pClient.Evaluate(ctx, userID, domain.UserType, tagName...)

	if err != nil {
		return nil, buildErrorBadGateway(err)
	}

	return fillPolicyOutput(res.Result, res.FailedPolicies), nil
}

func fillPolicyOutput(isAuth bool, failedRestricts []string) *domain.PolicyOutput {
	return &domain.PolicyOutput{
		IsAuthorized:    isAuth,
		RestrictsFailed: failedRestricts,
	}
}
