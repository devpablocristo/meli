package domain

type TransactionsNewsFeedBulkInput struct {
	AuthorizationsNewsFeed map[string]TransactionIDNewsFeed
	SiteID                 string
}

type TransactionsBulkOutput struct {
	AuthorizationIDsWithErr map[string]error
}

type TransactionIDNewsFeed string

type BlockedUserInput struct {
	KycIdentificationID string
	SiteID              string
}
