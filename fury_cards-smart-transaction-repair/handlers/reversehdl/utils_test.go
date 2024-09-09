package reversehdl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/melisource/fury_go-platform/pkg/fury"
	"github.com/stretchr/testify/assert"
)

// ReparationService
type mockReparationService struct {
	ReverseStub func(ctx context.Context, reversal domain.ReverseTransactionInput, wrapFields *field.WrappedFields) error
}

func (m *mockReparationService) Reverse(
	ctx context.Context,
	reversal domain.ReverseTransactionInput,
	wrapFields *field.WrappedFields,
) error {
	return m.ReverseStub(ctx, reversal, wrapFields)
}

// LogService
type mockLogService struct {
	log.LogService
	ErrorlnStub func(ctx context.Context, msg string, fields ...log.Field)
	AnyStub     func(key string, value interface{}) log.Field
	InfoStub    func(c context.Context, msg string, fields ...log.Field)
	InfolnStub  func(ctx context.Context, msg string, fields ...log.Field)
	SetStacktub func(stacktrace string) log.Field
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

func buildRoutes(
	writer http.ResponseWriter,
	request *http.Request,
	service *mockReparationService,
	mockLogService mockLogService) {

	app, _ := fury.NewWebApplication(fury.WithErrorHandler(web.DefaultErrorHandler))
	router := app.Group("/cards/smart-transaction-repair")
	v1 := router.Group("/v1")

	hdl := NewHandler(service, mockLogService)

	NewRouter(hdl).AddRoutesV1(v1)

	app.ServeHTTP(writer, request)
}

func mockBody(t *testing.T, reparationRequest reparationRequest) *bytes.Buffer {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(reparationRequest)
	assert.NoError(t, err)

	return &body
}

func assertFieldsFromWrap(t *testing.T, f log.Field) {
	switch f.Key {
	case "processing_times":
		maps := f.Interface.(map[string]int64)
		_, found := maps["handler.reverse"]
		assert.True(t, found)
	case "fields":
		fields := f.Interface.(map[string]interface{})
		_, found := fields["input_parameters"]
		assert.True(t, found)

	default:
	}
}
