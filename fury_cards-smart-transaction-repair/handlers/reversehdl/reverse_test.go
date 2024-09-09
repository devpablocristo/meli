package reversehdl

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	utils "github.com/melisource/fury_cards-go-toolkit/pkg/furyutils/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

const (
	basePath                                     = "/cards/smart-transaction-repair/v1/reverse"
	keytype, keyerror, keyinfo                   = "type", "error", "detail"
	keyfields, keyprocessingtimes, keystacktrace = "fields", "processing_times", "stacktrace"
)

func Test_Reverse(t *testing.T) {
	t.Run("should create a reverse with success", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodPost, basePath+"/100", mockBody(t, reparationRequest{UserID: 123, SiteID: "MLM", FaqID: "FaqID"}))

		request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
		request.Header.Add("x-client-id", "123456")
		request.Header.Add("X-Api-Client-Application", "app-test")
		request.Header.Add("X-Api-Client-Scope", "beta")
		writer := httptest.NewRecorder()

		service := &mockReparationService{
			ReverseStub: func(ctx context.Context, reverseInput domain.ReverseTransactionInput, wrapFields *field.WrappedFields) error {
				utilstest.AssertFieldsNoEmptyFromStruct(t, reverseInput)
				return nil
			},
		}

		mockLogService := mockLogService{
			InfoStub: func(c context.Context, msg string, fields ...log.Field) {
				logFields := []string{keyprocessingtimes, keyfields}
				keysLog := []string{}
				for _, f := range fields {
					assert.Contains(t, logFields, f.Key)
					assertFieldsFromWrap(t, f)
					keysLog = append(keysLog, f.Key)
				}
				for _, lf := range logFields {
					assert.Contains(t, keysLog, lf)
				}
			},
		}

		buildRoutes(writer, request, service, mockLogService)

		resp := writer.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_Reverse_Error_Conflict(t *testing.T) {
	t.Run("should return status code 409 if repair already done", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, basePath+"/99", mockBody(t, reparationRequest{UserID: 123, SiteID: "MLM"}))
		request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
		request.Header.Add("x-client-id", "123456")
		writer := httptest.NewRecorder()

		service := &mockReparationService{
			ReverseStub: func(ctx context.Context, reversal domain.ReverseTransactionInput, wrapFields *field.WrappedFields) error {
				return domain.BuildErrorAlreadyRepaired(errors.New("reverse already done"), "invalid request")
			},
		}

		mockLogService := mockLogService{
			InfolnStub: func(c context.Context, msg string, fields ...log.Field) {
				errExpected := gnierrors.Wrap(
					domain.AlreadyRepaired, errors.New(`"reverse already done"`), "invalid request", false)

				assert.Equal(t, "invalid request", msg)
				assert.Equal(t, keytype, fields[0].Key)
				assert.Equal(t, keyinfo, fields[1].Key)
				assert.Equal(t, gnierrors.Cause(errExpected).Error(), fields[1].Interface.(string))

				logFields := []string{keytype, keyprocessingtimes, keyfields, keyinfo}
				keysLog := []string{}
				for _, f := range fields {
					assert.Contains(t, logFields, f.Key)
					keysLog = append(keysLog, f.Key)
					assertFieldsFromWrap(t, f)
				}
				for _, lf := range logFields {
					assert.Contains(t, keysLog, lf)
				}

			},
			AnyStub: func(key string, value interface{}) log.Field {
				return log.Field{Key: key, Interface: value}
			},
		}

		buildRoutes(writer, request, service, mockLogService)

		resp := writer.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}

func Test_Reverse_Error_Unauthorized(t *testing.T) {
	t.Run("should return status code 401 if request is manual (without header fury)", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, basePath+"/100", mockBody(t, reparationRequest{UserID: 123, SiteID: "MLM"}))
		writer := httptest.NewRecorder()

		service := &mockReparationService{}

		mockLogService := mockLogService{
			ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
				errExpected := gnierrors.Wrap(
					domain.Unauthorized, errors.New("request is not authorized"), "invalid request", false)

				assert.Equal(t, "invalid request", msg)
				assert.Equal(t, keytype, fields[0].Key)
				assert.Equal(t, keyerror, fields[1].Key)
				assert.Equal(t, gnierrors.Cause(errExpected).Error(), fields[1].Interface.(string))

				var logFields = []string{keytype, "unauthorized_request", keyerror, keyfields, keystacktrace}
				keysLog := []string{}
				for _, f := range fields {
					assert.Contains(t, logFields, f.Key)
					keysLog = append(keysLog, f.Key)
					assertFieldsFromWrap(t, f)
				}
				for _, lf := range logFields {
					assert.Contains(t, keysLog, lf)
				}
			},
			SetStacktub: func(stacktrace string) log.Field {
				return log.Field{Key: keystacktrace, String: stacktrace}
			},
			AnyStub: func(key string, value interface{}) log.Field {
				switch key {
				case "unauthorized_request":
					assert.Equal(t, value, "request is not authorized: manual request not allowed")
				case keytype:
					assert.Equal(t, "unauthorized", value)
				case keyerror:
				default:
					panic(key)
				}
				return log.Field{Key: key, Interface: value}
			},
		}

		buildRoutes(writer, request, service, mockLogService)

		resp := writer.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func Test_Reverse_Error_Body_UnprocessableEntity(t *testing.T) {
	t.Run("should return status code 422", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, basePath+"/99", mockBody(t, reparationRequest{UserID: 0, SiteID: "MLM"}))
		request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
		writer := httptest.NewRecorder()

		service := &mockReparationService{}

		mockLogService := mockLogService{
			ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
				errExpected := gnierrors.Wrap(
					domain.UnprocessableEntity,
					errors.New("unprocessable_entity: validation_error: invalid fields: UserID"),
					"invalid request",
					false)

				assert.Equal(t, "invalid request", msg)
				assert.Equal(t, keytype, fields[0].Key)
				assert.Equal(t, keyerror, fields[1].Key)
				assert.Equal(t, gnierrors.Cause(errExpected).Error(), fields[1].Interface.(string))

				var logFields = []string{keytype, keyerror, keyfields, keystacktrace}
				keysLog := []string{}
				for _, f := range fields {
					assert.Contains(t, logFields, f.Key)
					keysLog = append(keysLog, f.Key)
					assertFieldsFromWrap(t, f)
				}
				for _, lf := range logFields {
					assert.Contains(t, keysLog, lf)
				}
			},
			SetStacktub: func(stacktrace string) log.Field {
				return log.Field{Key: keystacktrace, String: stacktrace}
			},
			AnyStub: func(key string, value interface{}) log.Field {
				switch key {
				case keytype:
					assert.Equal(t, "bad_request", value)
				case keyerror:
				default:
					panic(key)
				}
				return log.Field{Key: key, Interface: value}
			},
		}

		buildRoutes(writer, request, service, mockLogService)

		resp := writer.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func Test_Reverse_Error_Header_Required(t *testing.T) {
	t.Run("should return status code 400", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, basePath+"/100", mockBody(t, reparationRequest{UserID: 123, SiteID: "MLM"}))
		request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
		writer := httptest.NewRecorder()

		service := &mockReparationService{}

		mockLogService := mockLogService{
			ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
				errExpected := gnierrors.Wrap(
					domain.BadRequest, errors.New("x-client-id header is required"), "invalid request", false)

				assert.Equal(t, "invalid request", msg)
				assert.Equal(t, keytype, fields[0].Key)
				assert.Equal(t, keyerror, fields[1].Key)
				assert.Equal(t, gnierrors.Cause(errExpected).Error(), fields[1].Interface.(string))

				var logFields = []string{keytype, keyerror, keyfields, keystacktrace}
				keysLog := []string{}
				for _, f := range fields {
					assert.Contains(t, logFields, f.Key)
					keysLog = append(keysLog, f.Key)
					assertFieldsFromWrap(t, f)
				}
				for _, lf := range logFields {
					assert.Contains(t, keysLog, lf)
				}
			},
			SetStacktub: func(stacktrace string) log.Field {
				return log.Field{Key: keystacktrace, String: stacktrace}
			},
			AnyStub: func(key string, value interface{}) log.Field {
				switch key {
				case keytype:
					assert.Equal(t, "bad_request", value)
				case keyerror:
				default:
					panic(key)
				}
				return log.Field{Key: key, Interface: value}
			},
		}

		buildRoutes(writer, request, service, mockLogService)

		resp := writer.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestReverseHandler_getAndValidateRequest_validate_request(t *testing.T) {
	ctx := context.Background()
	r := ReverseHandler{}

	request := httptest.NewRequest(
		http.MethodPost, basePath+"/100", mockBody(t, reparationRequest{UserID: 123, SiteID: "MLM", FaqID: "faq_id"}))

	repRequest, err := r.getAndValidateRequest(ctx, request, "app_test")

	assert.NoError(t, err)
	utilstest.AssertFieldsNoEmptyFromStruct(t, repRequest)
}
