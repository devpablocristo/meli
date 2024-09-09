package payment

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentTransaction_Success(t *testing.T) {
	mockService := mockSearchHubClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			response := string(_paymentTransactionResponse)

			var dataPaymentTransaction paymentTransaction
			err := json.Unmarshal([]byte(response), &dataPaymentTransaction)
			assert.NoError(t, err)

			v := dataResponse.(*paymentTransaction)
			*v = dataPaymentTransaction
			return nil
		},
	}

	input := domain.TransactionInput{
		PaymentID: "54598550619",
	}

	wrapFields := &field.WrappedFields{Fields: field.NewWrappedFields().Fields}

	output, err := GetPaymentTransaction(context.TODO(), mockService, input, wrapFields)
	assert.NoError(t, err)
	assert.NotNil(t, output.TransactionData)

	utilstest.AssertFieldsNoEmptyFromStruct(t, *output.TransactionData)
}

func TestGetPaymentTransaction_Error(t *testing.T) {
	mockService := mockSearchHubClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			return mockErrorResponse("1117823961")
		},
	}

	input := domain.TransactionInput{
		PaymentID: "54598550619",
	}

	wrapFields := &field.WrappedFields{Fields: field.NewWrappedFields().Fields}

	output, err := GetPaymentTransaction(context.TODO(), mockService, input, wrapFields)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(mockErrorResponse("1117823961").Error()), "searchhub.get_payment_transaction.failure", false)

	assert.Empty(t, output)
	assertErrorFromResponse(t, errorExpected, err)
}
