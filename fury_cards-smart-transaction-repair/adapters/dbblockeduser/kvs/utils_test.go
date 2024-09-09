package kvs

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	kvsclient "github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/blockeduser.json
var _blockedUser []byte

//go:embed testdata/blockeduser-domain.json
var _blockedUserStruct []byte

type mockKvsClient struct {
	kvsclient.Service
	PutStub      func(ctx context.Context, key string, value interface{}) error
	GetValueStub func(ctx context.Context, key string, value interface{}) error
}

func (m mockKvsClient) Put(ctx context.Context, key string, value interface{}) error {
	return m.PutStub(ctx, key, value)
}

func (m mockKvsClient) GetValue(ctx context.Context, key string, value interface{}) error {
	return m.GetValueStub(ctx, key, value)
}

func mockStorageBlockedUser(t *testing.T) (blockedUser storage.BlockedUser) {
	err := json.Unmarshal(_blockedUserStruct, &blockedUser)
	assert.NoError(t, err)
	return
}
