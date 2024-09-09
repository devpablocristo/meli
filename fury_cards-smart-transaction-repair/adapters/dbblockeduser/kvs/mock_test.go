package kvs

import (
	"context"
	_ "embed"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/stretchr/testify/assert"
)

func Test_mock(t *testing.T) {
	mock := NewMock()
	err := mock.Save(context.TODO(), "123", storage.BlockedUser{})
	assert.NoError(t, err)

	output, err := mock.Get(context.TODO(), "123")
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
