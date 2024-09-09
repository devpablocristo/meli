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

func Test_blocked_BlockUserIfReparationExistsInBulk(t *testing.T) {
	t.Run("repair(s) exist: should block user(s) first time", func(t *testing.T) {
		ctx := context.TODO()

		mockReparationSearchRepo := mockReparationSearchRepo{
			GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
				return []storage.ReparationOutput{
					mockTransactionRepair01(),
					mockTransactionRepair02(),
				}, nil
			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
				assert.Nil(t, input.UnlockedAt)
				utilstest.AssertFieldsNoEmptyFromStruct(t, input)

				switch kycIdentificationID {
				case "abc":
					utilstest.AssertStructEqual(t, mockListFirstBlockedUser()[0], input)
				case "def":
					utilstest.AssertStructEqual(t, mockListFirstBlockedUser()[1], input)
				default:
					panic(kycIdentificationID)
				}

				return nil
			},
			GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
				return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "blockedlist database failure", false)
			},
		}

		service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

		output, err := service.BlockUserIfReparationExistsInBulk(
			ctx,
			domain.TransactionsNewsFeedBulkInput{
				SiteID:                 "MLM",
				AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_01": "cap01", "auth_02": "cap02"}},
			field.NewWrappedFields(),
		)

		assert.NoError(t, err)
		assert.Len(t, output.AuthorizationIDsWithErr, 0)
	})
}

func Test_blocked_BlockUserIfReparationExistsInBulk_Block_Error_Save(t *testing.T) {
	t.Run("repair(s) exist: should block user(s) first time, but fail", func(t *testing.T) {
		ctx := context.TODO()

		mockReparationSearchRepo := mockReparationSearchRepo{
			GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
				return []storage.ReparationOutput{
					mockTransactionRepair01(),
					mockTransactionRepair02(),
				}, nil

			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
				return gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)
			},
			GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
				return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "blockedlist database failure", false)
			},
		}

		service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

		output, err := service.BlockUserIfReparationExistsInBulk(
			ctx,
			domain.TransactionsNewsFeedBulkInput{
				SiteID:                 "MLM",
				AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_01": "cap01", "auth_02": "cap02"}},
			field.NewWrappedFields(),
		)

		assert.NoError(t, err)

		errExpected := gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)
		utilstest.AssertGnierrorsExpected(t, errExpected, output.AuthorizationIDsWithErr["auth_01"])
		utilstest.AssertGnierrorsExpected(t, errExpected, output.AuthorizationIDsWithErr["auth_02"])
	})
}

func Test_blocked_BlockUserIfReparationExistsInBulk_Already_Blocked(t *testing.T) {
	t.Run("repair(s) exist, but there is already a block for the same authorization", func(t *testing.T) {

		ctx := context.TODO()

		mockReparationSearchRepo := mockReparationSearchRepo{
			GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
				return []storage.ReparationOutput{
					mockTransactionRepair01(),
					mockTransactionRepair02(),
				}, nil

			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
				blockeUser := &storage.BlockedUser{
					KycIdentificationID: kycIdentificationID,
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
				}

				switch kycIdentificationID {
				case "abc":
					return blockeUser, nil
				case "def":
					blockeUser.CapturedRepairs = []storage.CapturedRepairs{
						{
							AuthorizationID: "auth_02",
							PaymentID:       "pay_02",
							TypeRepair:      domain.TypeReverse,
							RepairedAt:      time.Now(),
							CreatedAt:       time.Now(),
							UserID:          456,
						},
					}
					return blockeUser, nil

				default:
					panic(kycIdentificationID)
				}
			},
		}

		service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

		output, err := service.BlockUserIfReparationExistsInBulk(
			ctx,
			domain.TransactionsNewsFeedBulkInput{
				SiteID:                 "MLM",
				AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_01": "cap01", "auth_02": "cap02"}},
			field.NewWrappedFields(),
		)

		assert.NoError(t, err)
		assert.Len(t, output.AuthorizationIDsWithErr, 0)
	})
}

func Test_blocked_BlockUserIfReparationExistsInBulk_No_Block(t *testing.T) {
	ctx := context.TODO()

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
			assert.Contains(t, authorizationIDs, "auth_01")
			assert.Contains(t, authorizationIDs, "auth_02")
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "transaction repair error", false)
		},
	}

	mockBlockedlistRepo := mockBlockedlistRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

	output, err := service.BlockUserIfReparationExistsInBulk(
		ctx,
		domain.TransactionsNewsFeedBulkInput{
			SiteID:                 "MLM",
			AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_01": "cap01", "auth_02": "cap02"}},
		field.NewWrappedFields(),
	)

	assert.NoError(t, err)
	assert.Len(t, output.AuthorizationIDsWithErr, 0)
}

func Test_blocked_BlockUserIfReparationExistsInBulk_Error_GetTransactionRepair_Unknown(t *testing.T) {
	ctx := context.TODO()

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
			assert.Contains(t, authorizationIDs, "auth_01")
			assert.Contains(t, authorizationIDs, "auth_02")
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "transaction repair database failure", false)
		},
	}

	mockBlockedlistRepo := mockBlockedlistRepo{}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

	output, err := service.BlockUserIfReparationExistsInBulk(
		ctx,
		domain.TransactionsNewsFeedBulkInput{
			SiteID:                 "MLM",
			AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_01": "cap01", "auth_02": "cap02"}},
		field.NewWrappedFields(),
	)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"), "transaction repair database failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Nil(t, output)
}

