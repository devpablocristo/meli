package paymentcards

import "github.com/melisource/cards-smart-transaction-repair/core/domain"

func FillOutput(data TransactionWalletSearch) *domain.TransactionWalletOutput {
	return &domain.TransactionWalletOutput{
		TransactionData: fillTransactionData(*data.Payment),
		WalletData:      fillWalletData(*data.Wallet),
	}
}

func fillTransactionData(payment payment) *domain.TransactionData {
	return &domain.TransactionData{
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
				NumberID: payment.DebitTransaction.Operation.Card.NumberID,
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
			},
		},
		Provider: domain.Provider{
			ID: payment.DebitTransaction.Provider.ID,
		},
		Environment:  payment.DebitTransaction.Environment,
		SiteID:       payment.SiteID,
		StatusDetail: payment.StatusDetail,
		PayerID:      payment.PayerID,
	}
}

func fillWalletData(wallet wallet) *domain.WalletData {
	return &domain.WalletData{
		ID:       wallet.ID,
		CardsIDs: wallet.CardsIDs,
		Cards:    fillCards(wallet.Cards),
	}
}

func fillCards(cards []cardWallet) []domain.CardWallet {
	cardsWallet := []domain.CardWallet{}
	for _, card := range cards {
		cardsWallet = append(cardsWallet, domain.CardWallet{
			ID:             card.ID,
			BusinessMode:   card.BusinessMode,
			IsTest:         card.IsTest,
			Holder:         domain.CardHolder(card.Holder),
			IssuerAccounts: fillIssuerAccounts(card.IssuerAccounts),
		})
	}

	return cardsWallet
}

func fillIssuerAccounts(issuerAccounts []issuerAccounts) []domain.CardIssuerAccounts {
	cardIssuerAccounts := []domain.CardIssuerAccounts{}
	for _, accounts := range issuerAccounts {
		cardIssuerAccounts = append(cardIssuerAccounts, domain.CardIssuerAccounts{
			ID: accounts.ID,
		})
	}

	return cardIssuerAccounts
}
