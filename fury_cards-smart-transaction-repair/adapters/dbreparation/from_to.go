package dbreparation

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

func ToTransactionRepairOutput(transactionRepair TransactionRepairDB) *storage.ReparationOutput {
	return &storage.ReparationOutput{
		BaseReparation: storage.BaseReparation{
			AuthorizationID:     transactionRepair.AuthorizationID,
			PaymentID:           transactionRepair.PaymentID,
			KycIdentificationID: transactionRepair.KycIdentificationID,
			UserID:              transactionRepair.UserID,
			TransactionRepairID: transactionRepair.TransactionRepairID,
			SiteID:              transactionRepair.SiteID,
			Type:                transactionRepair.Type,
			CreatedAt:           transactionRepair.CreatedAt,
			FaqID:               transactionRepair.FaqID,
			Requester:           domain.Requester(transactionRepair.Requester),
			MonetaryTransaction: toMonetaryTransaction(transactionRepair.MonetaryTransaction),
		},
	}
}

func toMonetaryTransaction(monetaryTransaction MonetaryTransactionDB) storage.MonetaryTransaction {
	return storage.MonetaryTransaction{
		Billing:    storage.MonetaryValue(monetaryTransaction.Billing),
		Settlement: storage.MonetaryValue(monetaryTransaction.Settlement),
	}
}