func Test_blocked_BlockUserIfReparationExistsInBulk_Block_And_Append(t *testing.T) {
	t.Run("reparation exists, it is the 2º capture with reparation done", func(t *testing.T) {
		ctx := context.TODO()

		mockReparationSearchRepo := mockReparationSearchRepo{
			GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
				return []storage.ReparationOutput{
					{
						BaseReparation: storage.BaseReparation{
							KycIdentificationID: "abc",
							AuthorizationID:     "auth_02",
							UserID:              123,
							SiteID:              "MLM",
							PaymentID:           "pay_02",
							TransactionRepairID: "rev_02",
							Type:                domain.TypeReverse,
							CreatedAt:           time.Date(2022, 10, 12, 0, 0, 0, 0, time.UTC),
						},
					},
				}, nil
			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
				assert.Nil(t, input.UnlockedAt)
				utilstest.AssertFieldsNoEmptyFromStruct(t, input)
				utilstest.AssertStructEqual(t, mockSecondBlockUserWithNewCapture(), input)
				assert.Len(t, input.CapturedRepairs, 2)
				assert.Equal(t, input.CountCaptures, 2)

				return nil
			},
			GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
				return &storage.BlockedUser{
					KycIdentificationID: "abc",
					UnlockedAt:          nil,
					SiteID:              "MLM",
					CreatedAt:           time.Now().UTC(),
					LastBlockedAt:       time.Now().UTC(),
					CountCaptures:       1,
					CapturedRepairs: []storage.CapturedRepairs{
						{
							AuthorizationID: "auth_01",
							UserID:          123,
							PaymentID:       "pay_01",
							TypeRepair:      domain.TypeReverse,
							RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
							CreatedAt:       time.Date(2022, 10, 2, 0, 0, 0, 0, time.UTC),
						},
					},
				}, nil
			},
		}

		service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

		output, err := service.BlockUserIfReparationExistsInBulk(
			ctx,
			domain.TransactionsNewsFeedBulkInput{
				SiteID:                 "MLM",
				AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_02": "cap02", "auth_03": "cap_03"}},
			field.NewWrappedFields(),
		)

		assert.NoError(t, err)
		assert.Len(t, output.AuthorizationIDsWithErr, 0)
	})
}

func Test_blocked_BlockUserIfReparationExistsInBulk_Block_And_Append_Error_Save(t *testing.T) {
	t.Run("reparation exists, it is the 2º capture with reparation done, but fail", func(t *testing.T) {
		ctx := context.TODO()

		mockReparationSearchRepo := mockReparationSearchRepo{
			GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
				return []storage.ReparationOutput{
					{
						BaseReparation: storage.BaseReparation{
							KycIdentificationID: "abc",
							AuthorizationID:     authorizationIDs[0],
							UserID:              123,
							SiteID:              "MLM",
							PaymentID:           "pay_02",
							TransactionRepairID: "rev_02",
							Type:                domain.TypeReverse,
							CreatedAt:           time.Date(2022, 10, 12, 0, 0, 0, 0, time.UTC),
						},
					},
				}, nil
			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
				return gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)
			},
			GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
				return &storage.BlockedUser{
					KycIdentificationID: "abc",
					UnlockedAt:          nil,
					SiteID:              "MLM",
					CreatedAt:           time.Now().UTC(),
					LastBlockedAt:       time.Now().UTC(),
					CountCaptures:       1,
					CapturedRepairs: []storage.CapturedRepairs{
						{
							AuthorizationID: "auth_01",
							UserID:          123,
							PaymentID:       "pay_01",
							TypeRepair:      domain.TypeReverse,
							RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
							CreatedAt:       time.Date(2022, 10, 2, 0, 0, 0, 0, time.UTC),
						},
					},
				}, nil
			},
		}

		service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

		output, err := service.BlockUserIfReparationExistsInBulk(
			ctx,
			domain.TransactionsNewsFeedBulkInput{
				SiteID:                 "MLM",
				AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_02": "cap02"}},
			field.NewWrappedFields(),
		)

		assert.NoError(t, err)

		errorExpected := gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)
		utilstest.AssertGnierrorsExpected(t, errorExpected, output.AuthorizationIDsWithErr["auth_02"])
	})
}

func Test_blocked_BlockUserIfReparationExistsInBulk_Error_GetBlockedlist_Unknown(t *testing.T) {
	ctx := context.TODO()

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByAuthorizationIDsStub: func(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
			return []storage.ReparationOutput{
				mockTransactionRepair01(),
			}, nil
		},
	}

	mockBlockedlistRepo := mockBlockedlistRepo{
		GetStub: func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist database failure", false)
		},
	}

	service := New(mockReparationSearchRepo, mockBlockedlistRepo, mockBlockedUserSearchRepo{})

	output, err := service.BlockUserIfReparationExistsInBulk(
		ctx,
		domain.TransactionsNewsFeedBulkInput{
			SiteID:                 "MLM",
			AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{"auth_01": "cap01", "auth_02": "cap02"}},
		field.NewWrappedFields(),
	)

	assert.NoError(t, err)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"), "blockedlist database failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, output.AuthorizationIDsWithErr["auth_01"])
}
