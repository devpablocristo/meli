package validation

import (
	"context"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

func Test_validationService_Save_Success(t *testing.T) {
	inputResult := mockExpectedEligible("MLB")

	v := validationService{
		validationResultRepo: mockValidationResultRepo{
			SaveStub: func(ctx context.Context, paymentID string, validationResult storage.ValidationResult) error {
				assert.Equal(t, inputResult, validationResult)
				assert.Nil(t, validationResult.Reason)
				return nil
			},
		},
	}

	err := v.Save(context.TODO(), inputResult, field.NewWrappedFields())
	assert.NoError(t, err)
}

func Test_validationService_Save_Error(t *testing.T) {
	inputResult := mockExpectedEligible("MLB")

	v := validationService{
		validationResultRepo: mockValidationResultRepo{
			SaveStub: func(ctx context.Context, paymentID string, validationResult storage.ValidationResult) error {
				return domain.BuildErrorBadGateway(errors.New("any error"), "ds failed")
			},
		},
	}

	err := v.Save(context.TODO(), inputResult, field.NewWrappedFields())
	assert.Error(t, err)
}
