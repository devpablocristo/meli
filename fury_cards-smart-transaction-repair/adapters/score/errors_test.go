package score

import (
	"errors"
	"net/http"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_buildError_BadRequest(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), statusBadRequest)

	errorExpected := gnierrors.Wrap(domain.BadRequest,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_InternalServerError(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), statusInternalServerError)

	errorExpected := gnierrors.Wrap(domain.InternalServerError,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Locked(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), statusLocked)

	errorExpected := gnierrors.Wrap(domain.Locked,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Unauthorized(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), statusUnauthorized)

	errorExpected := gnierrors.Wrap(domain.Unauthorized,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Forbidden(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), statusForbidden)

	errorExpected := gnierrors.Wrap(domain.Forbidden,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_NotFound(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), statusNotFound)

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_BadGateway(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), statusBadGateway)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Default(t *testing.T) {
	mockScoreService := &score{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_score",
					String: value.(string),
				}
			},
		},
	}

	err := mockScoreService.buildError(errors.New("error"), http.StatusFailedDependency)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("error"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
