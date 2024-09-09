package ds

import (
	"encoding/json"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

// Max default value for the fury consumer bulk delivery mode.
const defaultSize = 20

func (d *dsTransactionRepair) GetListByAuthorizationIDs(authorizationIDs ...string) ([]storage.ReparationOutput, error) {

	queryBuilder := d.query.Ids(authorizationIDs)

	documents, err := d.repo.GetListByQueryBuilder(queryBuilder, dsclient.WithSize(defaultSize))
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
