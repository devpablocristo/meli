package transactions

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
	"github.com/melisource/fury_go-core/pkg/web"
)

const transactionCaptureType transactionType = "capture"

func (t TransactionsConsumer) ConsumeCaptures(w http.ResponseWriter, r *http.Request) error {
	wrapFields := field.NewWrappedFields()
	startTimer := wrapFields.Timers.Start(string(domain.TimerConsumerCaptures))
	start := time.Now()

	siteID := getSiteID(r)

	ctx, span := metrics.StartSpan(r.Context(), fmt.Sprint(domain.MetricRoot, "consumer_capture_transactions_", siteID))
	defer func() {
		span.Finish()
		addMetricTimerConsumerCaptures(ctx, siteID, time.Since(start))
	}()

	consumerRequest, err := t.getRequest(ctx, r, wrapFields, transactionCaptureType, siteID)
	if err != nil {
		return err
	}

	replaceAuthorizationID(consumerRequest)

	input := t.buildTransactionsInput(*consumerRequest, wrapFields, siteID, transactionCaptureType)

	output, err := t.service.BlockUserIfReparationExistsInBulk(ctx, *input, wrapFields)
	startTimer.Stop()

	bulkResponse := t.buildBulkResponse(ctx, consumerRequest, output, transactionCaptureType, input.SiteID, wrapFields, err)

	return web.EncodeJSON(w, bulkResponse, http.StatusOK)
}

// If a capture has filled its original_transaction_id, you must replace the authorization_id with this.
func replaceAuthorizationID(consumerRequest *msgConsumerRequest) {
	for i, msg := range consumerRequest.Messages {
		if len(strings.TrimSpace(msg.Body.OriginalTransactionID)) > 0 {
			consumerRequest.Messages[i].Body.AuthorizationID = msg.Body.OriginalTransactionID
		}
	}
}
