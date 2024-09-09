package payment

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchhubclient/v1"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/response_payment_transaction.json
var _paymentTransactionResponse []byte

type mockSearchHubClient struct {
	searchhubclient.Service
	SetDataByQueryRequestStub func(ctx context.Context, request graphql.Request, dataResponse interface{}) error
}

func (m mockSearchHubClient) SetDataByQueryRequest(
	ctx context.Context,
	request graphql.Request,
	dataResponse interface{},
) error {
	return m.SetDataByQueryRequestStub(ctx, request, dataResponse)
}

func mockErrorResponse(userID string) error {
	return fmt.Errorf(`{
		"errors":
			[{
				"message":"error [not_found]: Not found wallet with user_id: %s",
				"path":["wallet"],
				"extensions":{"code":"not_found","message":"Not found wallet with user_id: %s"}
			}],
		"status_code":200}`, userID, userID)
}

func assertErrorFromResponse(t *testing.T, errorExpected, err error) {
	assert.Error(t, err)
	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.JSONEq(t, gnierrors.Cause(errorExpected).Error(), gnierrors.Cause(err).Error())
}
