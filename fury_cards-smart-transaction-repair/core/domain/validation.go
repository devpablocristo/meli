package domain

type TypeReparation string

const (
	TypeReverse TypeReparation = "reverse"
)

type ValidationInput struct {
	PaymentID       string
	UserID          int64
	SiteID          string
	FaqID           string
	Type            TypeReparation
	TransactionData TransactionData
	UserKycData     SearchKycOutput
	Requester       Requester
}
