package domain

type SearchHubCountChargebackInput struct {
	UserID   string
	LastDays int
}

type SearchHubCountChargebackOutput struct {
	Total int
}
