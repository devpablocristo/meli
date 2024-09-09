package errors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

const keytype, keyinfo = "type", "detail"

func TestErrorAPI_Gnierror_without_attrs(t *testing.T) {
	gnierr := gnierrors.Wrap(domain.AlreadyRepaired, errors.New("cause-error"), "msg-err", false)

	mockLogService := mockLogService{
		InfolnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "msg-err", msg)
			assert.Equal(t, keytype, fields[0].Key)
			assert.Equal(t, keyinfo, fields[1].Key)
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	errApi := ErrAPI{Log: mockLogService}
	err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, field.NewWrappedFields())

	errExpected := &CustomError{
		Status:  409,
		Code:    "already_repaired",
		Message: "msg-err",
		Cause:   "cause-error",
	}

	assert.Equal(t, errExpected, err)
}

func TestErrorAPI_Gnierror_with_attrs(t *testing.T) {
	logSrv := log.New()

	logField1 := logSrv.Any("key1", "value1")
	logField2 := logSrv.Any("key3", "value3")

	attr1 := gnierrors.Attr{Key: "any1", Value: logField1}
	attr2 := gnierrors.Attr{Key: "any2", Value: logField2}

	gnierr := gnierrors.Wrap(domain.UnprocessableEntity, errors.New("cause-error"), "msg-err", false, attr1, attr2)

	mockLogService := mockLogService{
		InfolnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "msg-err", msg)
			assert.NotEmpty(t, fields)
			assert.Len(t, fields, 6)
			assert.Equal(t, keytype, fields[0].Key)
			assert.Equal(t, keyinfo, fields[1].Key)
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	wrapFields := field.NewWrappedFields()

	wrapFields.Fields.Add("payment_id", "100")
	wrapFields.Fields.Add("test_id", "x")
	start := wrapFields.Timers.Start("kvs")
	start.Stop()

	errApi := ErrAPI{Log: mockLogService}
	err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, wrapFields)

	errExpected := &CustomError{
		Status:  422,
		Code:    "unprocessable_entity",
		Message: "msg-err",
		Cause:   "cause-error",
	}

	assert.Equal(t, errExpected, err)
}

func TestErrorAPI_Gnierror_with_attr_diff(t *testing.T) {

	attr1 := gnierrors.Attr{Key: "warn", Value: "diff"}

	gnierr := gnierrors.Wrap(domain.BadGateway, errors.New("cause-error"), "msg-err", false, attr1)

	mockLogService := mockLogService{
		ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "msg-err", msg)
			assert.NotEmpty(t, fields)
			assert.Len(t, fields, 4)
			assert.Equal(t, keytype, fields[0].Key)
			assert.Equal(t, "error", fields[1].Key)
		},
		WarnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "conversion error for log field", msg)
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
		SetStacktub: func(stacktrace string) log.Field {
			return log.Field{Key: "stacktrace", String: "stacktrace"}
		},
	}

	errApi := ErrAPI{Log: mockLogService}
	err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, field.NewWrappedFields())

	errExpected := &CustomError{
		Status:  502,
		Code:    "bad_gateway",
		Message: "msg-err",
		Cause:   "cause-error",
	}

	assert.Equal(t, errExpected, err)
}

func TestErrorAPI_Error(t *testing.T) {

	mockLogService := mockLogService{
		ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "sorry, an unknown error occurred", msg)
			assert.NotEmpty(t, fields)
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
		SetStacktub: func(stacktrace string) log.Field {
			return log.Field{Key: "stacktrace", String: "stacktrace"}
		},
	}

	errApi := ErrAPI{Log: mockLogService}
	err := errApi.CreateAPIErrorAndLog(context.TODO(), errors.New("error-test"), field.NewWrappedFields())

	errExpected := &CustomError{
		Status:  500,
		Code:    "internal_server_error",
		Message: "sorry, an unknown error occurred",
		Cause:   "error-test",
	}

	assert.Equal(t, errExpected, err)
}

