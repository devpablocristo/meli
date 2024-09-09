package errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

type ErrAPI struct {
	Log log.LogService
}

// Create error and log.
func (e ErrAPI) CreateAPIErrorAndLog(ctx context.Context, err error, wrapFields *field.WrappedFields) error {
	errType := gnierrors.Type(err)
	var customErr CustomError
	logLevel := log.ErrorLevel

	defer func() {
		e.log(ctx, errType.String(), err, &customErr, logLevel, wrapFields)
	}()

	var status int

	switch errType {
	case domain.NotFound:
		status = http.StatusNotFound

	case domain.AlreadyRepaired, domain.AuthorizationAlreadyReversed:
		logLevel = log.InfoLevel
		status = http.StatusConflict

	case domain.UnprocessableEntity:
		logLevel = log.InfoLevel
		status = http.StatusUnprocessableEntity

	case domain.BadRequest:
		status = http.StatusBadRequest

	case domain.BadGateway:
		status = http.StatusBadGateway

	case domain.Unauthorized:
		status = http.StatusUnauthorized

	case domain.NotEligible:
		logLevel = log.InfoLevel

		err := e.buildErrorWithCustomCause(ctx, err, errType.String(), http.StatusUnprocessableEntity)
		customErr = *err.(*CustomError)
		return &customErr

	case domain.UnableRequest:
		logLevel = log.InfoLevel
		status = http.StatusUnprocessableEntity

	default:
		customErr = CustomError{
			Code:    domain.InternalServerError.String(),
			Status:  http.StatusInternalServerError,
			Message: "sorry, an unknown error occurred",
			Cause:   err.Error(),
		}
		return &customErr
	}

	customErr = CustomError{
		Code:    errType.String(),
		Status:  status,
		Message: gnierrors.Message(err),
		Cause:   e.removeDuplicateMessagesFromCause(err),
	}
	return &customErr
}

func (e ErrAPI) buildErrorWithCustomCause(ctx context.Context, gnierr error, errType string, status int) error {
	customError := CustomError{
		Code:    errType,
		Status:  status,
		Message: gnierrors.Message(gnierr),
	}

	switch v := gnierrors.Cause(gnierr).(type) {
	case *domain.ValidationCauseErrResponse:
		customError.Cause = causeCustom{v.Reason, v.CreationDatetime, v.ReasonDetail}
		return &customError

	default:
		e.Log.Warn(ctx, "error was not customized")
		customError.Cause = gnierr.Error()
		return &customError
	}
}

func (e ErrAPI) BuildErrorUnauthorizedAPI(msgError string) error {
	attr := domain.BuildAttrToLog(e.Log, domain.KeyFieldUnauthorized, "request is not authorized: manual request not allowed")

	errCause := errors.New("request is not authorized")
	return domain.BuildErrorUnauthorized(errCause, msgError, attr)
}

func (e ErrAPI) removeDuplicateMessagesFromCause(gnierr error) string {
	msg := gnierrors.Message(gnierr)
	cause := gnierrors.Cause(gnierr)
	return strings.Replace(cause.Error(), fmt.Sprint(msg, " : "), "", strings.Count(cause.Error(), msg))
}
