package payment

import (
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/stretchr/testify/assert"
)

func TestBuildRequestPaymentTransaction(t *testing.T) {
	input := domain.TransactionInput{
		PaymentID: "abc",
	}

	request := BuildRequestPaymentTransaction(input)

	assert.Equal(t, "abc", request.Params["id"])
}

func TestFillOutput(t *testing.T) {
	data := paymentTransaction{
		Payment: &payment{
			DebitTransaction: debitTransaction{
				Operation: operation{},
			},
		},
	}

	output := FillOutput(data)
	utilstest.AssertFieldsNoEmptyFromStruct(t, output)
}
