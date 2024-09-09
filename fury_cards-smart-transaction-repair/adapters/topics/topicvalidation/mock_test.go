package topicvalidation

import (
	"context"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

func Test_mock(t *testing.T) {
	mock := NewMock()

	field.NewWrappedFields()
	err := mock.Publish(context.Background(), storage.ValidationResult{}, field.NewWrappedFields())
	assert.NoError(t, err)
}
