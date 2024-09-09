package ds

import (
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"

	"github.com/stretchr/testify/assert"
)

func Test_dsBlockedUser_GetListByAuthorizationIDs(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			b, err := queryBuilder.Build().MarshalJSON()
			assert.NoError(t, err)

			queryExpected :=
				`{"and":[{"eq":{"field":"site_id","value":"MLM"}},
				{"in":{"field":"captured_repairs.authorization_id","value":["auth_1","auth_2"]}}]}`

			assert.JSONEq(t, queryExpected, string(b))

			return _blockedUserListDB, nil
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	blockedUsers, err := dsClient.GetListByAuthorizationIDs("MLM", storage.AuthorizationIDsSearch{"auth_1", "auth_2"})

	assert.NoError(t, err)
	assert.Len(t, blockedUsers, 2)
	assert.Equal(t, "abc", blockedUsers[0].KycIdentificationID)
	assert.Equal(t, "def", blockedUsers[1].KycIdentificationID)
	utilstest.AssertFieldsNoEmptyFromStruct(t, blockedUsers[0])
	utilstest.AssertFieldsNoEmptyFromStruct(t, blockedUsers[1])
}

func Test_dsBlockedUser_GetListByAuthorizationIDs_Error_BadGateway(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			return nil, errors.New("error unknown")
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	repairs, err := dsClient.GetListByAuthorizationIDs("auth_1", storage.AuthorizationIDsSearch{"auth_2"})

	assert.Error(t, err)
	assert.Len(t, repairs, 0)
}

func Test_dsBlockedUser_GetListByAuthorizationIDs_Error_Unmarshal(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			return []byte(`{`), nil
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	repairs, err := dsClient.GetListByAuthorizationIDs("auth_1", storage.AuthorizationIDsSearch{"auth_2"})

	assert.Error(t, err)
	assert.Len(t, repairs, 0)
}
