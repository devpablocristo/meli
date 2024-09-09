package validationresult

import (
	"context"
	_ "embed"
	"net/http"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/melisource/fury_go-platform/pkg/fury"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/validation-result_msg.json
var _validationResultMsgs []byte

//go:embed testdata/validation-result_domain.json
var _validationResultDomain []byte

// ValidationService.
type mockValidationService struct {
	ports.ValidationService
	SaveStub func(ctx context.Context, inputResult storage.ValidationResult, wrapFields *field.WrappedFields) error
}

func (m mockValidationService) Save(ctx context.Context, inputResult storage.ValidationResult, wrapFields *field.WrappedFields) error {
	return m.SaveStub(ctx, inputResult, wrapFields)
}

func buildRoutes(
	writer http.ResponseWriter,
	request *http.Request,
	service mockValidationService,
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

// LogService.
type mockLogService struct {
	log.LogService
	ErrorlnStub func(ctx context.Context, msg string, fields ...log.Field)
	AnyStub     func(key string, value interface{}) log.Field
	InfoStub    func(c context.Context, msg string, fields ...log.Field)
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

func assertFieldsFromWrap(t *testing.T, f log.Field) {
	switch f.Key {
	case "processing_times":
		maps := f.Interface.(map[string]int64)
		_, found := maps["consumer.validation_result"]
		assert.True(t, found)
	case "fields":
		fields := f.Interface.(map[string]interface{})
		_, found := fields["input_parameters"]
		assert.True(t, found)
	default:
	}
}
