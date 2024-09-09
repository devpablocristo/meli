package transactions

import (
	"context"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

const (
	msgError = "%s consumer %s transactions failed"
)

func (t TransactionsConsumer) buildErrorBadRequestAndLog(
	ctx context.Context,
	err error,
	transactionType transactionType,
	siteID string,
	wrapFields *field.WrappedFields,
) error {

	setFieldsInputParameters(wrapFields, siteID, transactionType, 0)

	gnierr := domain.BuildErrorBadRequest(err, fmt.Sprintf(msgError, siteID, transactionType))
	return t.errAPI.CreateAPIErrorAndLog(ctx, gnierr, wrapFields)
}

func (t TransactionsConsumer) reBuildErrorAndLog(
	ctx context.Context,
	gnierr error,
	transactionType transactionType,
	siteID string,
	wrapFields *field.WrappedFields,
) {

	msg := fmt.Sprintf(msgError, siteID, transactionType)
	gnierr = gnierrors.Wrap(gnierrors.Type(gnierr), gnierr, msg, false)

	_ = t.errAPI.CreateAPIErrorAndLog(ctx, gnierr, wrapFields)
}
