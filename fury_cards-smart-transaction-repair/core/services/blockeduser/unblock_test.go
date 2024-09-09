package blockeduser

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_blocked_UnblockUser_Success(t *testing.T) {
	ctx := context.TODO()

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			return &storage.BlockedUser{
				KycIdentificationID: "abc",
				UnlockedAt:          nil,
				SiteID:              "MLM",
				CreatedAt:           time.Now(),
				LastBlockedAt:       time.Now(),
				CountCaptures:       1,
				CapturedRepairs: []storage.CapturedRepairs{
					{
						AuthorizationID: "auth_01",
						UserID:          123,
						PaymentID:       "pay_01",
						TypeRepair:      domain.TypeReverse,
						RepairedAt:      time.Now(),
						CreatedAt:       time.Now(),
					},
				},
			}, nil
		},
		SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
			assert.NotNil(t, input.UnlockedAt)
			utilstest.AssertFieldsNoEmptyFromStruct(t, input)
			utilstest.AssertFieldsNoEmptyFromStruct(t, input.CapturedRepairs[0])

			return nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	err := service.UnblockUser(ctx, "123")

	assert.NoError(t, err)
}

func Test_blocked_UnblockUser_Error_Save(t *testing.T) {
	ctx := context.TODO()

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			return &storage.BlockedUser{
				KycIdentificationID: "abc",
				UnlockedAt:          nil,
				SiteID:              "MLM",
				CreatedAt:           time.Now(),
				LastBlockedAt:       time.Now(),
				CountCaptures:       1,
				CapturedRepairs: []storage.CapturedRepairs{
					{
						AuthorizationID: "auth_01",
						UserID:          123,
						PaymentID:       "pay_01",
						TypeRepair:      domain.TypeReverse,
						RepairedAt:      time.Now(),
						CreatedAt:       time.Now(),
					},
				},
			}, nil
		},
		SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
			return gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	err := service.UnblockUser(ctx, "123")

	errExpected := gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)

	utilstest.AssertGnierrorsExpected(t, errExpected, err)
}

func Test_blocked_UnblockUser_Error_NotFound(t *testing.T) {
	ctx := context.TODO()

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "blockedlist database failure", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	err := service.UnblockUser(ctx, "123")

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("not found"), "blockedlist database failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_blocked_UnblockUser_Error_Unknown(t *testing.T) {
	ctx := context.TODO()

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist database failure", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	err := service.UnblockUser(ctx, "123")

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"), "blockedlist database failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
