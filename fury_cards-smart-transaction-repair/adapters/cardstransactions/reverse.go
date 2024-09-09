package cardstransactions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
)

// header key
const (
	headerKeyClientID       = "X-Client-Id"
	headerKeySiteID         = "X-Site-Id"
	headerKeyEnvironment    = "X-Environment"
	headerKeyProvider       = "X-Provider"
	headerKeyProviderSiteID = "X-Provider-Site"
	headerKeyBetaScope      = "X-Beta-Scope"
)

// header value
const (
	headerValueContentType = "application/json"
)

// others
const (
	pathURL                                          = "/cards/transactions/debit/reversals"
	responseStatusApproved                           = "approved"
	responseStatusDeclined                           = "declined"
	responseStatusDetailAuthorizationAlreadyReversed = "authorization_already_reversed"
	typeSmartReverse                                 = "smart_reverse"
)

func (c cardsTransactions) Reverse(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
	header := c.fillHeader(*input.TransactionData, input.HeaderValueXClientID)
	err := validateHeader(header)
	if err != nil {
		return nil, c.buildError(err, statusBadRequest, nil)
	}

	request := fillRequest(*input.TransactionData)
	err = validateRequest(request)
	if err != nil {
		return nil, c.buildError(err, statusBadRequest, nil)
	}

	bodyRequest, err := c.json.Marshal(request)
	if err != nil {
		return nil, c.buildError(err, statusInternalServerError, bodyRequest)
	}

	baseURL := c.getBaseURL(strings.ToLower(input.TransactionData.Provider.ID))

	url := rest.CreateURL(baseURL, pathURL)

	response, err := c.httpClient.Do(ctx, "POST", url, headerValueContentType, bodyRequest, header)
	if err != nil {
		return nil, c.buildError(err, statusBadGateway, bodyRequest)
	}

	body, err := c.readAndValidateResponse(response, bodyRequest)
	if err != nil {
		return nil, err
	}

	var responseReversal responseReversal
	err = c.json.Unmarshal(body, &responseReversal)
	if err != nil {
		return nil, c.buildError(err, statusInternalServerError, bodyRequest)
	}

	err = c.checkResult(responseReversal, bodyRequest)
	if err != nil {
		return nil, err
	}

	return fillOutput(request), nil
}

func (c cardsTransactions) checkResult(responseReversal responseReversal, bodyRequest []byte) error {
	msgErrMarshalResponse := "reversal not approved. unfortunately it was not possible to convert the reason (%s)"
	switch {

	case responseReversal.Result.Status == responseStatusDeclined &&
		responseReversal.Result.StatusDetail != nil &&
		*responseReversal.Result.StatusDetail == responseStatusDetailAuthorizationAlreadyReversed:
		bodyResult, err := c.json.Marshal(responseReversal)
		if err != nil {
			bodyResult = []byte(fmt.Sprintf(msgErrMarshalResponse, err))
		}
		return c.buildErrorAuthorizationAlreadyReversed(bodyResponseError(bodyResult), bodyRequest)

	case responseReversal.Result.Status != responseStatusApproved:
		bodyResult, err := c.json.Marshal(responseReversal)
		if err != nil {
			err = fmt.Errorf(msgErrMarshalResponse, err)
			return c.buildError(err, statusForbidden, bodyRequest)
		}
		return c.buildError(bodyResponseError(bodyResult), statusForbidden, bodyRequest)

	default:
		return nil
	}
}

func (c cardsTransactions) readAndValidateResponse(response *http.Response, bodyRequest []byte) ([]byte, error) {
	defer response.Body.Close()
	body, err := c.ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, c.buildError(err, statusBadGateway, bodyRequest)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, c.buildError(bodyResponseError(body), statusCodeCardsTransactions(response.StatusCode), bodyRequest)
	}

	return body, nil
}

func (c cardsTransactions) getBaseURL(providerSite string) url.URL {
	baseURL := c.config.App.BaseURL
	if _, found := c.config.App.ProvidersWithLegacyURL[providerSite]; found {
		baseURL = c.config.App.BaseURLLegacy
	}

	return baseURL
}

func (c cardsTransactions) fillHeader(transaction domain.TransactionData, headerValueXClientID string) http.Header {
	header := http.Header{}
	header.Add(headerKeyClientID, headerValueXClientID)
	header.Add(headerKeySiteID, transaction.SiteID)
	header.Add(headerKeyEnvironment, transaction.Environment)
	header.Add(headerKeyProvider, transaction.Provider.ID)
	header.Add(headerKeyProviderSiteID, fmt.Sprintf(`%s-%s`, strings.ToLower(transaction.Provider.ID), strings.ToLower(transaction.SiteID)))

	if len(c.config.App.BetaScope) > 0 {
		header.Add(headerKeyBetaScope, c.config.App.BetaScope)
	}

	return header
}

func fillRequest(transactionData domain.TransactionData) requestReversal {
	return requestReversal{
		ID: fmt.Sprint(typeSmartReverse, "_", transactionData.Operation.Authorization.ID),
		Operation: operation{
			Type:                 setValueStringOrNil(transactionData.Operation.SubType),
			TransmissionDatetime: transactionData.Operation.TransmissionDatetime,
			IsAdvice:             transactionData.Operation.IsAdvice,
			Installments:         transactionData.Operation.Installments,
			AcquirerCode:         setValueStringOrNil(transactionData.Operation.AcquirerCode),
			Stan:                 setValueIntOrNil(transactionData.Operation.Stan),
			Authorization: authorization{
				ID:                   transactionData.Operation.Authorization.ID,
				Stan:                 setValueIntOrNil(transactionData.Operation.Authorization.Stan),
				TransmissionDatetime: transactionData.Operation.TransmissionDatetime,
				CardAcceptor: cardAcceptor{
					Terminal: setValueStringOrNil(transactionData.Operation.Authorization.CardAcceptor.Terminal),
				},
			},
			Transaction: transaction{
				Amount:        transactionData.Operation.Transaction.Amount,
				Currency:      transactionData.Operation.Transaction.Currency,
				DecimalDigits: transactionData.Operation.Transaction.DecimalDigits,
				TotalAmount:   transactionData.Operation.Transaction.TotalAmount,
			},
			Billing: billing{
				Amount:        transactionData.Operation.Billing.Amount,
				Currency:      transactionData.Operation.Billing.Currency,
				DecimalDigits: transactionData.Operation.Billing.DecimalDigits,
				TotalAmount:   transactionData.Operation.Billing.TotalAmount,
				Conversion:    conversion(transactionData.Operation.Billing.Conversion),
			},
			Settlement: settlement{
				Amount:        transactionData.Operation.Settlement.Amount,
				Currency:      transactionData.Operation.Settlement.Currency,
				DecimalDigits: transactionData.Operation.Settlement.DecimalDigits,
				Conversion:    conversion(transactionData.Operation.Settlement.Conversion),
			},
			Card: card{
				NumberId: transactionData.Operation.Card.NumberID,
				Country:  transactionData.Operation.Card.Country,
			},
		},
		Provider: provider(transactionData.Provider),
		Options: options{
			ReversalType: typeSmartReverse,
		},
	}
}

func setValueStringOrNil(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return &v
}

func setValueIntOrNil(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func fillOutput(requestReversal requestReversal) *domain.ReversalOutput {
	return &domain.ReversalOutput{
		ReverseID: requestReversal.ID,
	}
}
