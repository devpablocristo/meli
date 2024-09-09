package validationforcepublishhdl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	utils "github.com/melisource/fury_cards-go-toolkit/pkg/furyutils/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

const (
	basePath = "/cards/smart-transaction-repair/v1/validation"
)

func Test_Publish(t *testing.T) {
	t.Run("should publish a message with success", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodPost,
			basePath+"/force",
			mockBody(t, msgValidationPublish{
				KycIdentitificationID: "81a5c34c25b71fe97832170d483cbd46",
				PaymentID:             "123",
				AuthorizationID:       "AB33",
				Type:                  "reverse",
				UserID:                1309454095,
				SiteID:                "MLM",
				FaqID:                 "FaqID",
				Reason: map[domain.RuleType]publicDefault{
					"qty_reparation_per_period": {
						Actual:   "value1",
						Accepted: "value2",
					},
					"status_detail": {
						Actual:   123,
						Accepted: 456,
					},
				},
				CreatedAt: time.Now(),
			}))

		request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
		request.Header.Add("x-client-id", "123456")
		request.Header.Add("X-Api-Client-Application", "app-test")
		request.Header.Add("X-Api-Client-Scope", "beta")

		writer := httptest.NewRecorder()

		service := &mockValidationService{
			PublishStub: func(ctx context.Context, inputResult *storage.ValidationResult, wrapFields *field.WrappedFields) {
				assert.NotEmpty(t, inputResult)
			},
		}

		mockLogService := mockLogService{}

		buildRoutes(writer, request, service, mockLogService)

		resp := writer.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_Publish_Error_Body_UnprocessableEntity(t *testing.T) {
	t.Run("should return status code 422", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, basePath+"/force", mockBody(t, msgValidationPublish{}))
		request.Header.Add(utils.HeaderKeyAPISource, utils.HeaderValueFurySource)
		writer := httptest.NewRecorder()

		service := &mockValidationService{}

		mockLogService := mockLogService{}

		buildRoutes(writer, request, service, mockLogService)

		resp := writer.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}
