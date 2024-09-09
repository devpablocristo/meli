package score

import (
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type statusCodeScore int

const (
	statusBadRequest          statusCodeScore = http.StatusBadRequest
	statusInternalServerError statusCodeScore = http.StatusInternalServerError
	statusLocked              statusCodeScore = http.StatusLocked
	statusUnauthorized        statusCodeScore = http.StatusUnauthorized
	statusForbidden           statusCodeScore = http.StatusForbidden
	statusNotFound            statusCodeScore = http.StatusNotFound
	statusBadGateway          statusCodeScore = http.StatusBadGateway
)

const (
	msgErrorSearch = "request_score_failed"
)

type bodyResponseError []byte

func (b bodyResponseError) Error() string {
	return string(b)
}

func (s score) buildError(err error, statusCode statusCodeScore) error {
	switch statusCode {
	case statusBadRequest:
		return domain.BuildErrorBadRequest(err, msgErrorSearch)

	case statusInternalServerError:
		return domain.BuildErrorInternalServerError(err, msgErrorSearch)

	case statusLocked:
		return domain.BuildErrorLocked(err, msgErrorSearch)

	case statusUnauthorized:
		return domain.BuildErrorUnauthorized(err, msgErrorSearch)

	case statusForbidden:
		return domain.BuildErrorForbidden(err, msgErrorSearch)

	case statusBadGateway:
		return domain.BuildErrorBadGateway(err, msgErrorSearch)

	case statusNotFound:
		return domain.BuildErrorNotFoundGeneric(err, msgErrorSearch)

	default:
		return domain.BuildErrorBadGateway(err, msgErrorSearch)
	}
}
