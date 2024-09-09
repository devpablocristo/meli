package validationresult

import (
	"fmt"
	"net/http"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
	"github.com/melisource/fury_go-core/pkg/web"
)

func (v ValidationResultConsumer) ConsumeValidationResult(w http.ResponseWriter, r *http.Request) error {
	wrapFields := field.NewWrappedFields()
	startTimer := wrapFields.Timers.Start(string(domain.TimerConsumerValidationResult))
	start := time.Now()

	consumerRequest, err := v.getRequest(r.Context(), r, wrapFields)
	if err != nil {
		return err
	}

	addMetricValidationResultReceived(r.Context(), consumerRequest.SiteID)

	ctx, span := metrics.StartSpan(r.Context(), fmt.Sprint(domain.MetricRoot, consumerName))
	defer func() {
		span.Finish()
		addMetricTimerConsumerValidationResult(ctx, consumerRequest.SiteID, time.Since(start))
	}()

	input := v.buildValidationResultInput(*consumerRequest, wrapFields)

	err = v.service.Save(ctx, input, wrapFields)
	startTimer.Stop()
	if err != nil {
		return v.reBuildErrorAndLog(ctx, err, consumerRequest.SiteID, wrapFields)
	}

	addMetricValidationResultDone(ctx, consumerRequest.SiteID)

	v.log.Info(ctx, fmt.Sprintf("%s_consumed_validation_result_successfully", consumerRequest.SiteID), wrapFields.ToLogField(log.InfoLevel)...)

	return web.EncodeJSON(w, simpleMessageResponse{Message: "consumed validation result"}, http.StatusOK)
}
