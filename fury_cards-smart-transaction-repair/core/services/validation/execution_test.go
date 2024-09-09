package validation

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_eligibilityValidation_ExecuteValidation_Success(t *testing.T) {
	tests := []struct {
		name       string
		siteID     string
		parameters *domain.ConfigParametersValidation
	}{
		{
			name:       "local currency",
			siteID:     "MLB",
			parameters: mockParametersValidation(t),
		},
		{
			name:       "international currency",
			siteID:     "MLA",
			parameters: mockParametersValidation(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockValidationResultRepo := mockValidationResultRepo{}

			mockConfigurationsService := mockConfigurationsService{
				LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
					return mockParametersValidation(t), nil
				},
			}

			mockReparationSearchRepo := mockReparationSearchRepo{
				GetListByUserIDAndCreationPeriodStub: func(
					userID int64,
					gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {

					assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
					assert.Equal(t, int64(123), userID)
					return nil, nil
				},
			}

			mockLogService := mockLogService{
				InfoStub: func(c context.Context, msg string, fields ...log.Field) {
					rulesApplied := fields[0].Interface.(ruleLog)
					assert.Len(t, rulesApplied["100|"+tt.siteID], 11)
				},
				AnyStub: func(key string, value interface{}) log.Field {
					assert.NotEmpty(t, key)
					return log.Field{
						Key:       key,
						Interface: value,
					}
				},
			}

			mockBlockedUserService := mockBlockedUserService{
				CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
					return false, nil
				},
			}

			mockPolicy := mockPolicy{
				EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
					return &domain.PolicyOutput{
						IsAuthorized: true,
					}, nil
				},
			}

			mockSearchHub := mockSearchHub{
				CountChargebackFromUserStub: func(
					ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
					return &domain.SearchHubCountChargebackOutput{
						Total: 0,
					}, nil
				},
				GetTransactionHistoryStub: func(
					ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
					return &domain.TransactionHistoryOutput{
						IssuerAccounts: []domain.HistoryIssuerAccounts{
							{
								MatchTotal:   11,
								Transactions: []domain.HistoryOperation{{Settlement: domain.Settlement{Amount: 101, Currency: "USD"}}},
							},
						},
					}, nil
				},
			}

			mockEventValidationResult := mockEventValidationResult{
				PublishStub: func(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
					expectedValidationResult := mockExpectedEligible(tt.siteID)

					utilstest.AssertStructEqual(t, expectedValidationResult, input)
					assert.Nil(t, input.Reason)
					return nil
				},
			}

			mockScore := mockScore{
				GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
					return &domain.ScoreOutput{Score: 0.5}, nil
				},
			}

			eligibilityValidation := New(
				mockLogService,
				mockConfigurationsService,
				mockValidationResultRepo,
				mockReparationSearchRepo,
				mockBlockedUserService,
				mockPolicy,
				mockSearchHub,
				mockEventValidationResult,
				mockScore,
			)

			input := mockInputDefault(t, tt.siteID)

			err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

			assert.NoError(t, err)
		})
	}
}

func Test_eligibilityValidation_ExecuteValidation_Success_WithErrorPublishing(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		WarnlnStub: func(ctx context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "unpublished_event_validation_result_by_paymentID: 100", msg)
			assert.Equal(t, errors.New("error saving"), fields[0].Interface)
		},
		ErrStub: func(err error) log.Field {
			return log.Field{Interface: err}
		},
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 0,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{
						MatchTotal:   11,
						Transactions: []domain.HistoryOperation{{Settlement: domain.Settlement{Amount: 101, Currency: "USD"}}},
					},
				},
			}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{
		PublishStub: func(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
			utilstest.AssertStructEqual(t, mockExpectedEligible("MLB"), input)
			assert.Nil(t, input.Reason)
			return errors.New("error saving")
		},
	}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	assert.NoError(t, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_NotEligible_MaxAmountReparation(t *testing.T) {
	expectedValidation := mockExpected_NotEligible_RuleMaxAmount()

	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 0,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{
						MatchTotal:   11,
						Transactions: []domain.HistoryOperation{{Settlement: domain.Settlement{Amount: 101, Currency: "USD"}}},
					},
				},
			}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{
		PublishStub: func(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
			utilstest.AssertStructEqual(t, expectedValidation, input)
			utilstest.AssertFieldsNoEmptyFromStruct(t, input)
			return nil
		},
	}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")
	// Forced MaxAmountReparation
	input.TransactionData.Operation.Billing.Amount = 60365
	expectedValidation.TransactionEvaluated.Amount = 603.65

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	b, _ := json.Marshal(expectedValidation.Reason)
	errorExpected := gnierrors.Wrap(domain.NotEligible, errors.New(string(b)), "validation result", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)

	errResponseApi := gnierrors.Cause(err).(*domain.ValidationCauseErrResponse)
	utilstest.AssertFieldsNoEmptyFromStruct(t, *errResponseApi)
}

