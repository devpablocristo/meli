package transactions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

const (
	keyFieldAuthorizationIDsWithSuccessful = "authorization_ids_with_successful"
	keyFieldTransactionType                = "transaction_type"
	keyFieldTotalMsgsQ                     = "total_msgs"
	keyFieldResultStatus                   = "result_status"
)

func setFieldsInputParameters(wrapFields *field.WrappedFields, siteID string, transactionType transactionType, totalAuthorizationIDs int) {
	parameters := fmt.Sprintf(`{"%s":"%s","%s":"%s","%s":"%s"}`,
		domain.KeyFieldSiteID,
		siteID,
		keyFieldTransactionType,
		fmt.Sprint("transaction_", transactionType),
		keyFieldTotalMsgsQ,
		strconv.Itoa(totalAuthorizationIDs),
	)

	wrapFields.Fields.Add(string(domain.KeyFieldInputParameters), parameters)
}

func (t TransactionsConsumer) logResult(
	ctx context.Context,
	output domain.TransactionsBulkOutput,
	authorizationsIDsSuccess []string,
	transactionType transactionType,
	siteID string,
	wrapFields *field.WrappedFields,
) {
	logFields := []log.Field{wrapFields.Fields.ToLogField(log.InfoLevel)}

	if len(output.AuthorizationIDsWithErr) > 0 {
		logFields = append(logFields, t.log.Any(keyFieldResultStatus, "with_errors"))

		for auth, err := range output.AuthorizationIDsWithErr {
			logFields = append(logFields, t.log.Any(auth, fmt.Sprintf(`{"%s":"%s"}`, "error", err.Error())))
		}
	} else {
		logFields = append(logFields, t.log.Any(keyFieldResultStatus, "full_success"))
	}

	if len(authorizationsIDsSuccess) > 0 {
		logFields = append(logFields,
			t.log.Any(keyFieldAuthorizationIDsWithSuccessful, fmt.Sprintf(`{"%s"}`, strings.Join(authorizationsIDsSuccess, ", "))))
	}

	logFields = append(logFields, wrapFields.Timers.ToLogField())

	t.log.Infoln(ctx, fmt.Sprintf("%s_consumed_%s_transactions", siteID, transactionType), logFields...)
}
