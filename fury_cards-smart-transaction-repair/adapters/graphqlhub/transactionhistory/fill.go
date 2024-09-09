package transactionhistory

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

func FillOutput(data History) *domain.TransactionHistoryOutput {
	return &domain.TransactionHistoryOutput{
		IssuerAccounts: fillIssuerAccounts(data.Wallet.Cards),
	}
}

func fillIssuerAccounts(cards []card) []domain.HistoryIssuerAccounts {
	issuerAccounts := []domain.HistoryIssuerAccounts{}
	for _, card := range cards {
		for _, account := range card.IssuerAccounts {
			issuerAccounts = append(issuerAccounts, domain.HistoryIssuerAccounts{
				ID:           account.ID,
				MatchTotal:   account.Transactions.MatchTotal,
				Transactions: fillTransactionOperation(account.Transactions.Transactions),
			})
		}
	}

	return issuerAccounts
}

func fillTransactionOperation(transactions []transactionOperation) []domain.HistoryOperation {
	transactionsData := []domain.HistoryOperation{}
	for _, tr := range transactions {
		transactionsData = append(transactionsData, domain.HistoryOperation{
			Billing: domain.Billing{
				Amount:        tr.Operation.Billing.Amount,
				Currency:      tr.Operation.Billing.Currency,
				DecimalDigits: tr.Operation.Billing.DecimalDigits,
			},
			Settlement: domain.Settlement{
				Amount:        tr.Operation.Settlement.Amount,
				Currency:      tr.Operation.Settlement.Currency,
				DecimalDigits: tr.Operation.Settlement.DecimalDigits,
			},
		})
	}

	return transactionsData
}