func Test_eligibilityValidation_ExecuteValidation_Error_NotEligible_StatusDetailAllowed(t *testing.T) {
	expectedValidation := mockExpected_NotEligible_RuleStatusDetailAllowed()

	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 0,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{
				IssuerAccounts: []domain.HistoryIssuerAccounts{
					{
						MatchTotal:   11,
						Transactions: []domain.HistoryOperation{{Settlement: domain.Settlement{Amount: 101, Currency: "USD"}}},
					},
				},
			}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{
		PublishStub: func(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
			utilstest.AssertStructEqual(t, expectedValidation, input)
			utilstest.AssertFieldsNoEmptyFromStruct(t, input)
			return nil
		},
	}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")
	// Forced StatusDetail
	input.TransactionData.StatusDetail = "accredited"

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	b, _ := json.Marshal(expectedValidation.Reason)
	errorExpected := gnierrors.Wrap(domain.NotEligible, errors.New(string(b)), "validation result", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)

	errResponseApi := gnierrors.Cause(err).(*domain.ValidationCauseErrResponse)
	utilstest.AssertFieldsNoEmptyFromStruct(t, *errResponseApi)
}

func Test_eligibilityValidation_ExecuteValidation_Error_NotEligible_Full(t *testing.T) {
	expectedValidation := mockExpected_NotEligible_Rules()

	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return []storage.ReparationOutput{
				{},
				{},
				{},
			}, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return true, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized:    false,
				RestrictsFailed: []string{"FRAUDE_CRUCE", "FRAUDE_TC"},
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 1,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{
		PublishStub: func(ctx context.Context, input storage.ValidationResult, wrapFields *field.WrappedFields) error {
			utilstest.AssertStructEqual(t, expectedValidation, input)
			utilstest.AssertFieldsNoEmptyFromStruct(t, input, "ApprovedData")
			return nil
		},
	}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.4}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")
	// Forced StatusDetail
	input.TransactionData.StatusDetail = "accredited"

	// Forced MaxAmountReparation
	input.TransactionData.Operation.Billing.Amount = 59999999
	input.TransactionData.Operation.Billing.DecimalDigits = 5
	expectedValidation.TransactionEvaluated.Amount = 599.99999

	// Forced StatusDetail
	input.TransactionData.StatusDetail = "accredited"

	// Forced SubType
	input.TransactionData.Operation.SubType = "purchase extracash"

	// Forced Merchant
	input.TransactionData.Operation.Authorization.CardAcceptor.Name = "123Componentes"

	// Forced UserLifetime
	input.UserKycData.DateCreated = time.Now().UTC().AddDate(0, 0, 1).Truncate(time.Hour)

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	b, _ := json.Marshal(expectedValidation.Reason)
	errorExpected := gnierrors.Wrap(domain.NotEligible,
		errors.New(string(b)),
		"validation result", false)

	assert.Error(t, err)
	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.Equal(t, gnierrors.Cause(errorExpected).Error(), gnierrors.Cause(err).Error())

	errResponseApi := gnierrors.Cause(err).(*domain.ValidationCauseErrResponse)
	utilstest.AssertFieldsNoEmptyFromStruct(t, *errResponseApi)
}

