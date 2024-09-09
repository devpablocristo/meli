package ds

import (
	"errors"
	"testing"

	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/stretchr/testify/assert"
)

func Test_dsTransactionRepair_GetListByAuthorizationIDs_Success(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			b, err := queryBuilder.Build().MarshalJSON()
			assert.NoError(t, err)

			queryExpected := `{"ids":["auth_1","auth_2"]}`

			assert.JSONEq(t, queryExpected, string(b))

			return _transactionRepairDBList, nil
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	repairs, err := dsClient.GetListByAuthorizationIDs("auth_1", "auth_2")

	assert.NoError(t, err)
	assert.Len(t, repairs, 2)
	assert.Equal(t, "reverse_smart_auth_1", repairs[0].TransactionRepairID)
	assert.Equal(t, "reverse_smart_auth_2", repairs[1].TransactionRepairID)
	utilstest.AssertFieldsNoEmptyFromStruct(t, repairs[0])
	utilstest.AssertFieldsNoEmptyFromStruct(t, repairs[1])
}

func Test_dsTransactionRepair_GetListByAuthorizationIDs_Error_BadGateway(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			return nil, errors.New("error unknown")
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	repairs, err := dsClient.GetListByAuthorizationIDs("auth_1", "auth_2")

	assert.Error(t, err)
	assert.Len(t, repairs, 0)
}

func Test_dsTransactionRepair_GetListByAuthorizationIDs_Error_Unmarshal(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			return []byte(`{`), nil
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	repairs, err := dsClient.GetListByAuthorizationIDs("auth_1", "auth_2")

	assert.Error(t, err)
	assert.Len(t, repairs, 0)
}
