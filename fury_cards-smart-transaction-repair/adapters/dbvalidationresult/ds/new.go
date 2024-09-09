package ds

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	dsclient "github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
)

type dsValidationResult struct {
	repo dsclient.Service
}

func New(dsClient dsclient.Service) ports.ValidationResultRepository {
	return &dsValidationResult{
		repo: dsClient,
	}
}
