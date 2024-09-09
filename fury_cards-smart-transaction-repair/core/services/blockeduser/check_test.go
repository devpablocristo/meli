package blockeduser

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_blocked_CheckUserIsBlocked_True(t *testing.T) {
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
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	isBlocked, err := service.CheckUserIsBlocked(ctx, domain.BlockedUserInput{KycIdentificationID: "123"}, field.NewWrappedFields())

	assert.NoError(t, err)
	assert.True(t, isBlocked)
}

func Test_blocked_CheckUserIsBlocked_False(t *testing.T) {
	ctx := context.TODO()

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			now := new(time.Time)
			*now = time.Now().UTC()
			return &storage.BlockedUser{
				KycIdentificationID: "abc",
				UnlockedAt:          now,
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
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	isBlocked, err := service.CheckUserIsBlocked(ctx, domain.BlockedUserInput{KycIdentificationID: "123"}, field.NewWrappedFields())

	assert.NoError(t, err)
	assert.False(t, isBlocked)
}

func Test_blocked_CheckUserIsBlocked_False_NotFound(t *testing.T) {
	ctx := context.TODO()

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "blockedlist database failure", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	isBlocked, err := service.CheckUserIsBlocked(ctx, domain.BlockedUserInput{KycIdentificationID: "123"}, field.NewWrappedFields())

	assert.False(t, isBlocked)
	assert.NoError(t, err)
}

func Test_blocked_CheckUserIsBlocked_Error_Unknown(t *testing.T) {
	ctx := context.TODO()

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist database failure", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo)

	isBlocked, err := service.CheckUserIsBlocked(ctx, domain.BlockedUserInput{KycIdentificationID: "123"}, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"), "blockedlist database failure", false)

	assert.False(t, isBlocked)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
