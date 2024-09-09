package transactions

import (
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/melisource/fury_go-platform/pkg/fury"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/capture_msgs.json
var _captureMsgs []byte

//go:embed testdata/reversal_msgs.json
var _reverseMsgs []byte

// BlockedUserService
type mockBlockedUserService struct {
	ports.BlockedUserService

	BlockUserIfReparationExistsInBulkStub func(
		ctx context.Context,
		input domain.TransactionsNewsFeedBulkInput,
		wrapFields *field.WrappedFields) (*domain.TransactionsBulkOutput, error)

	UnblockUserByReversalClearingInBulkStub func(
		ctx context.Context,
		input domain.TransactionsNewsFeedBulkInput,
		wrapFields *field.WrappedFields) (*domain.TransactionsBulkOutput, error)
}

func (m mockBlockedUserService) BlockUserIfReparationExistsInBulk(
	ctx context.Context,
	input domain.TransactionsNewsFeedBulkInput,
	wrapFields *field.WrappedFields,
) (*domain.TransactionsBulkOutput, error) {
	return m.BlockUserIfReparationExistsInBulkStub(ctx, input, wrapFields)
}

func (m mockBlockedUserService) UnblockUserByReversalClearingInBulk(
	ctx context.Context,
	input domain.TransactionsNewsFeedBulkInput,
	wrapFields *field.WrappedFields,
) (*domain.TransactionsBulkOutput, error) {
	return m.UnblockUserByReversalClearingInBulkStub(ctx, input, wrapFields)
}

// LogService.
type mockLogService struct {
	log.LogService
	ErrorlnStub func(ctx context.Context, msg string, fields ...log.Field)
	AnyStub     func(key string, value interface{}) log.Field
	InfoStub    func(c context.Context, msg string, fields ...log.Field)
	SetStacktub func(stacktrace string) log.Field
	InfolnStub  func(ctx context.Context, msg string, fields ...log.Field)
}

func (m mockLogService) Errorln(c context.Context, msg string, fields ...log.Field) {
	m.ErrorlnStub(c, msg, fields...)
}
func (m mockLogService) Any(key string, value interface{}) log.Field {
	return m.AnyStub(key, value)
}
func (m mockLogService) Info(c context.Context, msg string, fields ...log.Field) {
	m.InfoStub(c, msg, fields...)
}
func (m mockLogService) SetStack(stacktrace string) log.Field {
	return m.SetStacktub(stacktrace)
}
func (m mockLogService) Infoln(c context.Context, msg string, fields ...log.Field) {
	m.InfolnStub(c, msg, fields...)
}

func assertCustomError(t *testing.T, resp *http.Response, expectedErrResponse errAPI.CustomError) {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	errResponse := errAPI.CustomError{}
	err = json.Unmarshal(body, &errResponse)

	assert.NoError(t, err)

	assert.Equal(t, expectedErrResponse, errResponse)
}

func assertFieldsFromWrap(t *testing.T, f log.Field, from string) {
	switch f.Key {
	case "processing_times":
		maps := f.Interface.(map[string]int64)
		_, found := maps[from]
		assert.True(t, found)
	case "fields":
		fields := f.Interface.(map[string]interface{})
		_, found := fields["input_parameters"]
		_, found2 := fields["authorization_ids"]
		assert.True(t, found)
		assert.True(t, found2)
	default:
	}
}

func buildRoutes(
	writer http.ResponseWriter,
	request *http.Request,
	service mockBlockedUserService,
	log mockLogService,

) {
	app, _ := fury.NewWebApplication(fury.WithErrorHandler(web.DefaultErrorHandler))
	router := app.Group("/cards/smart-transaction-repair")
	v1 := router.Group("/v1")

	consumer := NewConsumer(service, log)

	transactionsRouter := NewRouter(consumer)
	transactionsRouter.AddRoutesV1(v1)

	app.ServeHTTP(writer, request)
}
