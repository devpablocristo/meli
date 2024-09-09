package blockeduser

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

// BlockedUserRepository.
type mockBlockedlistRepo struct {
	GetStub  func(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error)
	SaveStub func(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error
}

func (m mockBlockedlistRepo) Get(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error) {
	return m.GetStub(ctx, kycIdentificationID)
}

func (m mockBlockedlistRepo) Save(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error {
	return m.SaveStub(ctx, kycIdentificationID, input)
}

// BlockedUserSearchRepository.
type mockBlockedUserSearchRepo struct {
	GetListByAuthorizationIDsStub func(siteID string, authorizationIDs storage.AuthorizationIDsSearch) ([]storage.BlockedUser, error)
}

func (m mockBlockedUserSearchRepo) GetListByAuthorizationIDs(
	siteID string,
	authorizationIDs storage.AuthorizationIDsSearch,
) ([]storage.BlockedUser, error) {

	return m.GetListByAuthorizationIDsStub(siteID, authorizationIDs)
}

// ReparationSearchRepository.
type mockReparationSearchRepo struct {
	ports.ReparationSearchRepository
	GetListByAuthorizationIDsStub func(authorizationIDs ...string) ([]storage.ReparationOutput, error)
}

func (m mockReparationSearchRepo) GetListByAuthorizationIDs(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
	return m.GetListByAuthorizationIDsStub(authorizationIDs...)
}

// Struct Mocks:

func mockListFirstBlockedUser() []storage.BlockedUser {
	return []storage.BlockedUser{
		*mockBlockeUser01(),
		*mockBlockeUser02(),
	}
}

func mockBlockedUsersByAuthorizationIDs(authorizationIDs storage.AuthorizationIDsSearch) []storage.BlockedUser {
	blockedUsers := []storage.BlockedUser{}

	for _, id := range authorizationIDs {
		blockedUser := mockBlockedUserByAuthorizataionID(id.(string))
		if blockedUser != nil {
			blockedUsers = append(blockedUsers, *blockedUser)
		}
	}
	return blockedUsers
}

func mockBlockedUserByAuthorizataionID(authorizationID string) *storage.BlockedUser {
	switch authorizationID {
	case "auth_01":
		return mockBlockeUser01()
	case "auth_02":
		return mockBlockeUser02()
	case "auth_03", "auth_03_1":
		return mockBlockeUser03()
	case "auth_04", "auth_04_1":
		return mockBlockeUser04()
	default:
		return nil
	}
}

func mockBlockeUser01() *storage.BlockedUser {
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
				PaymentID:       "pay_01",
				UserID:          123,
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

func mockBlockeUser02() *storage.BlockedUser {
	return &storage.BlockedUser{
		KycIdentificationID: "def",
		CountCaptures:       1,
		UnlockedAt:          nil,
		SiteID:              "MLM",
		CreatedAt:           time.Now().UTC(),
		LastBlockedAt:       time.Now().UTC(),
		CapturedRepairs: []storage.CapturedRepairs{
			{
				AuthorizationID: "auth_02",
				UserID:          456,
				PaymentID:       "pay_02",
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

func mockBlockeUser03() *storage.BlockedUser {
	return &storage.BlockedUser{
		KycIdentificationID: "ghi",
		CountCaptures:       2,
		UnlockedAt:          nil,
		SiteID:              "MLM",
		CreatedAt:           time.Now().UTC(),
		LastBlockedAt:       time.Now().UTC(),
		CapturedRepairs: []storage.CapturedRepairs{
			{
				AuthorizationID: "auth_03",
				UserID:          41,
				PaymentID:       "pay_03",
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now().UTC(),
			},
			{
				AuthorizationID: "auth_03_1",
				UserID:          42,
				PaymentID:       "pay_03_1",
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

func mockBlockeUser04() *storage.BlockedUser {
	return &storage.BlockedUser{
		KycIdentificationID: "jkl",
		UnlockedAt:          nil,
		SiteID:              "MLM",
		CreatedAt:           time.Now().UTC(),
		LastBlockedAt:       time.Now().UTC(),
		CountCaptures:       1,
		CapturedRepairs: []storage.CapturedRepairs{
			{
				AuthorizationID: "auth_04",
				PaymentID:       "pay_04",
				UserID:          10,
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now().UTC(),
				ReversalClearing: &storage.ReversalClearing{
					ReverseID: "rev_04",
					CreatedAt: time.Date(2022, 10, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				AuthorizationID: "auth_04_1",
				PaymentID:       "pay_04_1",
				UserID:          11,
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

func mockSecondBlockUserWithNewCapture() storage.BlockedUser {
	return storage.BlockedUser{
		KycIdentificationID: "abc",
		UnlockedAt:          nil,
		SiteID:              "MLM",
		CreatedAt:           time.Now().UTC(),
		LastBlockedAt:       time.Now().UTC(),
		CountCaptures:       2,
		CapturedRepairs: []storage.CapturedRepairs{
			{
				AuthorizationID: "auth_01",
				UserID:          123,
				PaymentID:       "pay_01",
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Date(2022, 10, 2, 0, 0, 0, 0, time.UTC),
			},
			{
				AuthorizationID: "auth_02",
				UserID:          123,
				PaymentID:       "pay_02",
				TypeRepair:      domain.TypeReverse,
				RepairedAt:      time.Date(2022, 10, 12, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

func mockTransactionRepair01() storage.ReparationOutput {
	return storage.ReparationOutput{
		BaseReparation: storage.BaseReparation{
			KycIdentificationID: "abc",
			AuthorizationID:     "auth_01",
			UserID:              123,
			SiteID:              "MLM",
			PaymentID:           "pay_01",
			TransactionRepairID: "reverse_smart_01",
			Type:                domain.TypeReverse,
			CreatedAt:           time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
			Requester:           domain.Requester{ClientApplication: "app_test"},
		},
	}
}

func mockTransactionRepair02() storage.ReparationOutput {
	return storage.ReparationOutput{
		BaseReparation: storage.BaseReparation{
			KycIdentificationID: "def",
			AuthorizationID:     "auth_02",
			UserID:              456,
			SiteID:              "MLM",
			PaymentID:           "pay_02",
			TransactionRepairID: "reverse_smart_02",
			Type:                domain.TypeReverse,
			CreatedAt:           time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC),
			FaqID:               "faq_test",
		},
	}
}
