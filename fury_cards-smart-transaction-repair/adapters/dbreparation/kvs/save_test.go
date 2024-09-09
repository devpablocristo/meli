package kvs

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbreparation"
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Save(t *testing.T) {
	ctx := context.TODO()

	input := mockReparation(t)

	mockKvsClient := mockKvsClient{
		PutStub: func(ctx context.Context, key string, value interface{}) error {
			reparation := value.(*dbreparation.TransactionRepairDB)
			bReparation, err := json.Marshal(reparation)
			assert.NoError(t, err)

			assert.Equal(t, "auth_1", key)
			assert.JSONEq(t, string(_transactionReparation), string(bReparation))

			return nil
		},
	}

	kvsRepo := New(mockKvsClient, mockLogService{})

	err := kvsRepo.Save(ctx, input)

	assert.Nil(t, err)
}

func Test_Save_Error_Unknow(t *testing.T) {
	ctx := context.TODO()

	input := mockReparation(t)

	mockKvsClient := mockKvsClient{
		PutStub: func(ctx context.Context, key string, value interface{}) error {
			return errors.New("some error unknow")
		},
	}

	kvsRepo := New(mockKvsClient, mockLogService{})

	err := kvsRepo.Save(ctx, input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some error unknow"), "transaction repair database failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_Save_Error_Validation(t *testing.T) {
	ctx := context.TODO()

	input := storage.ReparationInput{}

	mockLogService := mockLogService{
		WarnStub: func(c context.Context, msg string, fields ...log.Field) {},
	}

	kvsRepo := New(mockKvsClient{}, mockLogService)

	err := kvsRepo.Save(ctx, input)

	errorExpected := gnierrors.Wrap(domain.BadRequest,
		errors.New("required fields: AuthorizationID; PaymentID; KycIdentificationID; UserID;"), "transaction repair database failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
