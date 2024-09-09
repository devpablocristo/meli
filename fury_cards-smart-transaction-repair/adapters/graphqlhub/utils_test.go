package graphqlhub

import (
	"context"
	"fmt"
	"testing"

	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchhubclient/v1"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

type mockSearchHubClient struct {
	searchhubclient.Service
	CountChargebacksInTheLastDaysStub func(ctx context.Context, userID string, lastDays int) (*searchhubclient.OutputCountChargeback, error)
	SetDataByQueryRequestStub         func(ctx context.Context, request graphql.Request, dataResponse interface{}) error
}

func (m mockSearchHubClient) CountChargebacksInTheLastDays(
	ctx context.Context,
	userID string,
	lastDays int,
) (*searchhubclient.OutputCountChargeback, error) {
	return m.CountChargebacksInTheLastDaysStub(ctx, userID, lastDays)
}

func (m mockSearchHubClient) SetDataByQueryRequest(
	ctx context.Context,
	request graphql.Request,
	dataResponse interface{},
) error {
	return m.SetDataByQueryRequestStub(ctx, request, dataResponse)
}

func mockErrorNotFoundWallet(userID string) error {
	return fmt.Errorf(`{
		"errors":
			[{
				"message":"error [not_found]: Not found wallet with user_id: %s",
				"path":["wallet"],
				"extensions":{"code":"not_found","message":"Not found wallet with user_id: %s"}
			}],
		"status_code":200}`, userID, userID)
}

func assertErrorFromNotFoundWallet(t *testing.T, errorExpected, err error) {
	assert.Error(t, err)
	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.JSONEq(t, gnierrors.Cause(errorExpected).Error(), gnierrors.Cause(err).Error())
}
