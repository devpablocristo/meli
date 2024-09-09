package configurations

import (
	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

const (
	msgErrorConfigurations = "error in configurations"
	msgErrNotFoundProfile  = "no such file or directory"
)

func buildErrorNotFound(profile string, err error) error {
	return domain.BuildErrorNotFound(profile, err.Error(), msgErrorConfigurations)
}

func buildErrorBadGateway(err error) error {
	return domain.BuildErrorBadGateway(err, msgErrorConfigurations)
}
