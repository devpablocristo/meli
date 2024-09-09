package domain

type TransactionHistoryInput struct {
	UserID     string
	SizeSearch int
}

type TransactionHistoryOutput struct {
	IssuerAccounts []HistoryIssuerAccounts
}

type HistoryIssuerAccounts struct {
	ID           string
	MatchTotal   int
	Transactions []HistoryOperation
}

type HistoryOperation struct {
	Billing    Billing
	Settlement Settlement
}
