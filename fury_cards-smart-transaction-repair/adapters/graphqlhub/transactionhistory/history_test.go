package transactionhistory

import (
	"encoding/json"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestBuildRequestGetTransactionHistory(t *testing.T) {
	input := domain.TransactionHistoryInput{
		UserID:     "123",
		SizeSearch: 10,
	}

	t.Run("with size", func(t *testing.T) {
		queryOperationsExpected := `{
		"query": {
		  "and": [
			{
			  "eq": {
				"field": "type",
				"value": "authorization"
			  }
			},
			{
			  "exists": {
				"field": "capture_id"
			  }
			}
		  ]
		},
		"size": 10
	  }
	`
		request := BuildRequestGetTransactionHistory(input)

		bQueryParams, err := json.Marshal(request.QueryOperations["queryParams"])
		assert.NoError(t, err)

		assert.Equal(t, request.Params["id"], "123")
		assert.JSONEq(t, string(bQueryParams), queryOperationsExpected)
	})

	t.Run("with size default", func(t *testing.T) {
		input := domain.TransactionHistoryInput{
			UserID: "123",
		}

		request := BuildRequestGetTransactionHistory(input)

		assert.Equal(t, 100, request.QueryOperations["queryParams"].Size)
	})
}

func TestFillOutput(t *testing.T) {
	data := History{
		Wallet: &wallet{
			Cards: []card{
				{
					IssuerAccounts: []issuerAccount{
						{
							ID: "i1",
							Transactions: transactionAccounts{
								MatchTotal: 1,
								Transactions: []transactionOperation{
									{
										ID: "1", Operation: operation{
											Billing: billing{
												Amount:        10,
												Currency:      "AR",
												DecimalDigits: 0,
											},
											Settlement: settlement{
												Amount:        2,
												Currency:      "USD",
												DecimalDigits: 0,
											},
										},
									},
								},
							},
						},
					},
				},
				{
					IssuerAccounts: []issuerAccount{
						{
							ID: "i2",
							Transactions: transactionAccounts{
								MatchTotal: 1,
								Transactions: []transactionOperation{
									{
										ID: "1", Operation: operation{
											Billing: billing{
												Amount:        10,
												Currency:      "AR",
												DecimalDigits: 0,
											},
											Settlement: settlement{
												Amount:        2,
												Currency:      "USD",
												DecimalDigits: 0,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	output := FillOutput(data)
	assert.Len(t, output.IssuerAccounts, 2)
}
