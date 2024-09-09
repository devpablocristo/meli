package graphqltransaction

import (
	"context"
	"errors"

	_ "embed"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

type mock struct {
	JSON json.JSON
}

//go:embed testdata/transaction-data.json
var _transactionDataMock []byte

func (m mock) GetTransactionByPaymentID(
	ctx context.Context,
	input *domain.TransactionInput,
	wrapFields *field.WrappedFields) (*domain.TransactionOutput, error) {

	if input.PaymentID == "" {
		return nil, errors.New("not found")
	}

	var outputExpected *domain.TransactionOutput
	err := m.JSON.Unmarshal(_transactionDataMock, &outputExpected)

	return outputExpected, err
}

var _ ports.SearchTransaction = (*mock)(nil)

func NewMock() ports.SearchTransaction {
	return mock{
		json.NewJSON(),
	}
}
