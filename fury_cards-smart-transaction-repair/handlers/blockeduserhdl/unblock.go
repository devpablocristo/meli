package blockeduserhdl

import (
	"fmt"
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/handlers"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_go-core/pkg/web"
)

const paramIdentificationID = "kycIdentificationID"

func (b BlockedUserHandler) UnblockUser(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	kycIdentificationID, _ := web.Params(req).String(paramIdentificationID)

	wrapFields := field.NewWrappedFields()
	wrapFields.Fields.Add(paramIdentificationID, kycIdentificationID)

	err := b.service.UnblockUser(ctx, kycIdentificationID)
	if err != nil {
		return b.errAPI.CreateAPIErrorAndLog(ctx, err, wrapFields)
	}

	message := fmt.Sprintf("user %s successfully unlocked", kycIdentificationID)
	return web.EncodeJSON(w, handlers.SimpleMessageResponse{Message: message}, http.StatusOK)
}
