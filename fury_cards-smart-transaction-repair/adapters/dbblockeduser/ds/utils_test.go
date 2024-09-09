package ds

import (
	_ "embed"

	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

//go:embed testdata/blocked_user_list_db.json
var _blockedUserListDB []byte

// dsClientService
type mockDSClient struct {
	dsclient.Service
	GetListByQueryBuilderStub func(queryBuilder dsclient.QueryBuilder, opts ...dsclient.OptionParameters) ([]byte, error)
}

func (m mockDSClient) GetListByQueryBuilder(
	queryBuilder dsclient.QueryBuilder,
	opts ...dsclient.OptionParameters,
) ([]byte, error) {
	return m.GetListByQueryBuilderStub(queryBuilder)
}
