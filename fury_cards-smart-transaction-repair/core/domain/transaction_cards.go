package domain

type TransactionWalletInput struct {
	PaymentID string
	UserID    string
}

type TransactionWalletOutput struct {
	TransactionData *TransactionData
	WalletData      *WalletData
}

type WalletData struct {
	ID       string
	CardsIDs []string
	Cards    []CardWallet
}

type CardWallet struct {
	ID             string
	BusinessMode   string
	IsTest         bool
	Holder         CardHolder
	IssuerAccounts []CardIssuerAccounts
}

type CardHolder struct {
	KycIdentificationID string
	VersionID           string
	UserID              string
}

type CardIssuerAccounts struct {
	ID string
}
