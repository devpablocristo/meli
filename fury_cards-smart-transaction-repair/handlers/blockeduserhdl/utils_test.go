package blockeduserhdl

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/melisource/fury_go-platform/pkg/fury"
	"github.com/stretchr/testify/assert"
)

// BlockedUserService
type mockBlockedUserService struct {
	ports.BlockedUserService
	UnblockUserStub func(ctx context.Context, kycIdentificationID string) error
}

func (m mockBlockedUserService) UnblockUser(ctx context.Context, kycIdentificationID string) error {
	return m.UnblockUserStub(ctx, kycIdentificationID)
}

// LogService
type mockLogService struct {
	log.LogService
	ErrorlnStub func(ctx context.Context, msg string, fields ...log.Field)
	SetStacktub func(stacktrace string) log.Field
	AnyStub     func(key string, value interface{}) log.Field
}

func (m mockLogService) Errorln(c context.Context, msg string, fields ...log.Field) {
	m.ErrorlnStub(c, msg, fields...)
}
func (m mockLogService) Any(key string, value interface{}) log.Field {
	return m.AnyStub(key, value)
}
func (m mockLogService) SetStack(stacktrace string) log.Field {
	return m.SetStacktub(stacktrace)
}

func buildRoutes(
	writer http.ResponseWriter,
	request *http.Request,
	service *mockBlockedUserService,
	log mockLogService) {

	app, _ := fury.NewWebApplication(fury.WithErrorHandler(web.DefaultErrorHandler))
	router := app.Group("/cards/smart-transaction-repair")
	v1 := router.Group("/v1")

	hdl := NewHandler(service, log)

	NewRouter(hdl).AddRoutesV1(v1)

	app.ServeHTTP(writer, request)
}

func assertCustomError(t *testing.T, resp *http.Response, expectedErrResponse errAPI.CustomError) {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	errResponse := errAPI.CustomError{}
	err = json.Unmarshal(body, &errResponse)

	assert.NoError(t, err)

	assert.Equal(t, expectedErrResponse, errResponse)
}
