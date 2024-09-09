package pinghdl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/melisource/fury_go-platform/pkg/fury"
	"github.com/stretchr/testify/assert"
)

func TestPingHandlerRouter_AddRoutePing(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/cards/smart-transaction-repair/ping", nil)
	writer := httptest.NewRecorder()

	buildRoutes(writer, request)

	resp := writer.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func buildRoutes(
	writer http.ResponseWriter,
	request *http.Request,
) {

	app, _ := fury.NewWebApplication(fury.WithErrorHandler(web.DefaultErrorHandler))
	router := app.Group("/cards/smart-transaction-repair")

	NewRouter().AddRoutePing(router)

	app.ServeHTTP(writer, request)
}
