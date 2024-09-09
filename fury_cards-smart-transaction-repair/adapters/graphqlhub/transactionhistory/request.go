package transactionhistory

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
)

func BuildRequestGetTransactionHistory(transactionHistoryInput domain.TransactionHistoryInput) graphql.Request {
	sizeDefault := 100
	if transactionHistoryInput.SizeSearch > 0 {
		sizeDefault = transactionHistoryInput.SizeSearch
	}

	return graphql.Request{
		Query: queryTransactionHistory,
		Params: map[string]interface{}{
			"id": transactionHistoryInput.UserID,
		},
		QueryOperations: map[string]graphql.QueryParams{
			"queryParams": {
				Query: graphql.QueryOperation{
					And: []graphql.Clauses{
						{
							Eq: &graphql.ClauseEq{Field: "type", Value: "authorization"},
						},
						{
							Exists: &graphql.ClauseExists{Field: "capture_id"},
						},
					},
				},
				Size: sizeDefault,
			},
		},
	}
}
