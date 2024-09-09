package validationresult

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	utils "github.com/melisource/fury_cards-go-toolkit/pkg/furyutils/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

const (
	keyfields, keyprocessingTimes = "fields", "processing_times"
)

func TestValidationResultConsumer_ConsumeValidationResult_Success(t *testing.T) {
	reader := bytes.NewReader(_validationResultMsgs)
	request := httptest.NewRequest(http.MethodPost, "/cards/smart-transaction-repair/v1/consumer/validations", reader)
	request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
	writer := httptest.NewRecorder()

	service := mockValidationService{
		SaveStub: func(ctx context.Context, inputResult storage.ValidationResult, wrapFields *field.WrappedFields) error {

			var validationResultExpected storage.ValidationResult
			err := json.Unmarshal(_validationResultDomain, &validationResultExpected)
			assert.NoError(t, err)

			for k, v := range inputResult.Reason {
				assert.Equal(t, validationResultExpected.Reason[k].Actual, inputResult.Reason[k].Actual)

				switch k {
				case domain.RuleQtyReparationPerPeriodDays:
					_, cast := v.Accepted.(domain.ReasonResultPerPeriod)
					assert.True(t, cast)
					assert.Equal(t, &domain.ReasonResult{
						Actual:   float64(10),
						Accepted: domain.ReasonResultPerPeriod{Qty: 1, PeriodDays: 100}},
						inputResult.Reason[k])

				case domain.RuleQtyChargebackPerPeriodDays:
					_, cast := v.Accepted.(domain.ReasonResultPerPeriod)
					assert.True(t, cast)
					assert.Equal(t, &domain.ReasonResult{
						Actual:   float64(11),
						Accepted: domain.ReasonResultPerPeriod{Qty: 0, PeriodDays: 111}},
						inputResult.Reason[k])

				default:
					assert.Equal(t, validationResultExpected.Reason[k].Accepted, inputResult.Reason[k].Accepted)
				}
			}

			validationResultExpected.Reason = nil
			inputResult.Reason = nil

			assert.Equal(t, validationResultExpected, inputResult)

			return nil
		},
	}

	log := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {
			for _, f := range fields {
				assert.Contains(t,
					[]string{keyprocessingTimes, keyfields},
					f.Key,
				)
				assertFieldsFromWrap(t, f)
			}
		},
	}

	buildRoutes(writer, request, service, log)

	resp := writer.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestValidationResultConsumer_ConsumeValidationResult_ErrorService(t *testing.T) {
	reader := bytes.NewReader(_validationResultMsgs)
	request := httptest.NewRequest(http.MethodPost, "/cards/smart-transaction-repair/v1/consumer/validations", reader)
	request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
	writer := httptest.NewRecorder()

	service := mockValidationService{
		SaveStub: func(ctx context.Context, inputResult storage.ValidationResult, wrapFields *field.WrappedFields) error {
			return domain.BuildErrorBadGateway(errors.New("any error"), "timeout")
		},
	}

	log := mockLogService{
		ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
			for _, f := range fields {
				assert.Contains(t,
					[]string{keyprocessingTimes, keyfields, "type", "error", "stack"},
					f.Key,
				)
				assertFieldsFromWrap(t, f)
			}
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
		SetStacktub: func(stacktrace string) log.Field {
			return log.Field{Key: "stack", Interface: stacktrace}
		},
	}

	buildRoutes(writer, request, service, log)

	resp := writer.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
}

func TestValidationResultConsumer_ConsumeValidationResult_ErrorBadRequest(t *testing.T) {
	reader := bytes.NewReader(nil)
	request := httptest.NewRequest(http.MethodPost, "/cards/smart-transaction-repair/v1/consumer/validations", reader)
	request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
	writer := httptest.NewRecorder()

	service := mockValidationService{}

	log := mockLogService{
		ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
		SetStacktub: func(stacktrace string) log.Field {
			return log.Field{Key: "stack", Interface: stacktrace}
		},
	}

	buildRoutes(writer, request, service, log)

	resp := writer.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
