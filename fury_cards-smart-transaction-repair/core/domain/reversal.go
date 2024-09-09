package domain

type ReversalInput struct {
	TransactionData      *TransactionData
	HeaderValueXClientID string
}

type ReversalOutput struct {
	ReverseID string
}
