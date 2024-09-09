package ds

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_mock(t *testing.T) {
	m := NewMock()

	t.Run("GetByPaymentIDAndUserID", func(t *testing.T) {
		repair, err := m.GetByPaymentIDAndUserID("pay01", 123)

		assert.NoError(t, err)
		assert.NotNil(t, repair)
	})

	t.Run("GetListByUserIDAndCreationPeriod", func(t *testing.T) {
		list, err := m.GetListByUserIDAndCreationPeriod(123, time.Now(), time.Now())

		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("GetListByKycIentificationAndCreationPeriod", func(t *testing.T) {
		list, err := m.GetListByKycIentificationAndCreationPeriod("abc", time.Now(), time.Now())

		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("GetListByAuthorizationIDs", func(t *testing.T) {
		list, err := m.GetListByAuthorizationIDs("auth_1", "auth_2")

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, "reverse_smart_auth_1", list[0].TransactionRepairID)
		assert.Equal(t, "reverse_smart_auth_2", list[1].TransactionRepairID)
	})
}
