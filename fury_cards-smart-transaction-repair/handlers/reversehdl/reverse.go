package reversehdl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/handlers"
	utils "github.com/melisource/fury_cards-go-toolkit/pkg/furyutils/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/metrics/v1"
	"github.com/melisource/fury_go-core/pkg/web"
)

const (
	headerKeyClientID          = "x-client-id"
	headerKeyClientApplication = "X-Api-Client-Application"
	headerKeyClientScope       = "X-Api-Client-Scope"
)

func (r ReverseHandler) Reverse(w http.ResponseWriter, req *http.Request) error {
	ctx, span := metrics.StartSpan(req.Context(), fmt.Sprint(domain.MetricRoot, "handler.reverse"))
	defer span.Finish()

	wrapFields := field.NewWrappedFields()
	paymentID, _ := web.Params(req).String("id")
	clientApplication := req.Header.Get(headerKeyClientApplication)

	input, err := r.validateAndFillInput(ctx, req, paymentID, wrapFields, clientApplication)
	if err != nil {
		return err
	}

	r.setLogFieldsRequest(input.PaymentID, input.SiteID, input.FaqID, input.Requester.ClientApplication, input.UserID, wrapFields)

	err = r.reverse(ctx, *input, wrapFields, clientApplication)
	if err != nil {
		return r.errAPI.CreateAPIErrorAndLog(ctx, err, wrapFields)
	}

	r.log.Info(ctx, "reverse_successfully_requested", wrapFields.ToLogField(log.InfoLevel)...)

	message := "Reverse successfully requested"
	return web.EncodeJSON(w, handlers.SimpleMessageResponse{Message: message}, http.StatusOK)
}

func (r ReverseHandler) reverse(
	ctx context.Context,
	input domain.ReverseTransactionInput,
	wrapFields *field.WrappedFields,
	clientApplication string,
) error {
	start := time.Now()
	startTimer := wrapFields.Timers.Start(domain.TimerHandlerReverse)
	defer func() {
		addMetricTimerHandler(ctx, time.Since(start), input.SiteID, clientApplication)
		startTimer.Stop()
	}()

	return r.service.Reverse(ctx, input, wrapFields)
}

func (r ReverseHandler) validateAndFillInput(
	ctx context.Context,
	req *http.Request,
	paymentID string,
	wrapFields *field.WrappedFields,
	clientApplication string,
) (*domain.ReverseTransactionInput, error) {
	isSourceFury := utils.IsFuryRequest(req)
	if !isSourceFury {
		// try to get as much information as possible in the log.
		reparationReq, _ := r.getAndValidateRequest(ctx, req, clientApplication)

		addMetricUnauthorized(ctx, reparationReq.SiteID, clientApplication)

		r.setLogFieldsRequest(paymentID, reparationReq.SiteID, reparationReq.FaqID, clientApplication, reparationReq.UserID, wrapFields)
		err := r.errAPI.BuildErrorUnauthorizedAPI(msgError)
		return nil, r.errAPI.CreateAPIErrorAndLog(ctx, err, wrapFields)
	}

	reparationReq, err := r.getAndValidateRequest(ctx, req, clientApplication)
	if err != nil {
		r.setLogFieldsRequest(paymentID, reparationReq.SiteID, reparationReq.FaqID, clientApplication, reparationReq.UserID, wrapFields)
		return nil, r.errAPI.CreateAPIErrorAndLog(ctx, err, wrapFields)
	}

	headerValueXClientID, err := r.getAndValidateHeader(ctx, req, headerKeyClientID, reparationReq.SiteID, clientApplication)
	if err != nil {
		r.setLogFieldsRequest(paymentID, reparationReq.SiteID, reparationReq.FaqID, clientApplication, reparationReq.UserID, wrapFields)
		return nil, r.errAPI.CreateAPIErrorAndLog(ctx, err, wrapFields)
	}

	return &domain.ReverseTransactionInput{
		PaymentID:            paymentID,
		UserID:               reparationReq.UserID,
		HeaderValueXClientID: headerValueXClientID,
		SiteID:               strings.ToUpper(reparationReq.SiteID),
		FaqID:                reparationReq.FaqID,
		Requester: domain.Requester{
			ClientApplication: clientApplication,
			ClientScope:       req.Header.Get(headerKeyClientScope),
		},
	}, nil
}

func (r ReverseHandler) getAndValidateRequest(ctx context.Context, req *http.Request, clientApplication string) (reparationRequest, error) {
	var reparationReq reparationRequest

	if err := web.DecodeJSON(req, &reparationReq); err != nil {
		addMetricIvalidPayload(ctx, clientApplication)
		return reparationReq, r.buildErrorBadRequest(err)
	}

	return reparationReq, nil
}

func (r ReverseHandler) getAndValidateHeader(
	ctx context.Context,
	req *http.Request,
	headerKey, siteID, clientApplication string,
) (string, error) {
	valueHeader := req.Header.Get(headerKey)

	if strings.TrimSpace(valueHeader) == "" {
		addMetricHeaderRequired(ctx, siteID, clientApplication)
		return "", r.buildErrorBadRequest(fmt.Errorf(msgHeaderRequired, headerKey))
	}
	return valueHeader, nil
}

func (r ReverseHandler) setLogFieldsRequest(
	paymentID,
	siteID,
	faqID,
	clientApplication string,
	userID int64,
	wrapFields *field.WrappedFields,
) {
	logFieldsRequest := logFieldsRequest{
		SiteID:            siteID,
		PaymentID:         paymentID,
		UserID:            userID,
		FaqID:             faqID,
		ClientApplication: clientApplication,
	}

	b, _ := json.Marshal(logFieldsRequest)
	parameters := string(b)

	wrapFields.Fields.Add("input_parameters", parameters)
}
