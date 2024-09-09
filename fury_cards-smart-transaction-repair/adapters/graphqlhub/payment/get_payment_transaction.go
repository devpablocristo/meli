package payment

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchhubclient/v1"
	"github.com/melisource/fury_go-core/pkg/log"
)

func GetPaymentTransaction(
	ctx context.Context,
	searchHubCli searchhubclient.Service,
	transactionInput domain.TransactionInput,
	wrapFields *field.WrappedFields,
) (domain.TransactionOutput, error) {

	request := BuildRequestPaymentTransaction(transactionInput)

	var data paymentTransaction

	err := searchHubCli.SetDataByQueryRequest(ctx, request, &data)
	if err != nil {
		return domain.TransactionOutput{}, buildErrorBadGateway(err, msgErrFailGetPaymentTransaction)
	}

	body, _ := json.Marshal(data)
	wrapFields.Fields.AddAt(string(domain.KeyFieldBodyPaymentTransaction), strings.ReplaceAll(string(body), "\\", ""), log.ErrorLevel)

	return FillOutput(data), nil
}
