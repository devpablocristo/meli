package ports

import (
	"context"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

type ReparationService interface {
	Reverse(ctx context.Context, input domain.ReverseTransactionInput, wrapFields *field.WrappedFields) error
}

type ValidationService interface {
	ExecuteValidation(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error
	Save(ctx context.Context, inputResult storage.ValidationResult, wrapFields *field.WrappedFields) error
	PublishEvent(ctx context.Context, validationResult *storage.ValidationResult, wrapFields *field.WrappedFields)
}

type BlockedUserService interface {
	BlockUserIfReparationExistsInBulk(
		ctx context.Context,
		input domain.TransactionsNewsFeedBulkInput,
		wrapFields *field.WrappedFields) (*domain.TransactionsBulkOutput, error)
	CheckUserIsBlocked(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error)
	UnblockUser(ctx context.Context, kycIdentificationID string) error
	UnblockUserByReversalClearingInBulk(
		ctx context.Context,
		input domain.TransactionsNewsFeedBulkInput,
		wrapFields *field.WrappedFields) (*domain.TransactionsBulkOutput, error)
}
