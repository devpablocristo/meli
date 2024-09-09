package score

import "time"

type requestScore struct {
	Environment        string    `json:"environment"`
	SiteID             string    `json:"site_id"`
	StatusDetail       string    `json:"status_detail"`
	PayerID            int64     `json:"payer_id"`
	Type               string    `json:"type"`
	IsRecurringPayment bool      `json:"is_recurring_payment"`
	Operation          operation `json:"operation"`
	Provider           provider  `json:"provider"`
}

type operation struct {
	SubType                string        `json:"sub_type"`
	CreationDatetime       time.Time     `json:"creation_datetime"`
	TransmissionDatetime   time.Time     `json:"transmission_datetime"`
	IsAdvice               bool          `json:"is_advice"`
	Installments           int           `json:"installments"`
	AcquirerCode           string        `json:"acquirer_code"`
	Stan                   int           `json:"stan"`
	PaymentAge             int           `json:"payment_age"`
	AccountType            string        `json:"account_type"`
	AuthorizationIndicator string        `json:"authorization_indicator"`
	IsInternational        bool          `json:"is_international"`
	Authorization          authorization `json:"authorization"`
	Transaction            transaction   `json:"transaction"`
	Card                   card          `json:"card"`
	Billing                billing       `json:"billing"`
	Settlement             settlement    `json:"settlement"`
}

type provider struct {
	ID string `json:"id"`
}

type responseScore struct {
	Result string `json:"result"`
	Score  string `json:"score"`
}

type authorization struct {
	ID                   string       `json:"id"`
	Stan                 int          `json:"stan"`
	TransmissionDatetime time.Time    `json:"transmission_datetime"`
	CardAcceptor         cardAcceptor `json:"card_acceptor"`
}

type cardAcceptor struct {
	Terminal      string `json:"terminal"`
	Name          string `json:"name"`
	EntryMode     string `json:"entry_mode"`
	Type          string `json:"type"`
	TerminalMode  string `json:"terminal_mode"`
	PinCapability string `json:"pin_capability"`
}

type transaction struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	DecimalDigits int     `json:"decimal_digits"`
	TotalAmount   float64 `json:"total_amount"`
}

type billing struct {
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	DecimalDigits int        `json:"decimal_digits"`
	TotalAmount   float64    `json:"total_amount"`
	Conversion    conversion `json:"conversion"`
}

type settlement struct {
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	DecimalDigits int        `json:"decimal_digits"`
	TotalAmount   float64    `json:"total_amount"`
	Conversion    conversion `json:"conversion"`
}

type card struct {
	NumberId      string        `json:"number_id"`
	Country       string        `json:"country"`
	IssuerAccount issuerAccount `json:"issuer_account"`
	TokenInfo     tokenInfo     `json:"token_info"`
}

type tokenInfo struct {
	TokenWalletID string `json:"token_wallet_id"`
}

type issuerAccount struct {
	BusinessMode string `json:"business_mode"`
}

type conversion struct {
	Date          *time.Time `json:"date"`
	DecimalDigits *int       `json:"decimal_digits"`
	Rate          *float64   `json:"rate"`
	From          *string    `json:"from"`
}
