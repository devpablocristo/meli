package domain

import (
	"encoding/json"
	"fmt"
	"time"

	gnierrors "github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

var (
	InternalServerError          = gnierrors.Define("internal_server_error")
	NotFound                     = gnierrors.Define("not_found")
	BadRequest                   = gnierrors.Define("bad_request")
	BadGateway                   = gnierrors.Define("bad_gateway")
	DuplicatedKey                = gnierrors.Define("duplicated_key")
	Locked                       = gnierrors.Define("locked")
	Forbidden                    = gnierrors.Define("forbidden")
	Unauthorized                 = gnierrors.Define("unauthorized")
	UnprocessableEntity          = gnierrors.Define("unprocessable_entity")
	NotEligible                  = gnierrors.Define("not_eligible")
	AlreadyRepaired              = gnierrors.Define("already_repaired")
	UnavailableForLegalReasons   = gnierrors.Define("unavailable_for_legal_reasons")
	UnableRequest                = gnierrors.Define("unable_request")
	AuthorizationAlreadyReversed = gnierrors.Define("authorization_already_reversed")
)

type ValidationCauseErrResponse struct {
	Reason           string
	CreationDatetime time.Time
	ReasonDetail     Reason
}

func (v *ValidationCauseErrResponse) Error() string {
	b, _ := json.Marshal(v.ReasonDetail)
	return string(b)
}

func BuildErrorNotFound(id, messageCause, message string, attributes ...gnierrors.Attr) error {
	errCause := fmt.Errorf(fmt.Sprint(messageCause, " | ID: %s"), id)
	return gnierrors.Wrap(NotFound, errCause, message, false, attributes...)
}

func BuildErrorNotFoundGeneric(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(NotFound, err, message, false, attributes...)
}

func BuildErrorBadGateway(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(BadGateway, err, message, false, attributes...)
}

func BuildErrorBadRequest(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(BadRequest, err, message, false, attributes...)
}

func BuildErrorDuplicateKey(id string, message string, attributes ...gnierrors.Attr) error {
	errCause := fmt.Errorf("key %s exists", id)
	return gnierrors.Wrap(DuplicatedKey, errCause, message, false, attributes...)
}

func BuildErrorInternalServerError(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(InternalServerError, err, message, false, attributes...)
}

func BuildErrorLocked(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(Locked, err, message, false, attributes...)
}

func BuildErrorUnauthorized(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(Unauthorized, err, message, false, attributes...)
}

func BuildErrorForbidden(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(Forbidden, err, message, false, attributes...)
}

func BuildErrorAlreadyRepaired(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(AlreadyRepaired, err, message, false, attributes...)
}

func BuildErrorUnprocessableEntity(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(UnprocessableEntity, err, message, false, attributes...)
}

func BuildErrorNotEligible(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(NotEligible, err, message, false, attributes...)
}

func BuildErrorUnavailableForLegalReasons(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(UnavailableForLegalReasons, err, message, false, attributes...)
}

func BuildErrorUnableRequest(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(UnableRequest, err, message, false, attributes...)
}

func BuildErrorAuthorizationAlreadyReversed(err error, message string, attributes ...gnierrors.Attr) error {
	return gnierrors.Wrap(AuthorizationAlreadyReversed, err, message, false, attributes...)
}
