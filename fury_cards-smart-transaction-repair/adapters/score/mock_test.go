package score

import (
	"context"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mock_GetScore(t *testing.T) {
	m := NewMock()
	output, err := m.GetScore(context.TODO(), domain.ScoreInput{})

	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
