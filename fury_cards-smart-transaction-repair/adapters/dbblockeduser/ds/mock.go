package ds

import (
	_ "embed"
	"encoding/json"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

//go:embed testdata/blocked_user_list_domain.json
var _blockedUserListDomain []byte

type mock struct {
}

func (m mock) GetListByAuthorizationIDs(
	siteID string,
	authorizationIDs storage.AuthorizationIDsSearch,
) (list []storage.BlockedUser, err error) {

	err = json.Unmarshal(_blockedUserListDomain, &list)
	return list, err
}

func NewMock() ports.BlockedUserSearchRepository {
	return &mock{}
}
