package policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mock(t *testing.T) {
	mock := NewMock()

	res, err := mock.EvaluateWithUser(context.Background(), 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
