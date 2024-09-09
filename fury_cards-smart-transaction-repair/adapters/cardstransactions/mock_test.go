package cardstransactions

import (
	"context"
	_ "embed"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/stretchr/testify/assert"
)

func Test_mock_Reverse(t *testing.T) {
	m := NewMock()
	output, err := m.Reverse(context.TODO(), domain.ReversalInput{})

	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
