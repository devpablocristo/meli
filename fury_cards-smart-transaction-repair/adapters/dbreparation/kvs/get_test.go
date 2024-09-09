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

func Test_Get(t *testing.T) {
	mockKvsClient := mockKvsClient{
		GetValueStub: func(ctx context.Context, key string, value interface{}) error {
			err := json.Unmarshal(_transactionReparation, value)
			assert.NoError(t, err)
			assert.Equal(t, "auth_1", key)

			return nil
		},
	}

	kvsRepo := New(mockKvsClient, mockLogService{})

	output, err := kvsRepo.Get(context.TODO(), "auth_1")

	outputExpected := mockReparation(t)

	assert.NoError(t, err)
	assert.Equal(t, output.BaseReparation, outputExpected.BaseReparation)
}

func Test_Not_Get_Error_KeyNotFound(t *testing.T) {
	mockKvsClient := mockKvsClient{
		GetValueStub: func(ctx context.Context, key string, value interface{}) error {
			return kvs.ErrNotFound
		},
	}

	kvsRepo := New(mockKvsClient, mockLogService{})

	output, err := kvsRepo.Get(context.TODO(), "auth_1")

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("transaction repair not found | ID: auth_1"), "transaction repair database failure", false)

	assert.Nil(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_Not_Get_Error_Unknow(t *testing.T) {
	mockKvsClient := mockKvsClient{
		GetValueStub: func(ctx context.Context, key string, value interface{}) error {
			return errors.New("some error unknow")
		},
	}

	mockLogService := mockLogService{}

	kvsRepo := New(mockKvsClient, mockLogService)

	output, err := kvsRepo.Get(context.TODO(), "auth_1")

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some error unknow"), "transaction repair database failure", false)

	assert.Nil(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_Not_Get_Error_Validation(t *testing.T) {
	kvsRepo := New(mockKvsClient{}, mockLogService{})

	output, err := kvsRepo.Get(context.TODO(), "")

	errorExpected := gnierrors.Wrap(domain.BadRequest,
		errors.New("authorizationID cannot be null"), "transaction repair database failure", false)

	assert.Nil(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
