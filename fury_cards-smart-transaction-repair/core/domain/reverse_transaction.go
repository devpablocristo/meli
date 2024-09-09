package domain

type ReverseTransactionInput struct {
	PaymentID            string
	UserID               int64
	SiteID               string
	HeaderValueXClientID string
	FaqID                string
	Requester            Requester
}
