package kvs

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	kvs "github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_blockedlistRepository_Get(t *testing.T) {
	mockKvsClient := mockKvsClient{
		GetValueStub: func(ctx context.Context, key string, value interface{}) error {
			err := json.Unmarshal(_blockedUser, value)
			assert.NoError(t, err)
			assert.Equal(t, "123", key)

			return nil
		},
	}

	kvsRepo := New(mockKvsClient)

	output, err := kvsRepo.Get(context.TODO(), "123")

	outputExpected := mockStorageBlockedUser(t)

	assert.NoError(t, err)
	assert.Equal(t, &outputExpected, output)
}

func Test_blockedlistRepository_Get_Error_NotFound(t *testing.T) {
	mockKvsClient := mockKvsClient{
		GetValueStub: func(ctx context.Context, key string, value interface{}) error {
			return kvs.ErrNotFound
		},
	}

	kvsRepo := New(mockKvsClient)

	output, err := kvsRepo.Get(context.TODO(), "123")

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("user not found | ID: 123"), "blockeduser database failure", false)

	assert.Nil(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_blockedlistRepository_Get_Error_Unknown(t *testing.T) {
	mockKvsClient := mockKvsClient{
		GetValueStub: func(ctx context.Context, key string, value interface{}) error {
			return errors.New("some unknown error")
		},
	}

	kvsRepo := New(mockKvsClient)

	output, err := kvsRepo.Get(context.TODO(), "123")

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"), "blockeduser database failure", false)

	assert.Nil(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
