package transactions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	utils "github.com/melisource/fury_cards-go-toolkit/pkg/furyutils/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

const (
	keystacktrace, keyfields, keyprocessingTimes = "stacktrace", "fields", "processing_times"
)

func TestTransactionsConsumer_Success(t *testing.T) {
	tests := []struct {
		name                string
		url                 string
		body                []byte
		fromProcessingTimes string
		service             mockBlockedUserService
	}{
		{
			name:                "consumer capture",
			url:                 "/cards/smart-transaction-repair/v1/consumer/transactions/captures/mlm",
			body:                _captureMsgs,
			fromProcessingTimes: "consumer.captures",
			service: mockBlockedUserService{
				BlockUserIfReparationExistsInBulkStub: func(
					ctx context.Context,
					input domain.TransactionsNewsFeedBulkInput,
					wrapFields *field.WrappedFields,
				) (*domain.TransactionsBulkOutput, error) {

					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("capture_auth_mlb_test_adtorres_score-full-8"),
						input.AuthorizationsNewsFeed["forced_offline-capture_auth_mlb_test_adtorres_score-full-8-1669990748570"],
					)

					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("capture_auth_mlb_test_adtorres_score-full-8_2"),
						input.AuthorizationsNewsFeed["original_auth_02"],
					)

					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("capture_auth_03"),
						input.AuthorizationsNewsFeed["auth_03"],
					)

					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("capture_auth_04"),
						input.AuthorizationsNewsFeed["auth_04"],
					)

					return &domain.TransactionsBulkOutput{}, nil
				},
			},
		},
		{
			name:                "consumer reverse",
			url:                 "/cards/smart-transaction-repair/v1/consumer/transactions/reversals/mlm",
			body:                _reverseMsgs,
			fromProcessingTimes: "consumer.reversals",
			service: mockBlockedUserService{
				UnblockUserByReversalClearingInBulkStub: func(
					ctx context.Context,
					input domain.TransactionsNewsFeedBulkInput,
					wrapFields *field.WrappedFields,
				) (*domain.TransactionsBulkOutput, error) {
					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("reverse_auth_01"),
						input.AuthorizationsNewsFeed["auth_01"],
					)

					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("reverse_auth_02"),
						input.AuthorizationsNewsFeed["auth_02"],
					)

					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("reverse_auth_03"),
						input.AuthorizationsNewsFeed["auth_03"],
					)

					assert.Equal(
						t,
						domain.TransactionIDNewsFeed("reverse_auth_04"),
						input.AuthorizationsNewsFeed["auth_04"],
					)

					return &domain.TransactionsBulkOutput{}, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.url, reader)
			request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
			writer := httptest.NewRecorder()

			log := mockLogService{
				InfolnStub: func(c context.Context, msg string, fields ...log.Field) {
					for _, f := range fields {
						assert.Contains(t,
							[]string{keyprocessingTimes, keyfields, keyFieldAuthorizationIDsWithSuccessful, keyFieldResultStatus},
							f.Key,
						)
						assertFieldsFromWrap(t, f, tt.fromProcessingTimes)
					}
				},
				AnyStub: func(key string, value interface{}) log.Field {
					return log.Field{Key: key, Interface: value}
				},
			}

			buildRoutes(writer, request, tt.service, log)

			resp := writer.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestTransactionsConsumer_Error_AllMessages_Get_Timeout(t *testing.T) {
	tests := []struct {
		name                string
		url                 string
		body                []byte
		fromProcessingTimes string
		msgErr              string
		service             mockBlockedUserService
	}{
		{
			name:                "consumer capture",
			url:                 "/cards/smart-transaction-repair/v1/consumer/transactions/captures/mlm",
			body:                _captureMsgs,
			fromProcessingTimes: "consumer.captures",
			msgErr:              "MLM consumer capture transactions failed",
			service: mockBlockedUserService{
				BlockUserIfReparationExistsInBulkStub: func(
					ctx context.Context,
					input domain.TransactionsNewsFeedBulkInput,
					wrapFields *field.WrappedFields,
				) (*domain.TransactionsBulkOutput, error) {
					return nil,
						gnierrors.Wrap(domain.BadGateway, errors.New("timeout error"), "blockeduser service failure", false)
				},
			},
		},
		{
			name:                "consumer reverse",
			url:                 "/cards/smart-transaction-repair/v1/consumer/transactions/reversals/mlm",
			body:                _reverseMsgs,
			fromProcessingTimes: "consumer.reversals",
			msgErr:              "MLM consumer reversal transactions failed",
			service: mockBlockedUserService{
				UnblockUserByReversalClearingInBulkStub: func(
					ctx context.Context,
					input domain.TransactionsNewsFeedBulkInput,
					wrapFields *field.WrappedFields,
				) (*domain.TransactionsBulkOutput, error) {
					return nil,
						gnierrors.Wrap(domain.BadGateway, errors.New("timeout error"), "blockeduser service failure", false)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.url, reader)
			request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
			writer := httptest.NewRecorder()

			log := mockLogService{
				ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
					assert.Equal(t, tt.msgErr, msg)
					for _, f := range fields {
						assert.Contains(t, []string{"type", "error", keyfields, keyprocessingTimes, keystacktrace}, f.Key)
						assertFieldsFromWrap(t, f, tt.fromProcessingTimes)
					}
				},
				SetStacktub: func(stacktrace string) log.Field {
					return log.Field{Key: keystacktrace, String: stacktrace}
				},
				AnyStub: func(key string, value interface{}) log.Field {
					return log.Field{Key: key, Interface: value}
				},
			}

			buildRoutes(writer, request, tt.service, log)

			resp := writer.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			bodyResponse, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var bulkResponse bulkResponse
			err = json.Unmarshal(bodyResponse, &bulkResponse)
			assert.NoError(t, err)

			for index, item := range bulkResponse.Responses {
				switch index {
				case 0:
					assert.Equal(t, "01", item.ID)
					assert.Equal(t, http.StatusBadGateway, item.Code)
				case 1:
					assert.Equal(t, "02", item.ID)
					assert.Equal(t, http.StatusBadGateway, item.Code)
				}
			}
		})
	}
}

