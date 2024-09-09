package reparation

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

type reparation struct {
	log                  log.LogService
	validation           ports.ValidationService
	cardsTransactions    ports.CardsTransactions
	reparationRepo       ports.ReparationRepository
	searchHub            ports.SearchHub
	reparationSearchRepo ports.ReparationSearchRepository
	searchKyc            ports.SearchKyc
}

var _ ports.ReparationService = (*reparation)(nil)

func New(
	log log.LogService,
	validation ports.ValidationService,
	reparationRepo ports.ReparationRepository,
	searchHub ports.SearchHub,
	cardsTransactions ports.CardsTransactions,
	reparationSearchRepo ports.ReparationSearchRepository,
	searchKyc ports.SearchKyc,
) ports.ReparationService {

	return &reparation{
		reparationRepo:       reparationRepo,
		searchHub:            searchHub,
		cardsTransactions:    cardsTransactions,
		validation:           validation,
		reparationSearchRepo: reparationSearchRepo,
		log:                  log,
		searchKyc:            searchKyc,
	}
}
