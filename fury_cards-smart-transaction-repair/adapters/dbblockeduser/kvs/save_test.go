package kvs

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbblockeduser"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_blockedlistRepository_Save(t *testing.T) {
	ctx := context.TODO()

	input := mockStorageBlockedUser(t)

	mockKvsClient := mockKvsClient{
		PutStub: func(ctx context.Context, key string, value interface{}) error {
			blockedUserDB := value.(dbblockeduser.BlockedUserDB)

			b, err := json.Marshal(blockedUserDB)
			assert.NoError(t, err)

			assert.Equal(t, "123", key)
			assert.JSONEq(t, string(_blockedUser), string(b))

			return nil
		},
	}

	kvsRepo := New(mockKvsClient)

	err := kvsRepo.Save(ctx, "123", input)

	assert.Nil(t, err)
}

func Test_blockedlistRepository_Save_Error_Unknown(t *testing.T) {
	ctx := context.TODO()

	input := mockStorageBlockedUser(t)

	mockKvsClient := mockKvsClient{
		PutStub: func(ctx context.Context, key string, value interface{}) error {
			return errors.New("some unknown error")
		},
	}

	kvsRepo := New(mockKvsClient)

	err := kvsRepo.Save(ctx, "123", input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"), "blockeduser database failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