func TestTransactionsConsumer_Error_UnblockedUser(t *testing.T) {
	tests := []struct {
		name                string
		url                 string
		body                []byte
		fromProcessingTimes string
		msgErr              string
		service             mockBlockedUserService
	}{
		{
			name:                "consumer capture",
			url:                 "/cards/smart-transaction-repair/v1/consumer/transactions/captures/mlm",
			body:                _captureMsgs,
			fromProcessingTimes: "consumer.captures",
			msgErr:              "MLM_consumed_capture_transactions",
			service: mockBlockedUserService{
				BlockUserIfReparationExistsInBulkStub: func(
					ctx context.Context,
					input domain.TransactionsNewsFeedBulkInput,
					wrapFields *field.WrappedFields,
				) (*domain.TransactionsBulkOutput, error) {
					return &domain.TransactionsBulkOutput{
						AuthorizationIDsWithErr: map[string]error{
							"auth_03": gnierrors.Wrap(domain.BadGateway, errors.New("timeout error"), "blockeduser database failure", false),
							"auth_04": errors.New("some unknown error"),
						},
					}, nil
				},
			},
		},
		{
			name:                "consumer reverse",
			url:                 "/cards/smart-transaction-repair/v1/consumer/transactions/reversals/mlm",
			body:                _reverseMsgs,
			fromProcessingTimes: "consumer.reversals",
			msgErr:              "MLM_consumed_reversal_transactions",
			service: mockBlockedUserService{
				UnblockUserByReversalClearingInBulkStub: func(
					ctx context.Context,
					input domain.TransactionsNewsFeedBulkInput,
					wrapFields *field.WrappedFields,
				) (*domain.TransactionsBulkOutput, error) {
					return &domain.TransactionsBulkOutput{
						AuthorizationIDsWithErr: map[string]error{
							"auth_03": gnierrors.Wrap(domain.BadGateway, errors.New("timeout error"), "blockeduser database failure", false),
							"auth_04": errors.New("some unknown error"),
						},
					}, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.url, reader)
			request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
			writer := httptest.NewRecorder()

			log := mockLogService{
				InfolnStub: func(c context.Context, msg string, fields ...log.Field) {
					assert.Equal(t, tt.msgErr, msg)
					for _, f := range fields {
						logFields := []string{
							keyfields,
							keyprocessingTimes,
							"auth_03",
							"auth_04",
							keyFieldAuthorizationIDsWithSuccessful,
							keyFieldResultStatus,
						}
						assert.Contains(t, logFields, f.Key)
						assertFieldsFromWrap(t, f, tt.fromProcessingTimes)
					}
				},
				SetStacktub: func(stacktrace string) log.Field {
					return log.Field{Key: keystacktrace, String: stacktrace}
				},
				AnyStub: func(key string, value interface{}) log.Field {
					return log.Field{Key: key, Interface: value}
				},
			}

			buildRoutes(writer, request, tt.service, log)

			resp := writer.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			bodyResponse, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var bulkResponse bulkResponse
			err = json.Unmarshal(bodyResponse, &bulkResponse)
			assert.NoError(t, err)

			for index, item := range bulkResponse.Responses {
				switch index {
				case 0:
					assert.Equal(t, "01", item.ID)
					assert.Equal(t, http.StatusOK, item.Code)
				case 1:
					assert.Equal(t, "02", item.ID)
					assert.Equal(t, http.StatusOK, item.Code)
				case 2:
					assert.Equal(t, "03", item.ID)
					assert.Equal(t, http.StatusBadGateway, item.Code)
				case 3:
					assert.Equal(t, "04", item.ID)
					assert.Equal(t, http.StatusInternalServerError, item.Code)
				}
			}
		})
	}
}

func TestTransactionsConsumer_Error_DecodeJson(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		msgErr string
	}{
		{
			name:   "consumer capture",
			url:    "/cards/smart-transaction-repair/v1/consumer/transactions/captures/mlm",
			msgErr: "MLM consumer capture transactions failed",
		},
		{
			name:   "consumer reverse",
			url:    "/cards/smart-transaction-repair/v1/consumer/transactions/reversals/mlm",
			msgErr: "MLM consumer reversal transactions failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(nil)
			request := httptest.NewRequest(http.MethodPost, tt.url, reader)
			request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)

			writer := httptest.NewRecorder()

			service := mockBlockedUserService{}

			log := mockLogService{
				ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
					assert.Equal(t, tt.msgErr, msg)
				},
				SetStacktub: func(stacktrace string) log.Field {
					return log.Field{Key: keystacktrace, String: stacktrace}
				},
				AnyStub: func(key string, value interface{}) log.Field {
					return log.Field{Key: key, Interface: value}
				},
			}

			buildRoutes(writer, request, service, log)

			resp := writer.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			expectedErrResponse := errAPI.CustomError{
				Code:    "bad_request",
				Message: tt.msgErr,
				Cause:   "bad_request: EOF",
			}

			assertCustomError(t, resp, expectedErrResponse)
		})
	}
}
