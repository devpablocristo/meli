package validationresult

import (
	"context"
	"fmt"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

const (
	msgError = "consumer_validation_result_failed"
)

func (v ValidationResultConsumer) buildErrorBadRequestAndLog(
	ctx context.Context,
	err error,
	wrapFields *field.WrappedFields,
) error {

	gnierr := domain.BuildErrorBadRequest(err, msgError)
	return v.errAPI.CreateAPIErrorAndLog(ctx, gnierr, wrapFields)
}

func (v ValidationResultConsumer) reBuildErrorAndLog(
	ctx context.Context,
	gnierr error,
	siteID string,
	wrapFields *field.WrappedFields,
) error {

	msg := fmt.Sprintf(`%s_%s`, siteID, msgError)
	gnierr = gnierrors.Wrap(gnierrors.Type(gnierr), gnierr, msg, false)

	return v.errAPI.CreateAPIErrorAndLog(ctx, gnierr, wrapFields)
}
