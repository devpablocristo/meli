package validation

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

// Common.
const (
	oneDay = time.Hour * 24
)

type validationService struct {
	log                   log.LogService
	configurations        ports.Configurations
	validationResultRepo  ports.ValidationResultRepository
	reparationSearchRepo  ports.ReparationSearchRepository
	blockedUserService    ports.BlockedUserService
	policy                ports.Policy
	searchHub             ports.SearchHub
	eventValidationResult ports.EventValidationResult
	scoreService          ports.Score
}

var _ ports.ValidationService = (*validationService)(nil)

func New(
	log log.LogService,
	configurations ports.Configurations,
	validationResultRepo ports.ValidationResultRepository,
	reparationSearchRepo ports.ReparationSearchRepository,
	blockedUserService ports.BlockedUserService,
	policy ports.Policy,
	searchHub ports.SearchHub,
	eventValidationResult ports.EventValidationResult,
	scoreService ports.Score,
) ports.ValidationService {
	return &validationService{
		log:                   log,
		configurations:        configurations,
		validationResultRepo:  validationResultRepo,
		reparationSearchRepo:  reparationSearchRepo,
		blockedUserService:    blockedUserService,
		policy:                policy,
		searchHub:             searchHub,
		eventValidationResult: eventValidationResult,
		scoreService:          scoreService,
	}
}
