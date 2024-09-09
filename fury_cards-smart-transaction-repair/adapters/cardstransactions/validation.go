package cardstransactions

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func validateHeader(headers map[string][]string) error {
	var requiredHeaders string

	if value := headers[headerKeyClientID]; value[0] == "" {
		requiredHeaders = fmt.Sprint(requiredHeaders, headerKeyClientID, ", ")
	}

	if value := headers[headerKeySiteID]; value[0] == "" {
		requiredHeaders = fmt.Sprint(requiredHeaders, headerKeySiteID, ", ")
	}

	if value := headers[headerKeyEnvironment]; value[0] == "" {
		requiredHeaders = fmt.Sprint(requiredHeaders, headerKeyEnvironment, ", ")
	}

	if value := headers[headerKeyProvider]; value[0] == "" {
		requiredHeaders = fmt.Sprint(requiredHeaders, headerKeyProvider, ", ")
	}

	if requiredHeaders != "" {
		requiredHeaders = strings.TrimSuffix(requiredHeaders, ", ")
		msg := fmt.Sprint("required headers: ", strings.TrimSpace(requiredHeaders))
		return errors.New(msg)
	}

	return nil
}

func validateRequest(request requestReversal) error {
	var requiredFields string

	if isEmpty(request.ID) {
		requiredFields = fmt.Sprint(requiredFields, "id, ")
	}

	err := validateOperation(request.Operation)
	if err != nil {
		requiredFields = fmt.Sprint(requiredFields, err.Error())
	}

	err = validateProvider(request.Provider)
	if err != nil {
		requiredFields = fmt.Sprint(requiredFields, err.Error())
	}

	err = validateOptions(request.Options)
	if err != nil {
		requiredFields = fmt.Sprint(requiredFields, err.Error())
	}

	if requiredFields != "" {
		requiredFields = strings.TrimSuffix(requiredFields, ", ")
		msg := fmt.Sprint("required fields: ", strings.TrimSpace(requiredFields))
		return errors.New(msg)
	}
	return nil
}

func validateOperation(operation operation) error {
	var requiredFields string

	dateEmpty := time.Time{}
	if operation.TransmissionDatetime == dateEmpty {
		requiredFields = fmt.Sprint(requiredFields, "transmission_datetime, ")
	}

	if operation.Installments <= 0 {
		requiredFields = fmt.Sprint(requiredFields, "installments, ")
	}

	err := validateAuthorization(operation.Authorization)
	if err != nil {
		requiredFields = fmt.Sprint(requiredFields, err.Error())
	}

	err = validateTransaction(operation.Transaction)
	if err != nil {
		requiredFields = fmt.Sprint(requiredFields, err.Error())
	}

	err = validateCard(operation.Card)
	if err != nil {
		requiredFields = fmt.Sprint(requiredFields, err.Error())
	}

	if requiredFields != "" {
		requiredFields = strings.TrimSuffix(requiredFields, ", ")
		msg := fmt.Sprintf("operations: {%s}, ", strings.TrimSpace(requiredFields))
		return errors.New(msg)
	}

	return nil
}

func validateProvider(provider provider) error {
	if isEmpty(provider.ID) {
		return errors.New("provider: {id}, ")
	}
	return nil
}

func validateOptions(options options) error {
	if isEmpty(options.ReversalType) {
		return errors.New("options: {reversal_type}, ")
	}
	return nil
}

func validateAuthorization(authorization authorization) error {
	var requiredFields string

	if isEmpty(authorization.ID) {
		requiredFields = fmt.Sprint(requiredFields, "id, ")
	}

	dateEmpty := time.Time{}
	if authorization.TransmissionDatetime == dateEmpty {
		requiredFields = fmt.Sprint(requiredFields, "transmission_datetime, ")
	}

	if requiredFields != "" {
		requiredFields = strings.TrimSuffix(requiredFields, ", ")
		msg := fmt.Sprintf("authorization: {%s}, ", strings.TrimSpace(requiredFields))
		return errors.New(msg)
	}

	return nil
}

func validateTransaction(transaction transaction) error {
	var requiredFields string

	if isEmpty(transaction.Currency) {
		requiredFields = fmt.Sprint(requiredFields, "currency, ")
	}

	if requiredFields != "" {
		requiredFields = strings.TrimSuffix(requiredFields, ", ")
		msg := fmt.Sprintf("transaction: {%s}, ", strings.TrimSpace(requiredFields))
		return errors.New(msg)
	}

	return nil
}

func validateCard(card card) error {
	var requiredFields string

	if isEmpty(card.NumberId) {
		requiredFields = fmt.Sprint(requiredFields, "number_id, ")
	}

	if isEmpty(card.Country) {
		requiredFields = fmt.Sprint(requiredFields, "country, ")
	}

	if requiredFields != "" {
		requiredFields = strings.TrimSuffix(requiredFields, ", ")
		msg := fmt.Sprintf("card: {%s}, ", strings.TrimSpace(requiredFields))
		return errors.New(msg)
	}

	return nil
}

func isEmpty(v string) bool {
	return strings.TrimSpace(v) == ""
}
