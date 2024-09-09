package ds

import (
	"encoding/json"

	"github.com/melisource/cards-smart-transaction-repair/adapters/dbblockeduser"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

// Max default value for the fury consumer bulk delivery mode.
const defaultSize = 20

func (d *dsBlockedUser) GetListByAuthorizationIDs(
	siteID string,
	authorizationIDs storage.AuthorizationIDsSearch,
) ([]storage.BlockedUser, error) {

	queryBuilder := d.query.And(d.query.Eq("site_id", siteID)).And(d.query.InMultiple("captured_repairs.authorization_id", authorizationIDs))

	documents, err := d.repo.GetListByQueryBuilder(queryBuilder, dsclient.WithSize(defaultSize))
	if err != nil {
		return nil, buildErrorBadGateway(err)
	}

	listOutput := []storage.BlockedUser{}

	var blockedUsersDB []dbblockeduser.BlockedUserDB
	err = json.Unmarshal(documents, &blockedUsersDB)
	if err != nil {
		return nil, buildErrorBadGateway(err)
	}

	for _, blockedUserDB := range blockedUsersDB {
		listOutput = append(listOutput, *dbblockeduser.ToBlockedUserOutput(blockedUserDB))
	}

	return listOutput, nil
}
