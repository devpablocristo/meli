package policy

import (
	"context"

	"github.com/melisource/fury_cards-go-toolkit/pkg/policyclient/v1"
)

type mockPolicyClient struct {
	EvaluateStub func(ctx context.Context, evalID int64, memberType string, tagName ...string) (*policyclient.PolicyOutput, error)
}

func (m mockPolicyClient) Evaluate(
	ctx context.Context,
	evalID int64,
	memberType string,
	tagName ...string,
) (*policyclient.PolicyOutput, error) {
	return m.EvaluateStub(ctx, evalID, memberType, tagName...)
}
