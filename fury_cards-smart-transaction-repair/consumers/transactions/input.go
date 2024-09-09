package transactions

import (
	"context"
	"net/http"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_go-core/pkg/web"
)

func (t TransactionsConsumer) getRequest(
	ctx context.Context,
	r *http.Request,
	wrapFields *field.WrappedFields,
	transactionType transactionType,
	siteID string,
) (*msgConsumerRequest, error) {

	var request msgConsumerRequest

	if err := web.DecodeJSON(r, &request); err != nil {
		addMetricIvalidPayload(ctx, siteID, domain.MetricTagSourceQReversalTransactions)
		return nil, t.buildErrorBadRequestAndLog(ctx, err, transactionType, siteID, wrapFields)
	}

	return &request, nil
}

func (t TransactionsConsumer) buildTransactionsInput(
	request msgConsumerRequest,
	wrapFields *field.WrappedFields,
	siteID string,
	transactionType transactionType,
) *domain.TransactionsNewsFeedBulkInput {

	authorizationsIDs := []string{}
	authorizationsNewsFeed := map[string]domain.TransactionIDNewsFeed{}
	for _, msg := range request.Messages {
		authorizationsNewsFeed[msg.Body.AuthorizationID] = domain.TransactionIDNewsFeed(msg.Body.ID)
		authorizationsIDs = append(authorizationsIDs, msg.Body.AuthorizationID)
	}

	input := &domain.TransactionsNewsFeedBulkInput{
		AuthorizationsNewsFeed: authorizationsNewsFeed,
		SiteID:                 siteID,
	}

	setFieldsInputParameters(wrapFields, siteID, transactionType, len(authorizationsIDs))
	wrapFields.Fields.Add(string(domain.KeyFieldAuthorizationIDs), strings.Join(authorizationsIDs, ", "))

	return input
}

func getSiteID(r *http.Request) string {
	siteID, _ := web.Params(r).String("siteID")
	return strings.ToUpper(siteID)
}
