package transactions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
	"github.com/melisource/fury_go-core/pkg/web"
)

const transactionReversalType transactionType = "reversal"

func (t TransactionsConsumer) ConsumeReversals(w http.ResponseWriter, r *http.Request) error {
	wrapFields := field.NewWrappedFields()
	startTimer := wrapFields.Timers.Start(string(domain.TimerConsumerReversals))
	start := time.Now()

	siteID := getSiteID(r)

	ctx, span := metrics.StartSpan(r.Context(), fmt.Sprint(domain.MetricRoot, "consumer_reversal_transactions_", siteID))
	defer func() {
		span.Finish()
		addMetricTimerConsumerReverses(ctx, siteID, time.Since(start))
	}()

	consumerRequest, err := t.getRequest(ctx, r, wrapFields, transactionReversalType, siteID)
	if err != nil {
		return err
	}

	input := t.buildTransactionsInput(*consumerRequest, wrapFields, siteID, transactionReversalType)

	output, err := t.service.UnblockUserByReversalClearingInBulk(ctx, *input, wrapFields)
	startTimer.Stop()

	bulkResponse := t.buildBulkResponse(ctx, consumerRequest, output, transactionReversalType, input.SiteID, wrapFields, err)

	return web.EncodeJSON(w, bulkResponse, http.StatusOK)
}
