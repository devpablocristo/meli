package reversehdl

type logFieldsRequest struct {
	SiteID            string `json:"site_id"`
	PaymentID         string `json:"payment_id"`
	UserID            int64  `json:"user_id"`
	FaqID             string `json:"faq_id,omitempty"`
	ClientApplication string `json:"client_application"`
}
