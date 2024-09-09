package kvs

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

//go:embed testdata/transaction-reparation_domain.json
var _transactionReparationMock []byte

type mock struct{}

func (m mock) Save(ctx context.Context, transactionRepair storage.ReparationInput) error {
	return nil
}

func (m mock) Get(ctx context.Context, authorizationID string) (*storage.ReparationOutput, error) {
	var output *storage.ReparationOutput
	err := json.Unmarshal(_transactionReparationMock, &output)

	return output, err
}

func NewMock() mock {
	return mock{}
}
