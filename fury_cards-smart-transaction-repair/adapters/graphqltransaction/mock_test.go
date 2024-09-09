package graphqltransaction

import (
	"context"
	_ "embed"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

func Test_mock_GetTransactionByPaymentID_Success(t *testing.T) {
	mock := NewMock()

	input := &domain.TransactionInput{
		PaymentID: "100",
	}

	output, err := mock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}

func Test_mock_GetTransactionByPaymentID_Errors(t *testing.T) {
	mock := NewMock()

	t.Run("without paymentID", func(t *testing.T) {
		input := &domain.TransactionInput{}

		output, err := mock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

		assert.Error(t, err)
		assert.Empty(t, output)
	})
}
