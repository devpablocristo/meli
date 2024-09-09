package graphqlhub

import (
	"context"
	_ "embed"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

func Test_mock_CountChargebackFromUser(t *testing.T) {
	m := NewMock()
	output, err := m.CountChargebackFromUser(context.TODO(), domain.SearchHubCountChargebackInput{UserID: "12"})

	assert.NoError(t, err)
	assert.Equal(t, output.Total, 0)
}

func Test_mock_GetTransactionAndCards(t *testing.T) {
	m := NewMock()

	input := domain.TransactionWalletInput{
		UserID:    "123",
		PaymentID: "456",
	}

	output, err := m.GetTransactionAndCards(context.TODO(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func Test_mock_GetTransactionHistory(t *testing.T) {
	m := NewMock()

	input := domain.TransactionHistoryInput{
		UserID:     "123",
		SizeSearch: 100,
	}

	output, err := m.GetTransactionHistory(context.TODO(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func Test_mock_GetPaymentTransaction(t *testing.T) {
	m := mock{}
	wrapFields := &field.WrappedFields{Fields: field.NewWrappedFields().Fields}

	response, err := m.GetPaymentTransaction(context.Background(), domain.TransactionInput{PaymentID: "123"}, wrapFields)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}
