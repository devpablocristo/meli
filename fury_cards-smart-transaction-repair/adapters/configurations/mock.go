package configurations

import (
	"context"
	_ "embed"
	"encoding/json"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
)

//go:embed testdata/app-init_struct.json
var _applicationInitMock []byte

//go:embed testdata/parameters-validation_struct.json
var _parametersValidationMock []byte

type mock struct{}

func (m mock) LoadApplicationInit(ctx context.Context) (*domain.ConfigApplicationInit, error) {
	var configApplicationInit domain.ConfigApplicationInit

	err := json.Unmarshal(_applicationInitMock, &configApplicationInit)

	return &configApplicationInit, err
}

func (m mock) LoadParametersValidation(ctx context.Context) (*domain.ConfigParametersValidation, error) {
	var configParametersValidation domain.ConfigParametersValidation

	err := json.Unmarshal(_parametersValidationMock, &configParametersValidation)

	return &configParametersValidation, err
}

func (m mock) ResetTTLDefaultCache(ttl time.Duration) {}

var _ ports.Configurations = (*configService)(nil)

func NewMock() ports.Configurations {
	return mock{}
}
