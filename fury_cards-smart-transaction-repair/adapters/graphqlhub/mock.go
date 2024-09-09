package graphqlhub

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

type mock struct{}

func (m mock) CountChargebackFromUser(
	ctx context.Context,
	input domain.SearchHubCountChargebackInput,
) (*domain.SearchHubCountChargebackOutput, error) {
	return &domain.SearchHubCountChargebackOutput{}, nil
}

func (m mock) GetTransactionAndCards(
	ctx context.Context,
	transactionInput domain.TransactionWalletInput,
) (*domain.TransactionWalletOutput, error) {
	return &domain.TransactionWalletOutput{}, nil
}

func (m mock) GetTransactionHistory(
	ctx context.Context,
	transactionHistoryInput domain.TransactionHistoryInput,
) (*domain.TransactionHistoryOutput, error) {
	return &domain.TransactionHistoryOutput{}, nil
}

func (m mock) GetPaymentTransaction(
	ctx context.Context,
	transactionInput domain.TransactionInput,
	wrapFields *field.WrappedFields,
) (domain.TransactionOutput, error) {
	return domain.TransactionOutput{TransactionData: &domain.TransactionData{}}, nil
}

var _ ports.SearchHub = (*mock)(nil)

func NewMock() ports.SearchHub {
	return mock{}
}
