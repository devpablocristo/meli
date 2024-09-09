package ds

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_dsTransactionRepair_GetByPaymentIDAndUserID_Success(t *testing.T) {
	mockDSClient := mockDSClient{
		GetDocumentByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, v interface{}) error {
			b, err := queryBuilder.Build().MarshalJSON()
			assert.NoError(t, err)

			queryExpected := `{"and":[{"eq":{"field":"payment_id","value":"auth_1"}},{"and":[{"eq":{"field":"user_id","value":123}}]}]}`

			assert.JSONEq(t, queryExpected, string(b))

			return json.Unmarshal(_transactionRepairDB, v)
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	transactionRepair, err := dsClient.GetByPaymentIDAndUserID("auth_1", 123)

	assert.NoError(t, err)
	assert.NotEmpty(t, transactionRepair)
	assert.Equal(t, "auth_1", transactionRepair.AuthorizationID)
	assert.Equal(t, "abc", transactionRepair.KycIdentificationID)
	assert.Equal(t, int64(123), transactionRepair.UserID)
	assert.Equal(t, "1", transactionRepair.PaymentID)
	assert.Equal(t, "MLM", transactionRepair.SiteID)
	assert.Equal(t, "reverse_smart_auth_1", transactionRepair.TransactionRepairID)
	assert.Equal(t, domain.TypeReverse, transactionRepair.Type)
	assert.Equal(t, time.Date(2022, 11, 10, 0, 0, 0, 0, time.UTC), transactionRepair.CreatedAt)
}

func Test_dsTransactionRepair_GetByPaymentIDAndUserID_Error_NotFound(t *testing.T) {
	mockDSClient := mockDSClient{
		GetDocumentByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, v interface{}) error {
			return dsclient.ErrDocumentNotFound
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	transactionRepair, err := dsClient.GetByPaymentIDAndUserID("auth_1", 123)

	errorExpected := gnierrors.Wrap(domain.NotFound,
		errors.New("document not found | ID: auth_1"),
		"transaction repair document search failure", false)

	assert.Nil(t, transactionRepair)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_dsTransactionRepair_GetByPaymentIDAndUserID_Error_Unknown(t *testing.T) {
	mockDSClient := mockDSClient{
		GetDocumentByQueryBuilderStub: func(queryBuilder dsclient.QueryBuilder, v interface{}) error {
			return errors.New("some unknown error")
		},
	}

	dsQuery := dsclient.NewQuery()
	dsClient := New(mockDSClient, dsQuery)

	transactionRepair, err := dsClient.GetByPaymentIDAndUserID("auth_1", 123)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("some unknown error"),
		"transaction repair document search failure", false)

	assert.Nil(t, transactionRepair)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}
