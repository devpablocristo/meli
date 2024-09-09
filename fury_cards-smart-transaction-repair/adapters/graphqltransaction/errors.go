package graphqltransaction

import (
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	msgErrorSearch = "request transaction search failed"
)

type bodyResponseError []byte

func (b bodyResponseError) Error() string {
	return string(b)
}

func buildError(err error, statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return domain.BuildErrorBadRequest(err, msgErrorSearch)

	case http.StatusInternalServerError:
		return domain.BuildErrorInternalServerError(err, msgErrorSearch)

	default:
		return domain.BuildErrorBadGateway(err, msgErrorSearch)
	}
}