func Test_eligibilityValidation_ExecuteValidation_Error_GetParametersValidation(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return nil, gnierrors.Wrap(domain.BadGateway,
				errors.New(`not-found`), "error in configurations", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{}

	mockPolicy := mockPolicy{}

	mockSearchHub := mockSearchHub{}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`error in configurations : not-found`), "error in configurations", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_GetParametersValidation_Rules_Not_Found(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(_ context.Context) (*domain.ConfigParametersValidation, error) {
			return &domain.ConfigParametersValidation{
				RulesSite: map[string]*domain.ConfigRulesSite{
					"XXX": {},
				},
			}, nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{}

	mockPolicy := mockPolicy{}

	mockSearchHub := mockSearchHub{}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`rules not found for site(MLB) - parameters not configured`), "rules not found", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_LoadedData_TransactionRepair(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("unknown error"), "transaction repair document search failure", false)
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, gnierrors.Wrap(domain.BadGateway, errors.New("unknown error"), "blockedlist database failure", false)
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 0,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`transaction repair document search failure : unknown error`), "transaction repair document search failure", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_LoadedData_BlockedUser(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, gnierrors.Wrap(domain.BadGateway, errors.New("unknown error"), "blockedlist database failure", false)
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 0,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`unknown error`), "blockedlist database failure", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_LoadedData_GetPolicyAgentEvaluation(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(_ context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("unknown error"), "policy failure", false)
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 0,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`unknown error`), "policy failure", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_LoadedData_CountChargebackFromUser(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("unknown error"), "search-hub failure", false)
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`unknown error`), "search-hub failure", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_LoadedData_GetTransactionHistory(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("unknown error transaction history"), "search-hub failure", false)
		},
	}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`unknown error transaction history`), "search-hub failure", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_Error_LoadedData_GetScore(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(ctx context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			assert.Equal(t, gteCreatedAt, lteCreatedAt.Add(-time.Hour*24*30))
			assert.Equal(t, int64(123), userID)
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("unknown error score"), "service_score failure", false)
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`unknown error score`), "service_score failure", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_eligibilityValidation_ExecuteValidation_RunPrerequisitesForRuleValidation_Error(t *testing.T) {
	mockValidationResultRepo := mockValidationResultRepo{}

	mockConfigurationsService := mockConfigurationsService{
		LoadParametersValidationStub: func(_ context.Context) (*domain.ConfigParametersValidation, error) {
			return mockParametersValidation(t), nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetListByUserIDAndCreationPeriodStub: func(userID int64, gteCreatedAt, lteCreatedAt time.Time) ([]storage.ReparationOutput, error) {
			return nil, nil
		},
	}

	mockLogService := mockLogService{
		InfoStub: func(c context.Context, msg string, fields ...log.Field) {},
		AnyStub: func(key string, value interface{}) log.Field {
			return log.Field{Key: key, Interface: value}
		},
	}

	mockBlockedUserService := mockBlockedUserService{
		CheckUserIsBlockedStub: func(ctx context.Context, input domain.BlockedUserInput, wrapFields *field.WrappedFields) (bool, error) {
			return false, nil
		},
	}

	mockPolicy := mockPolicy{
		EvaluateWithUserStub: func(ctx context.Context, userID int64, tagName ...string) (*domain.PolicyOutput, error) {
			return &domain.PolicyOutput{
				IsAuthorized: true,
			}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		CountChargebackFromUserStub: func(
			ctx context.Context, input domain.SearchHubCountChargebackInput) (*domain.SearchHubCountChargebackOutput, error) {
			return &domain.SearchHubCountChargebackOutput{
				Total: 0,
			}, nil
		},
		GetTransactionHistoryStub: func(
			ctx context.Context, transactionHistoryInput domain.TransactionHistoryInput) (*domain.TransactionHistoryOutput, error) {
			return &domain.TransactionHistoryOutput{}, nil
		},
	}

	mockEventValidationResult := mockEventValidationResult{}

	mockScore := mockScore{
		GetScoreStub: func(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
			return &domain.ScoreOutput{Score: 0.5}, nil
		},
	}

	eligibilityValidation := New(
		mockLogService,
		mockConfigurationsService,
		mockValidationResultRepo,
		mockReparationSearchRepo,
		mockBlockedUserService,
		mockPolicy,
		mockSearchHub,
		mockEventValidationResult,
		mockScore,
	)

	input := mockInputDefault(t, "MLB")
	// Forced other currency
	input.TransactionData.Operation.Billing.Currency = "CRU"

	err := eligibilityValidation.ExecuteValidation(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.UnprocessableEntity,
		errors.New(`transaction billing currency(CRU) different from the parameters(BRL)`), "validation failed", false)

	assert.Error(t, err)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_fillTransactionEvaluated_WithoutConfig(t *testing.T) {
	t.Run("should return nil to storageTransactionEvaluated", func(t *testing.T) {
		storageTransactionEvaluated := fillTransactionEvaluated(domain.ValidationInput{}, nil)

		assert.Empty(t, storageTransactionEvaluated)
	})
}
