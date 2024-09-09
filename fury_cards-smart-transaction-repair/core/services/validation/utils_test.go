package validation

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

//go:embed testdata/parameters-validation_struct.json
var _parametersValidationStruct []byte

// LogService
type mockLogService struct {
	log.LogService
	ErrorStub  func(c context.Context, msg string, fields ...log.Field)
	ErrStub    func(err error) log.Field
	InfoStub   func(c context.Context, msg string, fields ...log.Field)
	WarnStub   func(c context.Context, msg string, fields ...log.Field)
	AnyStub    func(key string, value interface{}) log.Field
	WarnlnStub func(ctx context.Context, msg string, fields ...log.Field)
}

func (m mockLogService) Error(c context.Context, msg string, fields ...log.Field) {
	m.ErrorStub(c, msg, fields...)
}

func (m mockLogService) Err(err error) log.Field {
	return m.ErrStub(err)
}

func (m mockLogService) Info(c context.Context, msg string, fields ...log.Field) {
	m.InfoStub(c, msg, fields...)
}

func (m mockLogService) Warn(c context.Context, msg string, fields ...log.Field) {
	m.WarnStub(c, msg, fields...)
}

func (m mockLogService) Any(key string, value interface{}) log.Field {
	return m.AnyStub(key, value)
}

func (m mockLogService) Warnln(ctx context.Context, msg string, fields ...log.Field) {
	m.WarnlnStub(ctx, msg, fields...)
}

// ValidationResultRepository
type mockValidationResultRepo struct {
	SaveStub func(ctx context.Context, paymentID string, validationResult storage.ValidationResult) error
}

func (m mockValidationResultRepo) Save(ctx context.Context, paymentID string, validationResult storage.ValidationResult) error {
	return m.SaveStub(ctx, paymentID, validationResult)
}

// Configurations.
type mockConfigurationsService struct {
	ports.Configurations
	LoadParametersValidationStub func(ctx context.Context) (*domain.ConfigParametersValidation, error)
}

func (m mockConfigurationsService) LoadParametersValidation(ctx context.Context) (*domain.ConfigParametersValidation, error) {
	return m.LoadParametersValidationStub(ctx)
}

// ReparationSearchRepository.
type mockReparationSearchRepo struct {
	ports.ReparationSearchRepository
	GetListByUserIDAndCreationPeriodStub func(
		userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error)

	GetByPaymentIDAndUserIDStub func(
		paymentID string, userID int64) (*storage.ReparationOutput, error)

	GetListByKycIentificationAndCreationPeriodStub func(
		kycIdentificationID string, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error)
}

func (m mockReparationSearchRepo) GetListByUserIDAndCreationPeriod(
	userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
	return m.GetListByUserIDAndCreationPeriodStub(userID, gteCreatedAt, lteCreatedAt)
}

func (m mockReparationSearchRepo) GetByPaymentIDAndUserID(paymentID string, userID int64) (*storage.ReparationOutput, error) {
	return m.GetByPaymentIDAndUserIDStub(paymentID, userID)
}

func (m mockReparationSearchRepo) GetListByKycIentificationAndCreationPeriod(
	kycIdentificationID string,
	gteCreatedAt,
	lteCreatedAt time.Time,
) ([]storage.ReparationOutput, error) {
	return m.GetListByKycIentificationAndCreationPeriodStub(kycIdentificationID, gteCreatedAt, lteCreatedAt)
}

// BlockedUserService
type mockBlockedUserService struct {
	ports.BlockedUserService
	CheckUserIsBlockedStub func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error)
}

func (m mockBlockedUserService) CheckUserIsBlocked(
	ctx context.Context,
	input domain.BlockedUserInput,
	wrapFields *field.WrappedFields,
) (bool, error) {
	return m.CheckUserIsBlockedStub(ctx, input, wrapFields)
}

// Policy.
type mockPolicy struct {
	EvaluateWithUserStub func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error)
}

func (m mockPolicy) EvaluateWithUser(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
	return m.EvaluateWithUserStub(ctx, userID)
}

// SearchHub.
type mockSearchHub struct {
	ports.SearchHub
	CountChargebackFromUserStub func(
		ctx context.Context,
		input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error)
	GetTransactionHistoryStub func(
		ctx context.Context,
		transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error)
}

func (m mockSearchHub) CountChargebackFromUser(
	ctx context.Context,
	input domain.SearchHubCountChargebackInput,
) (*domain.SearchHubCountChargebackOutput, error) {
	return m.CountChargebackFromUserStub(ctx, input)
}

func (m mockSearchHub) GetTransactionHistory(
	ctx context.Context,
	transactionHistoryInput domain.TransactionHistoryInput,
) (*domain.TransactionHistoryOutput, error) {
	return m.GetTransactionHistoryStub(ctx, transactionHistoryInput)
}

// EventValidationResult.
type mockEventValidationResult struct {
	PublishStub func(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error
}

func (m mockEventValidationResult) Publish(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
	return m.PublishStub(ctx, input, wrapFields)
}

// Score.
type mockScore struct {
	GetScoreStub func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error)
}

func (m mockScore) GetScore(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
	return m.GetScoreStub(ctx, input)
}

// Mocks:

