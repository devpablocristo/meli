package transactionhistory

type History struct {
	Wallet *wallet `json:"wallet"`
}

type wallet struct {
	Cards []card `json:"cards"`
}

type card struct {
	IssuerAccounts []issuerAccount `json:"issuer_accounts"`
}

type issuerAccount struct {
	ID           string              `json:"id"`
	Transactions transactionAccounts `json:"transactions"`
}

type transactionAccounts struct {
	MatchTotal   int                    `json:"match_total"`
	Transactions []transactionOperation `json:"transactions"`
}

type transactionOperation struct {
	ID        string    `json:"id"`
	Operation operation `json:"operation"`
}

type operation struct {
	Billing    billing    `json:"billing"`
	Settlement settlement `json:"settlement"`
}

type billing struct {
	Amount        float64 `json:"amount"`
	DecimalDigits int     `json:"decimal_digits"`
	Currency      string  `json:"currency"`
}

type settlement struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	DecimalDigits int     `json:"decimal_digits"`
}
