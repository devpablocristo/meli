package ports

import (
	"context"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
)

type ReparationRepository interface {
	Save(ctx context.Context, transactionRepair storage.ReparationInput) error
	Get(ctx context.Context, authorizationID string) (*storage.ReparationOutput, error)
}

type ReparationSearchRepository interface {
	GetListByUserIDAndCreationPeriod(
		userID int64,
		gteCreatedAt,
		lteCreatedAt time.Time) ([]storage.ReparationOutput, error)

	GetListByKycIentificationAndCreationPeriod(
		kycIdentificationID string,
		gteCreatedAt,
		lteCreatedAt time.Time) ([]storage.ReparationOutput, error)

	GetByPaymentIDAndUserID(paymentID string, userID int64) (*storage.ReparationOutput, error)

	GetListByAuthorizationIDs(authorizationIDs ...string) ([]storage.ReparationOutput, error)
}

type BlockedUserRepository interface {
	Get(ctx context.Context, kycIdentificationID string) (*storage.BlockedUser, error)
	Save(ctx context.Context, kycIdentificationID string, input storage.BlockedUser) error
}

type BlockedUserSearchRepository interface {
	GetListByAuthorizationIDs(siteID string, authorizationIDs storage.AuthorizationIDsSearch) ([]storage.BlockedUser, error)
}

type ValidationResultRepository interface {
	Save(ctx context.Context, paymentID string, validationResult storage.ValidationResult) error
}

type SearchTransaction interface {
	GetTransactionByPaymentID(
		ctx context.Context,
		transactionInput *domain.TransactionInput,
		wrapFields *field.WrappedFields) (*domain.TransactionOutput, error)
}

type SearchHub interface {
	CountChargebackFromUser(
		ctx context.Context,
		input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error)

	GetTransactionAndCards(ctx context.Context, transactionInput domain.TransactionWalletInput) (*domain.TransactionWalletOutput, error)

	GetTransactionHistory(
		ctx context.Context,
		transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error)

	GetPaymentTransaction(
		ctx context.Context,
		transactionInput domain.TransactionInput,
		wrapFields *field.WrappedFields) (domain.TransactionOutput, error)
}

type SearchKyc interface {
	GetUserKyc(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error)
}

type Configurations interface {
	LoadApplicationInit(ctx context.Context) (*domain.ConfigApplicationInit, error)
	LoadParametersValidation(ctx context.Context) (*domain.ConfigParametersValidation, error)
	ResetTTLDefaultCache(ttl time.Duration)
}

type CardsTransactions interface {
	Reverse(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error)
}

type Policy interface {
	EvaluateWithUser(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error)
}

type EventValidationResult interface {
	Publish(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error
}

type Score interface {
	GetScore(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error)
}
