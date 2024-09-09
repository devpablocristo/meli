package cardstransactions

import "time"

type requestReversal struct {
	ID        string    `json:"id"`
	Operation operation `json:"operation"`
	Provider  provider  `json:"provider"`
	Options   options   `json:"options"`
}

type operation struct {
	Type                 *string       `json:"type"`
	TransmissionDatetime time.Time     `json:"transmission_datetime"`
	IsAdvice             bool          `json:"is_advice"`
	Installments         int           `json:"installments"`
	AcquirerCode         *string       `json:"acquirer_code"`
	Stan                 *int          `json:"stan"`
	Authorization        authorization `json:"authorization"`
	Transaction          transaction   `json:"transaction"`
	Settlement           settlement    `json:"settlement"`
	Billing              billing       `json:"billing"`
	Card                 card          `json:"card"`
}

type provider struct {
	ID string `json:"id"`
}

type options struct {
	ReversalType string `json:"reversal_type"`
}

type authorization struct {
	ID                   string       `json:"id"`
	Stan                 *int         `json:"stan"`
	TransmissionDatetime time.Time    `json:"transmission_datetime"`
	CardAcceptor         cardAcceptor `json:"card_acceptor"`
}

type cardAcceptor struct {
	Terminal *string `json:"terminal"`
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
	Conversion    conversion `json:"conversion"`
}

type card struct {
	NumberId string `json:"number_id"`
	Country  string `json:"country"`
}

type conversion struct {
	Date          *time.Time `json:"date"`
	DecimalDigits *int       `json:"decimal_digits"`
	Rate          *float64   `json:"rate"`
	From          *string    `json:"from"`
}

type responseReversal struct {
	Result responseReversalResult `json:"result"`
}

type responseReversalResult struct {
	Status          string  `json:"status"`
	StatusDetail    *string `json:"status_detail"`
	ProcessedAt     string  `json:"processed_at"`
	UnblockCard     *string `json:"unblock_card"`
	PartialApproval *string `json:"partial_approval"`
	Balance         *string `json:"balance"`
	Operation       *string `json:"operation"`
}
