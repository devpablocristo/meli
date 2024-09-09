package payment

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	hoursOneDay = 24
)

func FillOutput(data paymentTransaction) domain.TransactionOutput {
	return domain.TransactionOutput{
		TransactionData: fillTransactionData(*data.Payment),
	}
}

func fillTransactionData(payment payment) *domain.TransactionData {
	return &domain.TransactionData{
		Operation: domain.Operation{
			SubType:                payment.DebitTransaction.Operation.SubType,
			CreationDatetime:       payment.DebitTransaction.Operation.CreationDatetime,
			TransmissionDatetime:   payment.DebitTransaction.Operation.TransmissionDatetime,
			IsAdvice:               payment.DebitTransaction.Operation.IsAdvice,
			Installments:           payment.DebitTransaction.Operation.Installments,
			AcquirerCode:           payment.DebitTransaction.AcquirerCode,
			Stan:                   payment.DebitTransaction.Operation.Stan,
			AccountType:            payment.DebitTransaction.Operation.AccountType,
			AuthorizationIndicator: payment.DebitTransaction.Operation.AuthorizationIndicator,
			PaymentAge:             calculatePaymentAgeInDays(payment.DebitTransaction.Operation.CreationDatetime),
			IsInternational:        payment.DebitTransaction.Operation.IsInternational,
			Authorization: domain.Authorization{
				ID:                   payment.DebitTransaction.ID,
				Stan:                 payment.DebitTransaction.Operation.Stan,
				TransmissionDatetime: payment.DebitTransaction.Operation.TransmissionDatetime,
				CardAcceptor:         domain.CardAcceptor(payment.DebitTransaction.CardAcceptor),
			},
			Transaction: domain.Transaction{
				Amount:        payment.DebitTransaction.Operation.Transaction.Amount,
				Currency:      payment.DebitTransaction.Operation.Transaction.Currency,
				DecimalDigits: payment.DebitTransaction.Operation.Transaction.DecimalDigits,
				TotalAmount:   payment.DebitTransaction.Operation.Transaction.TotalAmount,
			},
			Card: domain.Card{
				NumberID: payment.DebitTransaction.Operation.Card.NumberID,
				Country:  payment.DebitTransaction.Operation.Card.Country,
				IssuerAccount: domain.IssuerAccount{
					CardBusinessMode: payment.DebitTransaction.Operation.Card.IssuerAccount.Card.BusinessMode,
				},
				TokenInfo: domain.TokenInfo{
					TokenWalletID: payment.DebitTransaction.Operation.Card.TokenInfo.TokenWalletID,
				},
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
			},
		},
		Provider: domain.Provider{
			ID: payment.DebitTransaction.Provider.ID,
		},
		Environment:        payment.DebitTransaction.Environment,
		SiteID:             payment.SiteID,
		StatusDetail:       payment.StatusDetail,
		PayerID:            payment.PayerID,
		Type:               payment.DebitTransaction.Type,
		IsRecurringPayment: payment.DebitTransaction.IsRecurringPayment,
	}
}

func calculatePaymentAgeInDays(creationDatetime time.Time) int {
	return int(time.Now().UTC().Sub(creationDatetime).Hours()) / hoursOneDay
}
