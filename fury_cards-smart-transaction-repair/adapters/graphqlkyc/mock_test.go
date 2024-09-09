package graphqlkyc

import (
	"context"
	_ "embed"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/stretchr/testify/assert"
)

func Test_mock_GetUserKyc(t *testing.T) {
	m := NewMock()
	output, err := m.GetUserKyc(context.TODO(), domain.SearchKycInput{UserID: "12"})

	assert.NoError(t, err)
	assert.Equal(t, "abc", output.KycIdentificationID)
}
