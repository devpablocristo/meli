package validationforcepublishhdl

import (
	"log"
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_go-core/pkg/web"
)

func (h ValidationForcePublishHandler) Publish(w http.ResponseWriter, r *http.Request) error {
	var reparationReq msgValidationPublish
	wrapFields := field.NewWrappedFields()

	if err := web.DecodeJSON(r, &reparationReq); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return err
	}

	validationResult := storage.ValidationResult{
		PaymentID:           reparationReq.PaymentID,
		KycIdentificationID: reparationReq.KycIdentitificationID,
		UserID:              reparationReq.UserID,
		IsEligible:          reparationReq.IsEligible,
		DateCreated:         reparationReq.CreatedAt,
		AuthorizationID:     reparationReq.AuthorizationID,
		Type:                reparationReq.Type,
		SiteID:              reparationReq.SiteID,
		FaqID:               reparationReq.FaqID,
		Requester: domain.Requester{
			ClientApplication: reparationReq.ClientApplication,
			ClientScope:       reparationReq.ClientScope,
		},
		TransactionEvaluated: storage.TransactionEvaluated{
			Amount:   reparationReq.Amount,
			Currency: reparationReq.Currency,
		},
		Reason: fillReason(reparationReq.Reason),
	}

	h.service.PublishEvent(r.Context(), &validationResult, wrapFields)

	return web.EncodeJSON(w, reparationReq, http.StatusOK)
}

func fillReason(requestReason map[domain.RuleType]publicDefault) domain.Reason {
	reason := domain.Reason{}

	for k, v := range requestReason {
		reason[k] = v.getValue()
	}

	return reason
}
