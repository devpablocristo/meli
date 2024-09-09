package graphqlkyc

import (
	"context"

	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchkycclient/v1"
)

type mockSearchKycClient struct {
	searchkycclient.Service
	SetDataByQueryRequestStub func(ctx context.Context, request graphql.Request, dataResponse interface{}) error
}

func (m mockSearchKycClient) SetDataByQueryRequest(
	ctx context.Context,
	request graphql.Request,
	dataResponse interface{},
) error {
	return m.SetDataByQueryRequestStub(ctx, request, dataResponse)
}
