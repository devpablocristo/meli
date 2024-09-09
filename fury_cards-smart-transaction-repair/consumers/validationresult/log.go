package validationresult

import (
	"encoding/json"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

type logFieldsRequest struct {
	SiteID            string `json:"site_id"`
	PaymentID         string `json:"payment_id"`
	UserID            int64  `json:"user_id"`
	FaqID             string `json:"faq_id,omitempty"`
	ClientApplication string `json:"client_application"`
}

func setFieldsInputParameters(wrapFields *field.WrappedFields, request msgValidationResult) {
	logFieldsRequest := logFieldsRequest{
		SiteID:            request.SiteID,
		PaymentID:         request.PaymentID,
		UserID:            request.UserID,
		FaqID:             request.FaqID,
		ClientApplication: request.ClientApplication,
	}

	b, _ := json.Marshal(logFieldsRequest)
	parameters := string(b)

	wrapFields.Fields.Add(string(domain.KeyFieldInputParameters), parameters)
}
