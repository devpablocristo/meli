package cardstransactions

import (
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type statusCodeCardsTransactions int

// expected errors
const (
	statusBadRequest          statusCodeCardsTransactions = http.StatusBadRequest
	statusInternalServerError statusCodeCardsTransactions = http.StatusInternalServerError
	statusLocked              statusCodeCardsTransactions = http.StatusLocked
	statusUnauthorized        statusCodeCardsTransactions = http.StatusUnauthorized
	statusForbidden           statusCodeCardsTransactions = http.StatusForbidden
	statusNotFound            statusCodeCardsTransactions = http.StatusNotFound
	statusBadGateway          statusCodeCardsTransactions = http.StatusBadGateway
)

const (
	msgErrorSearch = "request transaction reversal failed"
)

type bodyResponseError []byte

func (b bodyResponseError) Error() string {
	return string(b)
}

func (c cardsTransactions) buildError(err error, statusCode statusCodeCardsTransactions, bodyRequest []byte) error {
	attrBody := domain.BuildAttrToLog(c.log, domain.KeyFieldBodyReversal, string(bodyRequest))

	switch statusCode {

	case statusBadRequest:
		return domain.BuildErrorBadRequest(err, msgErrorSearch, attrBody)

	case statusInternalServerError:
		return domain.BuildErrorInternalServerError(err, msgErrorSearch, attrBody)

	case statusLocked:
		return domain.BuildErrorLocked(err, msgErrorSearch, attrBody)

	case statusUnauthorized:
		return domain.BuildErrorUnauthorized(err, msgErrorSearch, attrBody)

	case statusForbidden:
		return domain.BuildErrorForbidden(err, msgErrorSearch, attrBody)

	case statusBadGateway:
		return domain.BuildErrorBadGateway(err, msgErrorSearch, attrBody)

	case statusNotFound:
		return domain.BuildErrorNotFoundGeneric(err, msgErrorSearch, attrBody)

	default:
		return domain.BuildErrorBadGateway(err, msgErrorSearch, attrBody)
	}
}

func (c cardsTransactions) buildErrorAuthorizationAlreadyReversed(err error, bodyRequest []byte) error {
	attrBody := domain.BuildAttrToLog(c.log, domain.KeyFieldBodyReversal, string(bodyRequest))

	return domain.BuildErrorAuthorizationAlreadyReversed(err, msgErrorSearch, attrBody)
}
