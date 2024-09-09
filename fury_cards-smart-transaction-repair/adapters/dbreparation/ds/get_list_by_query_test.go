package ds

import (
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_dsTransactionRepair_GetListByUserIDAndCreationPeriod_Success(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			b, err := queryBuilder.Build().MarshalJSON()
			assert.NoError(t, err)

			queryExpected := `{
				"and":[{"eq":{"field":"user_id","value":1}},
				{"date_range":{"field":"created_at","gte":"2022-10-27"}},
				{"date_range":{"field":"created_at","lte":"2022-11-26"}}]
			}`

			assert.JSONEq(t, queryExpected, string(b))

			return _transactionRepairDBList, nil
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	gte, lte := mockGenerateDates()
	output, err := dsClient.GetListByUserIDAndCreationPeriod(1, gte, lte)
	assertCommonGetList(t, output, err)
}

func Test_dsTransactionRepair_GetListByKycIentificationAndCreationPeriod_Success(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			b, err := queryBuilder.Build().MarshalJSON()
			assert.NoError(t, err)

			queryExpected := `{
				"and":[{"eq":{"field":"kyc_identification_id","value":"abc"}},
				{"date_range":{"field":"created_at","gte":"2022-10-27"}},
				{"date_range":{"field":"created_at","lte":"2022-11-26"}}]
			}`

			assert.JSONEq(t, queryExpected, string(b))

			return _transactionRepairDBList, nil
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	gte, lte := mockGenerateDates()
	output, err := dsClient.GetListByKycIentificationAndCreationPeriod("abc", gte, lte)
	assertCommonGetList(t, output, err)
}

func assertCommonGetList(t *testing.T, output []storage.ReparationOutput, err error) {
	assert.NoError(t, err)
	assert.Len(t, output, 2)
	assert.Equal(t, "auth_1", output[0].AuthorizationID)
	assert.Equal(t, "abc", output[0].KycIdentificationID)
	assert.Equal(t, int64(123), output[0].UserID)
	assert.Equal(t, "1", output[0].PaymentID)
	assert.Equal(t, "MLM", output[0].SiteID)
	assert.Equal(t, "reverse_smart_auth_1", output[0].TransactionRepairID)
	assert.Equal(t, domain.TypeReverse, output[0].Type)
	assert.Equal(t, time.Date(2022, 11, 10, 0, 0, 0, 0, time.UTC), output[0].CreatedAt)
}

func Test_dsTransactionRepair_GetListByUserIDAndCreationPeriod_Err_Unmarshal(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			return []byte("}"), nil
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	gte, lte := mockGenerateDates()
	output, err := dsClient.GetListByUserIDAndCreationPeriod(1, gte, lte)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("invalid character '}' looking for beginning of value"),
		"transaction repair document search failure", false)

	assert.Nil(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_dsTransactionRepair_GetListByUserIDAndCreationPeriod_Err_Client(t *testing.T) {
	mockDSClient := mockDSClient{
		GetListByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error) {
			return nil, errors.New("some unknown error")
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	gte, lte := mockGenerateDates()
	output, err := dsClient.GetListByUserIDAndCreationPeriod(1, gte, lte)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"),
		"transaction repair document search failure", false)

	assert.Nil(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
