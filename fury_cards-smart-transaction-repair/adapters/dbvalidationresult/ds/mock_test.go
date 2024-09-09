package ds

import (
	"context"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/stretchr/testify/assert"
)

func Test_mock_Save(t *testing.T) {
	m := NewMock()

	err := m.Save(context.TODO(), "12", storage.ValidationResult{})

	assert.NoError(t, err)
}
