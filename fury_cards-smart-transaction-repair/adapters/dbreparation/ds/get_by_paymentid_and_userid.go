package ds

import (
	"github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

func (d *dsTransactionRepair) GetByPaymentIDAndUserID(paymentID string, userID int64) (*storage.ReparationOutput, error) {
	queryBuilder := d.query.And(
		d.query.Eq("payment_id", paymentID),
		d.query.And(
			d.query.Eq("user_id", userID),
		),
	)

	transactionRepair := dbreparation.TransactionRepairDB{}

	err := d.repo.GetDocumentByQueryBuilder(queryBuilder, &transactionRepair)
	if err != nil {
		if err == dsclient.ErrDocumentNotFound {
			return nil, buildErrorNotFound(paymentID, err)
		}
		return nil, buildErrorBadGateway(err)
	}

	return dbreparation.ToTransactionRepairOutput(transactionRepair), nil
}
