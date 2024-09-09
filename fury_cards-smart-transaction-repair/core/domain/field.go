package domain

import (
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	gnierrors "github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

type keyFieldLog string

const (
	KeyFieldPaymentID              keyFieldLog = "payment_id"
	KeyFieldUserID                 keyFieldLog = "user_id"
	KeyFieldSiteID                 keyFieldLog = "site_id"
	KeyFieldAuthorizationID        keyFieldLog = "authorization_id"
	KeyFieldAuthorizationIDs       keyFieldLog = "authorization_ids"
	KeyFieldBodyReversal           keyFieldLog = "request_reversal"
	KeyFieldBodyGraphqlTransaction keyFieldLog = "response_graphql_transaction"
	KeyFieldRulesApplied           keyFieldLog = "rules_applied"
	KeyFieldOriginalErrType        keyFieldLog = "original_error_type"
	KeyFieldUnauthorized           keyFieldLog = "unauthorized_request"
	KeyFieldInputParameters        keyFieldLog = "input_parameters"
	KeyFieldKycIDs                 keyFieldLog = "kyc_identification_ids"
	KeyFieldBodyPaymentTransaction keyFieldLog = "response_payment_transaction"
)

func BuildAttrToLog(logSrv log.LogService, key keyFieldLog, value interface{}) gnierrors.Attr {
	return gnierrors.Attr{
		Key:   string(key),
		Value: logSrv.Any(string(key), value),
	}
}
