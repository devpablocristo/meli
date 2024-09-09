package graphqlhub

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/adapters/graphqlhub/paymentcards"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_searchHub_GetTransactionAndCards(t *testing.T) {
	mockService := mockSearchHubClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			response := `{"payment":{},"wallet":{"cards":[{"id":"card1","issuer_accounts":[{"id":"account01"}]}]}}`
			var dataTransactionWalletSearch paymentcards.TransactionWalletSearch
			err := json.Unmarshal([]byte(response), &dataTransactionWalletSearch)
			assert.NoError(t, err)

			v := dataResponse.(*paymentcards.TransactionWalletSearch)
			*v = dataTransactionWalletSearch
			return nil
		},
	}

	service := New(mockService)

	input := domain.TransactionWalletInput{
		UserID:    "35315689",
		PaymentID: "54598550619",
	}

	output, err := service.GetTransactionAndCards(context.TODO(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output.TransactionData)
	assert.NotNil(t, output.WalletData)
}

func Test_searchHub_GetTransactionAndCards_Error(t *testing.T) {
	mockService := mockSearchHubClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			return mockErrorNotFoundWallet("1117823961")
		},
	}

	service := New(mockService)

	input := domain.TransactionWalletInput{
		UserID:    "1117823961",
		PaymentID: "54598550619",
	}

	output, err := service.GetTransactionAndCards(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(mockErrorNotFoundWallet("1117823961").Error()), "search-hub get_transaction_and_cards failure", false)

	assert.Nil(t, output)
	assertErrorFromNotFoundWallet(t, errorExpected, err)
}
