package graphqlhub

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/adapters/graphqlhub/transactionhistory"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_searchHub_GetTransactionHistory_Success(t *testing.T) {
	mockService := mockSearchHubClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			response := `{"wallet": {
							"cards": [
								{
									"issuer_accounts": [
										{
											"id": "GKALLNGKUBIEYQXVYBHPIYDLDMAWOADJIJSEDZAW",
											"transactions": {									
												"transactions": [
													{											
														"operation": {
															"settlement": {
																"currency": "USD",
																"amount": null,
																"decimal_digits": null
															},
															"billing": {
																"amount": 366,
																"currency": "BRL",
																"decimal_digits": 2
															}
														}										
													},
													{											
														"operation": {
															"settlement": {
																"currency": null,
																"amount": null,
																"decimal_digits": null
															},
															"billing": {
																"amount": 335,
																"currency": "BRL",
																"decimal_digits": 2
															}
														}
													}
												]
											}
										}
									]
								}
							]
						}
					}`

			var history transactionhistory.History
			err := json.Unmarshal([]byte(response), &history)
			assert.NoError(t, err)

			v := dataResponse.(*transactionhistory.History)
			*v = history
			return nil
		},
	}

	service := New(mockService)

	input := domain.TransactionHistoryInput{
		UserID:     "1117823961",
		SizeSearch: 100,
	}

	output, err := service.GetTransactionHistory(context.TODO(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.IssuerAccounts, 1)
	assert.Len(t, output.IssuerAccounts[0].Transactions, 2)
}

func Test_searchHub_GetTransactionHistory_Error(t *testing.T) {
	mockService := mockSearchHubClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			return mockErrorNotFoundWallet("1117823961666")
		},
	}

	service := New(mockService)

	input := domain.TransactionHistoryInput{
		UserID:     "1117823961666",
		SizeSearch: 100,
	}

	output, err := service.GetTransactionHistory(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(mockErrorNotFoundWallet("1117823961666").Error()), "search-hub get_transaction_history failure", false)

	assert.Nil(t, output)
	assertErrorFromNotFoundWallet(t, errorExpected, err)
}
