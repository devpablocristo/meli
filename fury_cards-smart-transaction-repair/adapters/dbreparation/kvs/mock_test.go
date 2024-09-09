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
	err := mock.Save(context.TODO(), storage.ReparationInput{})
	assert.NoError(t, err)

	output, err := mock.Get(context.TODO(), "auth_1")
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
