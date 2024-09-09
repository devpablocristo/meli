package configurations

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/app-init_struct.json
var _applicationInit_struct []byte

//go:embed testdata/parameters-validation_struct.json
var _parametersValidation_struct []byte

func mockConfigApplicationInitExpected(t *testing.T) *domain.ConfigApplicationInit {
	var configApplicationInit domain.ConfigApplicationInit

	err := json.Unmarshal(_applicationInit_struct, &configApplicationInit)
	assert.NoError(t, err)

	return &configApplicationInit
}

func mockConfigParametersValidationExpected(t *testing.T) *domain.ConfigParametersValidation {
	var configParametersValidation domain.ConfigParametersValidation

	err := json.Unmarshal(_parametersValidation_struct, &configParametersValidation)
	assert.NoError(t, err)

	return &configParametersValidation
}
