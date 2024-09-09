package reversehdl

type reparationRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	SiteID string `json:"site_id" validate:"required"`
	FaqID  string `json:"faq_id"`
}
