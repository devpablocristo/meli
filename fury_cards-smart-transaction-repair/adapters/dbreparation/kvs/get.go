package kvs

import (
	"context"
	"errors"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	kvsclient "github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
)

func (k kvsRepository) Get(ctx context.Context, authorizationID string) (*storage.ReparationOutput, error) {
	if err := validateGet(authorizationID); err != nil {
		return nil, err
	}
	var transactionRepair dbreparation.TransactionRepairDB
	err := k.repo.GetValue(ctx, authorizationID, &transactionRepair)
	if err != nil {
		if err == kvsclient.ErrNotFound {
			return nil, buildErrorNotFound(authorizationID)
		}
		return nil, buildErrorBadGateway(err)
	}
	return dbreparation.ToTransactionRepairOutput(transactionRepair), nil
}

func validateGet(authorizationID string) error {
	if authorizationID == "" {
		return buildErrorBadRequest(errors.New("authorizationID cannot be null"))
	}
	return nil
}
