package ds

import (
	_ "embed"
	"time"

	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

//go:embed testdata/transaction-repair-list-db.json
var _transactionRepairDBList []byte

// dsClientService
type mockDSClient struct {
	dsclient.Service
	GetListByQueryBuilderStub     func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error)
	GetDocumentByQueryBuilderStub func(queryBuilder dsclient.QueryBuilder, v interface{}) error
}

func (m mockDSClient) GetListByQueryBuilder(
	queryBuilder dsclient.QueryBuilder,
	opts ...dsclient.OptionParameters,
) ([]byte, error) {
	return m.GetListByQueryBuilderStub(queryBuilder)
}

func (m mockDSClient) GetDocumentByQueryBuilder(queryBuilder dsclient.QueryBuilder, v interface{}) error {
	return m.GetDocumentByQueryBuilderStub(queryBuilder, v)
}

func mockGenerateDates() (time.Time, time.Time) {
	dateActual := time.Date(2022, 11, 26, 0, 0, 0, 0, time.UTC)
	qtyDays := 30
	oneDay := time.Hour * 24
	gte := dateActual
	return gte.Add(-oneDay * time.Duration(qtyDays)), dateActual
}
