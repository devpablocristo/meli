package graphqltransaction

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

const (
	contentType = "application/json"
	pathURL     = "/cards/transactions/search"
	hoursOneDay = 24
)

func (t transactionSearch) GetTransactionByPaymentID(
	ctx context.Context,
	input *domain.TransactionInput,
	wrapFields *field.WrappedFields) (*domain.TransactionOutput, error) {

	request := requestTransaction{
		Query: queryPretty,
		Variables: Variables{
			ID: input.PaymentID,
		},
	}

	body, err := t.json.Marshal(request)
	if err != nil {
		return nil, buildError(err, http.StatusInternalServerError)
	}

	url := rest.CreateURL(t.config.BaseURL, pathURL)

	response, err := t.httpClient.Do(ctx, "POST", url, contentType, body, nil)
	if err != nil {
		return nil, buildError(err, http.StatusBadGateway)
	}

	body, err = t.readAndValidateResponse(response)
	if err != nil {
		return nil, err
	}

	wrapFields.Fields.AddAt(string(domain.KeyFieldBodyGraphqlTransaction), strings.ReplaceAll(string(body), "\\", ""), log.ErrorLevel)

	var transaction responseTransaction
	err = t.json.Unmarshal(body, &transaction)
	if err != nil {
		return nil, buildError(err, http.StatusInternalServerError)
	}

	err = t.validateTransactionDataResponse(transaction)
	if err != nil {
		return nil, err
	}

	return fillTransactionOutput(*transaction.Data.Payment), nil
}

func (t transactionSearch) readAndValidateResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	body, err := t.ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, buildError(err, http.StatusBadGateway)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, buildError(bodyResponseError(body), response.StatusCode)
	}

	return body, nil
}

func fillTransactionOutput(payment payment) *domain.TransactionOutput {
	return &domain.TransactionOutput{
		TransactionData: &domain.TransactionData{
			Operation: domain.Operation{
				SubType:              payment.DebitTransaction.Operation.SubType,
				CreationDatetime:     payment.DebitTransaction.Operation.CreationDatetime,
				TransmissionDatetime: payment.DebitTransaction.Operation.TransmissionDatetime,
				IsAdvice:             payment.DebitTransaction.Operation.IsAdvice,
				Installments:         payment.DebitTransaction.Operation.Installments,
				AcquirerCode:         payment.DebitTransaction.AcquirerCode,
				Stan:                 payment.DebitTransaction.Operation.Stan,
				Authorization: domain.Authorization{
					ID:                   payment.DebitTransaction.ID,
					Stan:                 payment.DebitTransaction.Operation.Stan,
					TransmissionDatetime: payment.DebitTransaction.Operation.TransmissionDatetime,
					CardAcceptor: domain.CardAcceptor{
						Terminal: payment.DebitTransaction.CardAcceptor.Terminal,
						Name:     payment.DebitTransaction.CardAcceptor.Name,
					},
				},
				Transaction: domain.Transaction{
					Amount:        payment.DebitTransaction.Operation.Transaction.Amount,
					Currency:      payment.DebitTransaction.Operation.Transaction.Currency,
					DecimalDigits: payment.DebitTransaction.Operation.Transaction.DecimalDigits,
					TotalAmount:   payment.DebitTransaction.Operation.Transaction.TotalAmount,
				},
				Card: domain.Card{
					NumberID: payment.DebitTransaction.Operation.Card.NumberId,
					Country:  payment.DebitTransaction.Operation.Card.Country,
				},
				Billing: domain.Billing{
					Amount:        payment.DebitTransaction.Operation.Billing.Amount,
					Currency:      payment.DebitTransaction.Operation.Billing.Currency,
					DecimalDigits: payment.DebitTransaction.Operation.Billing.DecimalDigits,
					TotalAmount:   payment.DebitTransaction.Operation.Billing.TotalAmount,
					Conversion:    domain.Conversion(payment.DebitTransaction.Operation.Billing.Conversion),
				},
				Settlement: domain.Settlement{
					Amount:        payment.DebitTransaction.Operation.Settlement.Amount,
					Currency:      payment.DebitTransaction.Operation.Settlement.Currency,
					DecimalDigits: payment.DebitTransaction.Operation.Settlement.DecimalDigits,
					TotalAmount:   payment.DebitTransaction.Operation.Settlement.TotalAmount,
					Conversion:    domain.Conversion(payment.DebitTransaction.Operation.Settlement.Conversion),
				},
				PaymentAge: calculatePaymentAgeInDays(payment.DebitTransaction.Operation.CreationDatetime),
			},
			Provider: domain.Provider{
				ID: payment.DebitTransaction.Provider.ID,
			},
			Environment:  payment.DebitTransaction.Environment,
			SiteID:       payment.SiteID,
			StatusDetail: payment.StatusDetail,
			PayerID:      payment.PayerID,
			Type:         payment.DebitTransaction.Type,
		},
	}
}

func (t transactionSearch) validateTransactionDataResponse(transaction responseTransaction) error {
	if transaction.Data.Payment == nil && len(transaction.Err) > 0 {
		b, err := t.json.Marshal(transaction.Err)
		if err != nil {
			return buildError(err, http.StatusInternalServerError)
		}

		return buildError(errors.New(string(b)), http.StatusBadGateway)
	}
	return nil
}

func calculatePaymentAgeInDays(creationDatetime time.Time) int {
	return int(time.Now().UTC().Sub(creationDatetime).Hours()) / hoursOneDay
}
