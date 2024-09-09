package ds

import (
	_ "embed"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/stretchr/testify/assert"
)

func Test_mock_GetListByAuthorizationIDs(t *testing.T) {
	m := NewMock()

	t.Run("GetListByAuthorizationIDs", func(t *testing.T) {
		blockedUsers, err := m.GetListByAuthorizationIDs("MLM", storage.AuthorizationIDsSearch{"auth_01"})

		assert.NoError(t, err)
		assert.NotNil(t, blockedUsers)
	})
}
