package reparation

import (
	"context"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	logpkg "github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

const (
	msgErrorRequest = "invalid request"
)

func buildErrorAlreadyRepaired(cause error, msg string) error {
	return domain.BuildErrorAlreadyRepaired(cause, msg)
}

func buildErrorBadRequest(cause error, msg string) error {
	return domain.BuildErrorBadRequest(cause, msg)
}

func buildErrorUnableRequest(cause error, msg string) error {
	return domain.BuildErrorUnableRequest(cause, msg)
}

func buildErrorAuthorizationAlreadyReversed(cause error, msg string) error {
	return domain.BuildErrorAuthorizationAlreadyReversed(cause, msg)
}

func reBuildError(ctx context.Context, log logpkg.LogService, errorType string, gnierror error) error {
	msg := gnierrors.Message(gnierror)
	attrErrType := domain.BuildAttrToLog(log, domain.KeyFieldOriginalErrType, gnierrors.Type(gnierror).String())

	switch errorType {

	case domain.BadGateway.String():
		return domain.BuildErrorBadGateway(gnierror, msg, attrErrType)

	default:
		log.Warn(ctx, fmt.Sprint("unmapped error type: ", errorType))
		return domain.BuildErrorBadGateway(gnierror, msg, attrErrType)
	}
}

func buildErrorFromValidationService(gnierror error) error {
	return gnierror
}
