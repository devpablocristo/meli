package blockeduser

import (
	"context"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_blocked_UnblockUserByReversalClearingInBulk_UnblockAll(t *testing.T) {
	t.Run("the user(s) is found, and all his authorizations have associated reverses: should unblock", func(t *testing.T) {
		ctx := context.TODO()

		inputBulk := domain.TransactionsNewsFeedBulkInput{
			SiteID: "MLM",
			AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{
				"auth_02":   "rev02",
				"auth_03":   "rev03",
				"auth_03_1": "rev03_1",
				"auth_04_1": "rev04_1",
				"auth_x":    "revx",
			},
		}

		mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{
			GetListByAuthorizationIDsStub: func(siteID string, authorizationIDs storage.AuthorizationIDsSearch) ([]storage.BlockedUser, error) {
				return mockBlockedUsersByAuthorizationIDs(authorizationIDs), nil
			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
				assert.NotNil(t, input.UnlockedAt)
				utilstest.AssertFieldsNoEmptyFromStruct(t, input)
				for _, cap := range input.CapturedRepairs {
					notNil := assert.NotNil(t, cap.ReversalClearing)
					if notNil {
						utilstest.AssertFieldsNoEmptyFromStruct(t, *cap.ReversalClearing)
					}
				}
				return nil
			},
		}

		service := New(mockReparationSearchRepo{}, mockBlockedlistRepo, mockBlockedUserSearchRepo)

		output, err := service.UnblockUserByReversalClearingInBulk(ctx, inputBulk, field.NewWrappedFields())

		assert.NoError(t, err)
		assert.Len(t, output.AuthorizationIDsWithErr, 0)
	})
}

func Test_blocked_UnblockUserByReversalClearingInBulk_NotUnblock(t *testing.T) {
	t.Run("the user(s) is found, but only with an authorization associated with the reverse: should not unblock", func(t *testing.T) {
		ctx := context.TODO()

		inputBulk := domain.TransactionsNewsFeedBulkInput{
			SiteID: "MLM",
			AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{
				"auth_03": "rev03",
				"auth_04": "rev04",
			},
		}

		mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{
			GetListByAuthorizationIDsStub: func(siteID string, authorizationIDs storage.AuthorizationIDsSearch) ([]storage.BlockedUser, error) {
				return mockBlockedUsersByAuthorizationIDs(authorizationIDs), nil
			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
				assert.Nil(t, input.UnlockedAt)
				utilstest.AssertFieldsNoEmptyFromStruct(t, input)
				for _, cap := range input.CapturedRepairs {
					if _, found := inputBulk.AuthorizationsNewsFeed[cap.AuthorizationID]; found {
						notNil := assert.NotNil(t, cap.ReversalClearing)
						if notNil {
							utilstest.AssertFieldsNoEmptyFromStruct(t, *cap.ReversalClearing)
						}
					} else {
						assert.Nil(t, cap.ReversalClearing)
					}
				}
				return nil
			},
		}

		service := New(mockReparationSearchRepo{}, mockBlockedlistRepo, mockBlockedUserSearchRepo)

		output, err := service.UnblockUserByReversalClearingInBulk(ctx, inputBulk, field.NewWrappedFields())

		assert.NoError(t, err)
		assert.Len(t, output.AuthorizationIDsWithErr, 0)
	})
}

func Test_blocked_UnblockUserByReversalClearingInBulk_Error_Save(t *testing.T) {
	t.Run("the user(s) is found: should unblock user(s), but fail", func(t *testing.T) {
		ctx := context.TODO()

		inputBulk := domain.TransactionsNewsFeedBulkInput{
			SiteID: "MLM",
			AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{
				"auth_01": "rev01",
			},
		}

		mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{
			GetListByAuthorizationIDsStub: func(siteID string, authorizationIDs storage.AuthorizationIDsSearch) ([]storage.BlockedUser, error) {
				return mockBlockedUsersByAuthorizationIDs(authorizationIDs), nil
			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{
			SaveStub: func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
				return gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)
			},
		}

		service := New(mockReparationSearchRepo{}, mockBlockedlistRepo, mockBlockedUserSearchRepo)

		output, err := service.UnblockUserByReversalClearingInBulk(ctx, inputBulk, field.NewWrappedFields())

		errExpected := gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)

		assert.NoError(t, err)
		utilstest.AssertGnierrorsExpected(t, errExpected, output.AuthorizationIDsWithErr["auth_01"])
	})
}

func Test_blocked_UnblockUserByReversalClearingInBulk_Error_GetBulk(t *testing.T) {
	t.Run("error in get, should return an error", func(t *testing.T) {
		ctx := context.TODO()

		inputBulk := domain.TransactionsNewsFeedBulkInput{
			SiteID: "MLM",
			AuthorizationsNewsFeed: map[string]domain.TransactionIDNewsFeed{
				"auth_01": "rev01",
			},
		}

		mockBlockedUserSearchRepo := mockBlockedUserSearchRepo{
			GetListByAuthorizationIDsStub: func(siteID string, authorizationIDs storage.AuthorizationIDsSearch) ([]storage.BlockedUser, error) {
				return nil, gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)
			},
		}

		mockBlockedlistRepo := mockBlockedlistRepo{}

		service := New(mockReparationSearchRepo{}, mockBlockedlistRepo, mockBlockedUserSearchRepo)

		output, err := service.UnblockUserByReversalClearingInBulk(ctx, inputBulk, field.NewWrappedFields())

		errExpected := gnierrors.Wrap(domain.BadGateway, errors.New("some unknown error"), "blockedlist repo failed", false)

		utilstest.AssertGnierrorsExpected(t, errExpected, err)
		assert.Nil(t, output)
	})
}
