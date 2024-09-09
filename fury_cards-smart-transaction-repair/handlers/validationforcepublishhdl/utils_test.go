package validationforcepublishhdl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/melisource/fury_go-platform/pkg/fury"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

type mockValidationService struct {
	PublishStub func(ctx context.Context, inputResult *storage.ValidationResult, wrapFields *field.WrappedFields)
}

func (m mockValidationService) ExecuteValidation(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
	return nil
}

func (m mockValidationService) Save(ctx context.Context, inputResult storage.ValidationResult, wrapFields *field.WrappedFields) error {
	return nil
}

func (m mockValidationService) PublishEvent(ctx context.Context, validationResult *storage.ValidationResult, wrapFields *field.WrappedFields) {
	m.PublishStub(ctx, validationResult, wrapFields)
}

// LogService
type mockLogService struct {
	log.LogService
}

func buildRoutes(w http.ResponseWriter, req *http.Request, service *mockValidationService, log mockLogService) {
	app, _ := fury.NewWebApplication(fury.WithErrorHandler(web.DefaultErrorHandler))
	router := app.Group("/cards/smart-transaction-repair")
	v1 := router.Group("/v1")

	hdl := NewHandler(service, log)

	NewRouter(hdl).AddRoutesV1(v1)

	app.ServeHTTP(w, req)
}

func mockBody(t *testing.T, reparationRequest msgValidationPublish) *bytes.Buffer {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(reparationRequest)
	assert.NoError(t, err)

	return &body
}
