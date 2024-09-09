package reparation

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/transaction-data.json
var _transactionData []byte

// LogService.
type mockLogService struct {
	log.LogService
	ErrorStub  func(c context.Context, msg string, fields ...log.Field)
	ErrStub    func(err error) log.Field
	StringStub func(key string, val string) log.Field
	WarnStub   func(c context.Context, msg string, fields ...log.Field)
	InfoStub   func(c context.Context, msg string, fields ...log.Field)
	AnyStub    func(key string, value interface{}) log.Field
}

func (m mockLogService) Error(c context.Context, msg string, fields ...log.Field) {
	m.ErrorStub(c, msg, fields...)
}

func (m mockLogService) Err(err error) log.Field {
	return m.ErrStub(err)
}

func (m mockLogService) String(key string, val string) log.Field {
	return m.StringStub(key, val)
}

func (m mockLogService) Warn(c context.Context, msg string, fields ...log.Field) {
	m.WarnStub(c, msg, fields...)
}

func (m mockLogService) Info(c context.Context, msg string, fields ...log.Field) {
	m.InfoStub(c, msg, fields...)
}

func (m mockLogService) Any(key string, value interface{}) log.Field {
	return m.AnyStub(key, value)
}

// ReparationRepository.
type mockTransactionRepairRepo struct {
	SaveStub func(ctx context.Context, transactionRepair storage.ReparationInput) error
	GetStub  func(ctx context.Context, authorizationID string) (*storage.ReparationOutput, error)
}

func (m mockTransactionRepairRepo) Save(ctx context.Context, transactionRepair storage.ReparationInput) error {
	return m.SaveStub(ctx, transactionRepair)
}

func (m mockTransactionRepairRepo) Get(ctx context.Context, paymentID string) (*storage.ReparationOutput, error) {
	return m.GetStub(ctx, paymentID)
}

// SearchHub.
type mockSearchHub struct {
	ports.SearchHub
	GetPaymentTransactionStub func(
		ctx context.Context,
		transactionInput domain.TransactionInput,
		wrapFields *field.WrappedFields,
	) (domain.TransactionOutput, error)
}

func (m mockSearchHub) GetPaymentTransaction(
	ctx context.Context, transactionInput domain.TransactionInput, wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
	return m.GetPaymentTransactionStub(ctx, transactionInput, wrapFields)
}

// CardsTransactions.
type mockCardsTransactions struct {
	ReverseStub func(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error)
}

func (m mockCardsTransactions) Reverse(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
	return m.ReverseStub(ctx, input)
}

// ValidationService
type mockValidationService struct {
	ports.ValidationService
	ExecuteValidationStub func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error
}

func (m mockValidationService) ExecuteValidation(
	ctx context.Context,
	input domain.ValidationInput,
	wrapFields *field.WrappedFields,
) error {
	return m.ExecuteValidationStub(ctx, input, wrapFields)
}

// ReparationSearchRepository.
type mockReparationSearchRepo struct {
	ports.ReparationSearchRepository
	GetByPaymentIDAndUserIDStub func(paymentID string, userID int64) (*storage.ReparationOutput, error)
}

func (m mockReparationSearchRepo) GetByPaymentIDAndUserID(paymentID string, userID int64) (*storage.ReparationOutput, error) {
	return m.GetByPaymentIDAndUserIDStub(paymentID, userID)
}

// SearchKyc.
type mockSearchKyc struct {
	GetUserKycStub func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error)
}

func (m mockSearchKyc) GetUserKyc(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
	return m.GetUserKycStub(ctx, input)
}

// Data mocks:

func mockTransactionDataDefault(t *testing.T) (transactionData *domain.TransactionData) {
	err := json.Unmarshal(_transactionData, &transactionData)
	assert.NoError(t, err)
	return
}

func mockBaseTransactionRepairExpected() storage.BaseReparation {
	return storage.BaseReparation{
		AuthorizationID:     "auth_mlb_test_aotorres_FyLRqgiKVMw9Vb8y737ppvNBsYpvkf",
		PaymentID:           "100",
		KycIdentificationID: "abc",
		UserID:              123,
		TransactionRepairID: "reverse_auth_mlb_test_aotorres_FyLRqgiKVMw9Vb8y737ppvNBsYpvkf",
		SiteID:              "MLB",
		Type:                domain.TypeReverse,
		CreatedAt:           time.Now().UTC(),
		FaqID:               "faq_01",
		Requester: domain.Requester{
			ClientApplication: "card-admin",
			ClientScope:       "beta",
		},
		MonetaryTransaction: storage.MonetaryTransaction{
			Billing:    storage.MonetaryValue{Amount: 350.00, Currency: "BRL"},
			Settlement: storage.MonetaryValue{Amount: 18.99, Currency: "USD"},
		},
	}
}

func mockValidationInputExpected(t *testing.T) domain.ValidationInput {
	transactionData := *mockTransactionDataDefault(t)
	return domain.ValidationInput{
		PaymentID:       "100",
		UserID:          123,
		Type:            domain.TypeReverse,
		SiteID:          transactionData.SiteID,
		FaqID:           "faq_01",
		TransactionData: transactionData,
		UserKycData: domain.SearchKycOutput{
			KycIdentificationID: "abc",
			DateCreated:         time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC),
		},
		Requester: domain.Requester{
			ClientApplication: "card-admin",
			ClientScope:       "beta",
		},
	}
}