func mockTransactionDataDefault(t *testing.T) (transactionData domain.TransactionData) {
	err := json.Unmarshal(_transactionData, &transactionData)
	assert.NoError(t, err)
	return
}

func mockParametersValidation(t *testing.T) (parameters *domain.ConfigParametersValidation) {
	err := json.Unmarshal(_parametersValidationStruct, &parameters)
	assert.NoError(t, err)
	return
}

func mockInputDefault(t *testing.T, siteID string) domain.ValidationInput {
	transactionData := mockTransactionDataDefault(t)
	if siteID == "MLA" {
		transactionData.SiteID = siteID
	}
	return domain.ValidationInput{
		PaymentID:       "100",
		UserID:          123,
		Type:            domain.TypeReverse,
		SiteID:          siteID,
		FaqID:           "faq_01",
		TransactionData: transactionData,
		UserKycData: domain.SearchKycOutput{
			KycIdentificationID: "abc",
			DateCreated:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Requester: domain.Requester{
			ClientApplication: "card-admin",
			ClientScope:       "beta",
		},
	}
}

func mockExpectedEligible(siteID string) storage.ValidationResult {
	var transactionEvaluated storage.TransactionEvaluated
	var score *float64
	switch siteID {
	case "MLB":
		transactionEvaluated = storage.TransactionEvaluated{
			Amount:   350.00,
			Currency: "BRL",
		}
		score = new(float64)
		*score = 0.5
	case "MLA":
		transactionEvaluated = storage.TransactionEvaluated{
			Amount:   18.99,
			Currency: "USD",
		}
	}

	return storage.ValidationResult{
		PaymentID:           "100",
		KycIdentificationID: "abc",
		UserID:              123,
		IsEligible:          true,
		DateCreated:         time.Now().UTC(),
		AuthorizationID:     "auth_mlb_test_aotorres_FyLRqgiKVMw9Vb8y737ppvNBsYpvkf",
		Type:                domain.TypeReverse,
		SiteID:              siteID,
		FaqID:               "faq_01",
		Requester: domain.Requester{
			ClientApplication: "card-admin",
			ClientScope:       "beta",
		},
		TransactionEvaluated: transactionEvaluated,
		Score:                score,
	}
}

func mockExpected_NotEligible_RuleMaxAmount() storage.ValidationResult {
	validationResult := mockExpectedEligible("MLB")
	validationResult.IsEligible = false
	validationResult.Reason = domain.Reason{
		domain.RuleMaxAmountReparation: &domain.ReasonResult{
			Actual:   603.65,
			Accepted: 500,
		},
	}

	return validationResult
}

func mockExpected_NotEligible_RuleStatusDetailAllowed() storage.ValidationResult {
	validationResult := mockExpectedEligible("MLB")
	validationResult.IsEligible = false
	validationResult.Reason = domain.Reason{
		domain.RuleStatusDetailAllowed: &domain.ReasonResult{
			Actual:   "accredited",
			Accepted: []string{"pending_capture"},
		},
	}

	return validationResult
}

func mockExpected_NotEligible_Rules() storage.ValidationResult {
	userIsBlocked := new(bool)
	*userIsBlocked = true

	score04 := new(float64)
	*score04 = 0.4

	validationResult := mockExpectedEligible("MLB")
	validationResult.IsEligible = false
	validationResult.Score = score04
	validationResult.Reason = domain.Reason{
		domain.RuleMaxAmountReparation: &domain.ReasonResult{
			Actual:   599.99999,
			Accepted: 500,
		},
		domain.RuleStatusDetailAllowed: &domain.ReasonResult{
			Actual:   "accredited",
			Accepted: []string{"pending_capture"},
		},
		domain.RuleQtyReparationPerPeriodDays: &domain.ReasonResult{
			Actual: 3,
			Accepted: domain.ReasonResultPerPeriod{
				Qty:        2,
				PeriodDays: 30,
			},
		},
		domain.RuleBlockeduser: &domain.ReasonResult{
			Actual: true,
		},
		domain.RuleRestrictions: &domain.ReasonResult{
			Actual: []string{"FRAUDE_CRUCE", "FRAUDE_TC"},
		},
		domain.RuleQtyChargebackPerPeriodDays: &domain.ReasonResult{
			Actual: 1,
			Accepted: domain.ReasonResultPerPeriod{
				Qty:        0,
				PeriodDays: 90,
			},
		},
		domain.RuleSubTypeAllowed: &domain.ReasonResult{
			Actual:   "purchase extracash",
			Accepted: []string{"purchase"},
		},
		domain.RuleMerchantBlockedlist: &domain.ReasonResult{
			Actual: "123Componentes",
		},
		domain.RuleUserLifetime: &domain.ReasonResult{
			Actual:   time.Now().UTC().AddDate(0, 0, 1).Truncate(time.Hour),
			Accepted: 10,
		},
		domain.RuleMinTotalHistoricalAmount: &domain.ReasonResult{
			Actual:   float64(0),
			Accepted: float64(100),
		},
		domain.RuleMinTransactionsQty: &domain.ReasonResult{
			Actual:   0,
			Accepted: 10,
		},
		domain.RuleScore: &domain.ReasonResult{
			Actual:   0.4,
			Accepted: 0.5,
		},
	}

	return validationResult
}
