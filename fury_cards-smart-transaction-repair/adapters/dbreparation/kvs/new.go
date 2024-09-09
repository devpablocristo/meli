package kvs

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	kvs "github.com/melisource/fury_cards-go-toolkit/pkg/kvsclient/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

type kvsRepository struct {
	repo kvs.Service
	log  log.LogService
}

func New(ksvClient kvs.Service, log log.LogService) ports.ReparationRepository {
	return &kvsRepository{
		repo: ksvClient,
		log:  log,
	}
}
