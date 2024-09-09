package paymentcards

import (
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/stretchr/testify/assert"
)

func TestBuildRequestTransactionAndCards(t *testing.T) {
	input := domain.TransactionWalletInput{
		UserID:    "123",
		PaymentID: "abc",
	}

	request := BuildRequestTransactionAndCards(input)

	assert.Equal(t, "123", request.Params["user_id"])
	assert.Equal(t, "abc", request.Params["payment_id"])
}

func TestFillOutput(t *testing.T) {
	data := TransactionWalletSearch{
		Payment: &payment{
			DebitTransaction: debitTransaction{
				Operation: operation{},
			},
		},
		Wallet: &wallet{
			Cards: []cardWallet{
				{
					IssuerAccounts: []issuerAccounts{
						{
							ID: "i1",
						},
					},
				},
			},
		},
	}

	output := FillOutput(data)
	utilstest.AssertFieldsNoEmptyFromStruct(t, *output)
}
