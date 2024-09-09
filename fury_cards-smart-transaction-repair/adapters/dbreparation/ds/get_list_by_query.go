package ds

import (
	"encoding/json"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

const (
	format = "2006-01-02"
)

func (d *dsTransactionRepair) GetListByUserIDAndCreationPeriod(
	userID int64,
	gteCreatedAt,
	lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {

	queryBuilder := d.query.And(
		d.query.Eq("user_id", userID),
		d.query.DateRange("created_at").Gte(gteCreatedAt.Format(format)),
		d.query.DateRange("created_at").Lte(lteCreatedAt.Format(format)),
	)

	return d.GetList(queryBuilder)
}

func (d *dsTransactionRepair) GetListByKycIentificationAndCreationPeriod(
	kycIdentificationID string,
	gteCreatedAt,
	lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {

	queryBuilder := d.query.And(
		d.query.Eq("kyc_identification_id", kycIdentificationID),
		d.query.DateRange("created_at").Gte(gteCreatedAt.Format(format)),
		d.query.DateRange("created_at").Lte(lteCreatedAt.Format(format)),
	)

	return d.GetList(queryBuilder)
}

func (d *dsTransactionRepair) GetList(queryBuilder dsclient.QueryBuilder) ([]storage.ReparationOutput, error) {
	documents, err := d.repo.GetListByQueryBuilder(queryBuilder)
	if err != nil {
		return nil, buildErrorBadGateway(err)
	}

	listOutput := []storage.ReparationOutput{}

	var transactionsRepairDB []dbreparation.TransactionRepairDB
	err = json.Unmarshal(documents, &transactionsRepairDB)
	if err != nil {
		return nil, buildErrorBadGateway(err)
	}

	for _, transactionRepairDB := range transactionsRepairDB {
		listOutput = append(listOutput, *dbreparation.ToTransactionRepairOutput(transactionRepairDB))
	}

	return listOutput, nil
}
