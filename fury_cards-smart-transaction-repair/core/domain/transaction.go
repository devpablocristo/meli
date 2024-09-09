package domain

import "time"

type TransactionInput struct {
	PaymentID string
}

type TransactionOutput struct {
	TransactionData *TransactionData
}

type TransactionData struct {
	Operation          Operation
	Provider           Provider
	Environment        string
	SiteID             string
	StatusDetail       string
	PayerID            int64
	Type               string
	IsRecurringPayment bool
}

type Operation struct {
	SubType                string
	CreationDatetime       time.Time
	TransmissionDatetime   time.Time
	Installments           int
	AcquirerCode           string
	Stan                   int
	Authorization          Authorization
	Transaction            Transaction
	Card                   Card
	IsAdvice               bool
	Billing                Billing
	Settlement             Settlement
	PaymentAge             int
	AccountType            string
	AuthorizationIndicator string
	IsInternational        bool
}

type Provider struct {
	ID string
}

type Options struct {
	ReversalType string
}

type Authorization struct {
	ID                   string
	Stan                 int
	TransmissionDatetime time.Time
	CardAcceptor         CardAcceptor
}

type Transaction struct {
	Amount        float64
	Currency      string
	DecimalDigits int
	TotalAmount   float64
}

type Billing struct {
	Amount        float64
	Currency      string
	DecimalDigits int
	TotalAmount   float64
	Conversion    Conversion
}

type Settlement struct {
	Amount        float64
	Currency      string
	DecimalDigits int
	TotalAmount   float64
	Conversion    Conversion
}

type Card struct {
	NumberID      string
	Country       string
	IssuerAccount IssuerAccount
	TokenInfo     TokenInfo
}

type CardAcceptor struct {
	Terminal      string
	Name          string
	EntryMode     string
	Type          string
	TerminalMode  string
	PinCapability string
}

type Conversion struct {
	Date          *time.Time
	DecimalDigits *int
	Rate          *float64
	From          *string
}

type IssuerAccount struct {
	CardBusinessMode string
}

type TokenInfo struct {
	TokenWalletID string
}
