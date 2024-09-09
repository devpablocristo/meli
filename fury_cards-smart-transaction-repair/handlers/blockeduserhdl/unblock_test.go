package blockeduserhdl

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	errAPI "github.com/melisource/cards-smart-transaction-repair/cmd/errors"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func TestBlockedUserHandler_UnblockUser(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/cards/smart-transaction-repair/v1/blockedlist/users/abc/unblock", nil)
	writer := httptest.NewRecorder()

	service := &mockBlockedUserService{
		UnblockUserStub: func(ctx context.Context, kycIdentificationID string) error {
			assert.Equal(t, "abc", kycIdentificationID)
			return nil
		},
	}

	log := mockLogService{}

	buildRoutes(writer, request, service, log)

	resp := writer.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, string(body), `{"message":"user abc successfully unlocked"}`)
}

func TestBlockedUserHandler_UnblockUser_Error_NotFound(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/cards/smart-transaction-repair/v1/blockedlist/users/abc/unblock", nil)
	writer := httptest.NewRecorder()

	service := &mockBlockedUserService{
		UnblockUserStub: func(ctx context.Context, kycIdentificationID string) error {
			assert.Equal(t, "abc", kycIdentificationID)
			return gnierrors.Wrap(domain.NotFound, errors.New("not found"), "blockeduser service failure", false)
		},
	}

	log := mockLogService{
		ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "blockeduser service failure", msg)
		},
		SetStacktub: func(stacktrace string) log.Field {
			return log.Field{Key: "stacktrace", String: "stacktrace"}
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	buildRoutes(writer, request, service, log)

	resp := writer.Result()
	defer resp.Body.Close()

	expectedErrResponse := errAPI.CustomError{
		Code:    "not_found",
		Message: "blockeduser service failure",
		Cause:   "not found",
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	assertCustomError(t, resp, expectedErrResponse)
}
