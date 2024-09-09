package kvs

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

//go:embed testdata/blockeduser.json
var _blockedUserMock []byte

type mock struct{}

func (m mock) Save(ctx context.Context, kycIdentificationID string, blockedUser storage.BlockedUser) error {
	return nil
}

func (m mock) Get(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
	var output *storage.BlockedUser
	err := json.Unmarshal(_blockedUserMock, &output)

	return output, err
}

func NewMock() mock {
	return mock{}
}
