package transactions

type bulkResponse struct {
	Responses []msgResponse `json:"responses"`
}

type msgResponse struct {
	ID   string `json:"id"`
	Code int    `json:"code"`
}

type msgConsumerRequest struct {
	Messages []msgBody `json:"messages"`
}

type msgBody struct {
	ID   string         `json:"id"`
	Body msgTransaction `json:"msg"`
}

type msgTransaction struct {
	ID                    string `json:"id"`
	AuthorizationID       string `json:"authorization_id"`
	Type                  string `json:"type"`
	OriginalTransactionID string `json:"original_authorization_id"`
}
