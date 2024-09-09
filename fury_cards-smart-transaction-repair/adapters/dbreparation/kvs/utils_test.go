package kvs

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/transaction-reparation.json
var _transactionReparation []byte

//go:embed testdata/transaction-reparation_domain.json
var _transactionReparationDomain []byte

// KvsService
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

// LogService
type mockLogService struct {
	log.LogService
	WarnStub func(c context.Context, msg string, fields ...log.Field)
}

func (m mockLogService) Warn(c context.Context, msg string, fields ...log.Field) {
	m.WarnStub(c, msg, fields...)
}

func mockReparation(t *testing.T) (input storage.ReparationInput) {
	err := json.Unmarshal(_transactionReparationDomain, &input)
	assert.NoError(t, err)
	return
}
