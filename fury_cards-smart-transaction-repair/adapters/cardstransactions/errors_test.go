package cardstransactions

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
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), statusBadRequest, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.BadRequest,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_InternalServerError(t *testing.T) {
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), statusInternalServerError, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.InternalServerError,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Locked(t *testing.T) {
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), statusLocked, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.Locked,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Unauthorized(t *testing.T) {
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), statusUnauthorized, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.Unauthorized,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Forbidden(t *testing.T) {
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), statusForbidden, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.Forbidden,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_NotFound(t *testing.T) {
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), statusNotFound, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_BadGateway(t *testing.T) {
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), statusBadGateway, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_Default(t *testing.T) {
	mockReversalService := &cardsTransactions{
		log: mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		},
	}

	err := mockReversalService.buildError(errors.New("error"), http.StatusFailedDependency, []byte(`{}`))

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("error"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