func TestErrorAPI_Gnierror_others_types(t *testing.T) {
	mockLogService := mockLogService{
		ErrorlnStub: func(c context.Context, msg string, fields ...log.Field) {},
		InfolnStub:  func(c context.Context, msg string, fields ...log.Field) {},
		SetStacktub: func(stacktrace string) log.Field {
			return log.Field{}
		},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
		WarnStub: func(c context.Context, msg string, fields ...log.Field) {},
	}

	errApi := ErrAPI{Log: mockLogService}

	t.Run("bad_request", func(t *testing.T) {
		gnierr := gnierrors.Wrap(domain.BadRequest, errors.New("cause-error"), "msg-err", false)

		wrappedFields := field.NewWrappedFields()
		wrappedFields.Fields.Add(string(domain.KeyFieldBodyGraphqlTransaction), "{}")

		err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, wrappedFields)

		errExpected := &CustomError{
			Status:  400,
			Code:    "bad_request",
			Message: "msg-err",
			Cause:   "cause-error",
		}

		assert.Equal(t, errExpected, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		gnierr := gnierrors.Wrap(domain.Unauthorized, errors.New("cause-error"), "msg-err", false)

		err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, field.NewWrappedFields())

		errExpected := &CustomError{
			Status:  401,
			Code:    "unauthorized",
			Message: "msg-err",
			Cause:   "cause-error",
		}

		assert.Equal(t, errExpected, err)
	})

	t.Run("not_eligible", func(t *testing.T) {
		creationDatetime := time.Now()
		reasonDetail := domain.Reason{domain.RuleMaxAmountReparation: &domain.ReasonResult{
			Actual:   2500,
			Accepted: 100,
		}}

		causeErr := &domain.ValidationCauseErrResponse{
			Reason:           "reason",
			CreationDatetime: creationDatetime,
			ReasonDetail:     reasonDetail,
		}
		gnierr := gnierrors.Wrap(domain.NotEligible, causeErr, "msg-err", false)

		err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, field.NewWrappedFields())

		errExpected := &CustomError{
			Status:  422,
			Code:    "not_eligible",
			Message: "msg-err",
			Cause:   causeCustom{"reason", creationDatetime, reasonDetail},
		}

		assert.Equal(t, errExpected, err)
	})

	t.Run("not_eligible - unknown", func(t *testing.T) {
		gnierr := gnierrors.Wrap(domain.NotEligible, errors.New("x"), "msg-err", false)

		err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, field.NewWrappedFields())

		errExpected := &CustomError{
			Status:  422,
			Code:    "not_eligible",
			Message: "msg-err",
			Cause:   gnierr.Error(),
		}

		assert.Equal(t, errExpected, err)
	})

	t.Run("not_found", func(t *testing.T) {
		gnierr := gnierrors.Wrap(domain.NotFound, errors.New("cause-error"), "msg-err", false)

		err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, field.NewWrappedFields())

		errExpected := &CustomError{
			Status:  404,
			Code:    "not_found",
			Message: "msg-err",
			Cause:   "cause-error",
		}

		assert.Equal(t, errExpected, err)
	})

	t.Run("unable_request", func(t *testing.T) {
		gnierr := gnierrors.Wrap(domain.UnableRequest, errors.New("cause-error"), "msg-err", false)

		err := errApi.CreateAPIErrorAndLog(context.TODO(), gnierr, field.NewWrappedFields())

		errExpected := &CustomError{
			Status:  422,
			Code:    "unable_request",
			Message: "msg-err",
			Cause:   "cause-error",
		}

		assert.Equal(t, errExpected, err)
	})
}

func TestCustomError(t *testing.T) {
	err := &CustomError{
		Status:  409,
		Code:    "conflict",
		Message: "msg-err",
		Cause:   "cause-error",
	}

	bErr, errMarshal := err.MarshalJSON()

	assert.NoError(t, errMarshal)
	assert.Equal(t, `{"code":"conflict","message":"msg-err","cause":"cause-error"}`, string(bErr))
	assert.Equal(t, "conflict: msg-err: cause-error:", err.Error())
	assert.Equal(t, 409, err.StatusCode())
}

func TestErrAPI_BuildErrorUnauthorizedAPI(t *testing.T) {
	mockLogService := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.Equal(t, key, "unauthorized_request")
			assert.Equal(t, value, "request is not authorized: manual request not allowed")
			return log.Field{Key: key, Interface: value}
		},
	}

	errApi := ErrAPI{Log: mockLogService}

	gnierr := gnierrors.Wrap(domain.BadGateway, errors.New("request is not authorized"), "invalid request", false)
	err := errApi.BuildErrorUnauthorizedAPI("invalid request")

	assert.EqualError(t, err, gnierr.Error())
}
