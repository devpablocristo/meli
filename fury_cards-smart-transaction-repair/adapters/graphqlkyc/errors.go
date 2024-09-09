package graphqlkyc

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	msgErrFailGetUser = "graphql-search-kyc get_user_kyc failure"
)

func buildErrorBadGateway(err error, msgErr string) error {
	return domain.BuildErrorBadGateway(err, msgErr)
}

func buildErrorUnavailableForLegalReasons(err error, msgErr string) error {
	return domain.BuildErrorUnavailableForLegalReasons(err, msgErr)
}
