package errors

import (
	"context"
	"encoding/json"

	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

func (e ErrAPI) log(
	ctx context.Context,
	errType string,
	gnierr error,
	customErr *CustomError,
	logLevel log.Level,
	wrapFields *field.WrappedFields,
) {
	if errType == "" {
		errType = "unknown"
	}

	logFields := []log.Field{e.Log.Any("type", errType)}
	switch logLevel {
	case log.InfoLevel:
		b, _ := json.Marshal(customErr.Cause)

		logFields = append(logFields, e.Log.Any("detail", string(b)))
		logFields = e.appendCommonFields(ctx, logLevel, gnierr, wrapFields, logFields)
		e.Log.Infoln(ctx, customErr.Message, logFields...)
	default:
		logFields = append(logFields, e.Log.Any("error", customErr.Cause))
		logFields = e.appendCommonFields(ctx, logLevel, gnierr, wrapFields, logFields)
		logFields = append(logFields, e.Log.SetStack(gnierrors.StringStacktrace(gnierr)))
		e.Log.Errorln(ctx, customErr.Message, logFields...)
	}
}

func (e ErrAPI) appendCommonFields(
	ctx context.Context,
	logLevel log.Level,
	gnierr error,
	wrapFields *field.WrappedFields,
	logFields []log.Field,
) []log.Field {
	if field := wrapFields.Fields.ToLogField(logLevel); field != (log.Field{}) {
		logFields = append(logFields, wrapFields.Fields.ToLogField(logLevel))
	}

	logFields = append(logFields, e.addAttrsLogs(ctx, gnierr)...)

	if cast, ok := wrapFields.Timers.ToLogField().Interface.(map[string]int64); ok && len(cast) > 0 {
		logFields = append(logFields, wrapFields.Timers.ToLogField())
	}

	return logFields
}

func (e ErrAPI) addAttrsLogs(ctx context.Context, gnierr error) []log.Field {
	attrs := gnierrors.Attributes(gnierr)
	logFields := []log.Field{}

	for k, attr := range attrs {
		if lg, found := attr.(log.Field); found {
			logFields = append(logFields, lg)
		} else {
			e.Log.Warn(ctx, "conversion error for log field", e.Log.Any(k, k))
		}
	}
	return logFields
}
