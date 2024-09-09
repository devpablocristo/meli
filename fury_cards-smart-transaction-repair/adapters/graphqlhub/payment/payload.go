package payment

import "time"

type paymentTransaction struct {
	Payment *payment `json:"payment"`
}

type payment struct {
	SiteID           string           `json:"site_id"`
	StatusDetail     string           `json:"status_detail"`
	Status           string           `json:"status"`
	DebitTransaction debitTransaction `json:"debit_transaction"`
	PayerID          int64            `json:"payer_id"`
}

type debitTransaction struct {
	ID                 string       `json:"id"`
	AcquirerCode       string       `json:"acquirer_code"`
	CaptureDatetime    time.Time    `json:"capture_datetime"`
	CaptureID          string       `json:"capture_id"`
	Environment        string       `json:"environment"`
	Operation          operation    `json:"operation"`
	ReversalsIDs       []string     `json:"reversal_ids"`
	Status             string       `json:"status"`
	Type               string       `json:"type"`
	CardAcceptor       cardAcceptor `json:"card_acceptor"`
	Provider           provider     `json:"provider"`
	IsRecurringPayment bool         `json:"is_recurring_payment"`
}

type operation struct {
	AccountType            string      `json:"account_type"`
	AuthorizationIndicator string      `json:"authorization_indicator"`
	IsInternational        bool        `json:"is_international"`
	CreationDatetime       time.Time   `json:"creation_datetime"`
	SubType                string      `json:"subtype"`
	TransmissionDatetime   time.Time   `json:"transmission_datetime"`
	IsAdvice               bool        `json:"is_advice"`
	ExpirationDatetime     time.Time   `json:"expiration_date"`
	Installments           int         `json:"installments"`
	AcquirerCode           string      `json:"acquirer_code"`
	Stan                   int         `json:"stan"`
	Transaction            transaction `json:"transaction"`
	Card                   card        `json:"card"`
	Billing                billing     `json:"billing"`
	Settlement             settlement  `json:"settlement"`
}

type transaction struct {
	Amount         float64 `json:"amount"`
	TotalAmount    float64 `json:"total_amount"`
	DecimalDigits  int     `json:"decimal_digits"`
	Currency       string  `json:"currency"`
	CapturedAmount float64 `json:"captured_amount"`
}

type billing struct {
	Amount        float64    `json:"amount"`
	DecimalDigits int        `json:"decimal_digits"`
	Currency      string     `json:"currency"`
	TotalAmount   float64    `json:"total_amount"`
	Conversion    conversion `json:"conversion"`
}

type settlement struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	DecimalDigits int     `json:"decimal_digits"`
}

type card struct {
	NumberID      string        `json:"number_id"`
	Country       string        `json:"country"`
	IssuerAccount issuerAccount `json:"issuer_account"`
	TokenInfo     tokenInfo     `json:"token_info"`
}

type cardAcceptor struct {
	Terminal      string `json:"terminal"`
	Name          string `json:"name"`
	EntryMode     string `json:"entry_mode"`
	Type          string `json:"type"`
	TerminalMode  string `json:"terminal_mode"`
	PinCapability string `json:"pin_capability"`
}

type provider struct {
	ID string `json:"id"`
}

type conversion struct {
	Date          *time.Time `json:"date"`
	DecimalDigits *int       `json:"decimal_digits"`
	Rate          *float64   `json:"rate"`
	From          *string    `json:"from"`
}

type issuerAccount struct {
	ID   string            `json:"id"`
	Card cardIssuerAccount `json:"card"`
}

type cardIssuerAccount struct {
	BusinessMode string `json:"business_mode"`
}

type tokenInfo struct {
	TokenWalletID string `json:"token_wallet_id"`
}
