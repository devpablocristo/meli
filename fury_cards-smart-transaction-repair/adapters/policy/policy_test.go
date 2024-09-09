package policy

import (
	"context"
	"errors"
	"testing"

	"github.com/melisource/fury_cards-go-toolkit/pkg/policyclient/v1"
	"github.com/stretchr/testify/assert"
)

func TestPolicy_EvaluateWithUser_Success_withAllowed(t *testing.T) {
	m := mockPolicyClient{
		EvaluateStub: func(ctx context.Context, evalID int64, memberType string, tagName ...string) (*policyclient.PolicyOutput, error) {
			assert.Equal(t, "user", memberType)
			return &policyclient.PolicyOutput{
				Result:  true,
				Message: "OK",
			}, nil
		},
	}
	c := New(m)

	res, err := c.EvaluateWithUser(context.Background(), 1, "")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, true, res.IsAuthorized)
	assert.Equal(t, 0, len(res.RestrictsFailed))
}

func TestPolicy_EvaluateWithUser_Success_withNotAllowed(t *testing.T) {
	m := mockPolicyClient{
		EvaluateStub: func(ctx context.Context, evalID int64, memberType string, tagName ...string) (*policyclient.PolicyOutput, error) {
			assert.Equal(t, "user", memberType)
			return &policyclient.PolicyOutput{
				Result:         false,
				Message:        "OK",
				FailedPolicies: []string{"id_policy"},
			}, nil
		},
	}
	c := New(m)

	res, err := c.EvaluateWithUser(context.Background(), 1, "")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, false, res.IsAuthorized)
	assert.Equal(t, 1, len(res.RestrictsFailed))
}

func TestPolicy_EvaluateWithUser_Failed(t *testing.T) {
	m := mockPolicyClient{
		EvaluateStub: func(ctx context.Context, evalID int64, memberType string, tagName ...string) (*policyclient.PolicyOutput, error) {
			assert.Equal(t, "user", memberType)
			return nil, errors.New("error in evaluate")
		},
	}
	c := New(m)

	res, err := c.EvaluateWithUser(context.Background(), 1, "")

	assert.Nil(t, res)
	assert.Equal(t, "transaction repair policy validated failure : error in evaluate", err.Error())
}
