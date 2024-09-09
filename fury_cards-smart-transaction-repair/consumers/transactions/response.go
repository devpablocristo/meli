package transactions

import (
	"context"
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (t TransactionsConsumer) buildBulkResponse(
	ctx context.Context,
	consumerRequest *msgConsumerRequest,
	output *domain.TransactionsBulkOutput,
	transactionType transactionType,
	siteID string,
	wrapFields *field.WrappedFields,
	gnierr error,
) bulkResponse {

	bulk := bulkResponse{}

	if gnierr != nil {
		wrapFields.Fields.Add(keyFieldResultStatus, "full_error")
		t.reBuildErrorAndLog(ctx, gnierr, transactionType, siteID, wrapFields)
		bulk.Responses = fillAllErrorResponseMessages(consumerRequest, gnierr)
		return bulk
	}

	msgsResponse, authorizationsIDsSuccess := fillResponseMessages(consumerRequest, *output)
	bulk.Responses = msgsResponse

	t.logResult(ctx, *output, authorizationsIDsSuccess, transactionType, siteID, wrapFields)
	return bulk
}

func fillResponseMessages(
	consumerRequest *msgConsumerRequest,
	output domain.TransactionsBulkOutput,
) (msgsResponse []msgResponse, authorizationsIDsSuccess []string) {

	msgsResponse = []msgResponse{}
	for _, msg := range consumerRequest.Messages {
		if err, found := output.AuthorizationIDsWithErr[msg.Body.AuthorizationID]; found {
			msgsResponse = append(msgsResponse, msgResponse{ID: msg.ID, Code: setCode(err)})
			continue
		}
		msgsResponse = append(msgsResponse, msgResponse{ID: msg.ID, Code: http.StatusOK})
		authorizationsIDsSuccess = append(authorizationsIDsSuccess, msg.Body.AuthorizationID)
	}

	return
}

func fillAllErrorResponseMessages(consumerRequest *msgConsumerRequest, gnierr error) []msgResponse {
	msgsResponse := []msgResponse{}

	for _, msg := range consumerRequest.Messages {
		msgsResponse = append(msgsResponse, msgResponse{ID: msg.ID, Code: setCode(gnierr)})
	}

	return msgsResponse
}

func setCode(gnierr error) int {
	errType := gnierrors.Type(gnierr)

	switch errType {
	case domain.BadGateway:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
