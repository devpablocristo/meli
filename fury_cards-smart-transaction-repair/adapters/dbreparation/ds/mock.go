package ds

import (
	_ "embed"
	"encoding/json"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
)

//go:embed testdata/transaction-repair-list-domain.json
var _transactionRepairDomainList []byte

//go:embed testdata/transaction-repair-db.json
var _transactionRepairDB []byte

type mock struct{}

func (m mock) GetListByUserIDAndCreationPeriod(
	userID int64,
	gteCreatedAt,
	lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {

	var list []storage.ReparationOutput

	return list, json.Unmarshal(_transactionRepairDomainList, &list)
}

func (m mock) GetListByKycIentificationAndCreationPeriod(
	kycIdentificationID string,
	gteCreatedAt,
	lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {

	var list []storage.ReparationOutput

	return list, json.Unmarshal(_transactionRepairDomainList, &list)
}

func (m mock) GetByPaymentIDAndUserID(paymentID string, userID int64) (*storage.ReparationOutput, error) {
	list, err := m.GetListByUserIDAndCreationPeriod(userID, time.Now(), time.Now())
	return &list[0], err
}

func (m mock) GetListByAuthorizationIDs(authorizationIDs ...string) ([]storage.ReparationOutput, error) {
	var list []storage.ReparationOutput

	return list, json.Unmarshal(_transactionRepairDomainList, &list)
}

func NewMock() ports.ReparationSearchRepository {
	return &mock{}
}
