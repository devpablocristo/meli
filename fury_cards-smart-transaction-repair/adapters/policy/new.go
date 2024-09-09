package policy

import (
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/policyclient/v1"
)

type evaluate struct {
	pClient policyclient.PolicyClient
}

func New(policyClient policyclient.PolicyClient) ports.Policy {
	return &evaluate{pClient: policyClient}
}
