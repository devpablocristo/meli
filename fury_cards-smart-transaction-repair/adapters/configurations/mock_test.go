package configurations

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mock_LoadApplicationInit_Success(t *testing.T) {
	mock := NewMock()
	appInit, err := mock.LoadApplicationInit(context.TODO())

	assert.NoError(t, err)
	assert.NotEmpty(t, appInit)
}

func Test_mock_LoadParametersValidation_Success(t *testing.T) {
	mock := NewMock()
	parametersValidation, err := mock.LoadParametersValidation(context.TODO())

	assert.NoError(t, err)
	assert.NotEmpty(t, parametersValidation)
}
