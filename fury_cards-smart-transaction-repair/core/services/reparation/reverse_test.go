package reparation

import (
	"context"
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

func Test_reparation_Reverse_Success(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{
		SaveStub: func(ctx context.Context, transactionRepair storage.ReparationInput) error {
			utilstest.AssertStructEqual(t, mockBaseTransactionRepairExpected(), transactionRepair.BaseReparation)
			utilstest.AssertFieldsNoEmptyFromStruct(t, transactionRepair.BaseReparation)
			return nil
		},
	}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			_ *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: mockTransactionDataDefault(t),
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{
		ReverseStub: func(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
			assert.Equal(t, "123456", input.HeaderValueXClientID)
			utilstest.AssertFieldsNoEmptyFromStruct(t, input)

			return &domain.ReversalOutput{
				ReverseID: "reverse_auth_mlb_test_aotorres_FyLRqgiKVMw9Vb8y737ppvNBsYpvkf",
			}, nil
		},
	}

	mockValidationService := mockValidationService{
		ExecuteValidationStub: func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
			utilstest.AssertStructEqual(t, mockValidationInputExpected(t), input)
			utilstest.AssertFieldsNoEmptyFromStruct(t, input)
			return nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "adapter..", false)
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(
		context.TODO(),
		domain.ReverseTransactionInput{
			PaymentID:            "100",
			UserID:               123,
			SiteID:               "MLB",
			HeaderValueXClientID: "123456",
			FaqID:                "faq_01",
			Requester: domain.Requester{
				ClientApplication: "card-admin",
				ClientScope:       "beta",
			}},
		field.NewWrappedFields(),
	)

	assert.NoError(t, err)
}

func Test_reparation_Reverse_Success_WithError_GetTransactionReparation(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{
		SaveStub: func(ctx context.Context, transactionRepair storage.ReparationInput) error {
			assert.Equal(t, int64(123), transactionRepair.UserID)
			assert.Equal(t, "100", transactionRepair.PaymentID)
			return nil
		},
	}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{Type: "authorization"},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{
		ReverseStub: func(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
			assert.NotNil(t, input.TransactionData)
			return &domain.ReversalOutput{
				ReverseID: "ReverseID",
			}, nil
		},
	}

	mockValidationService := mockValidationService{
		ExecuteValidationStub: func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
			return nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.BadGateway, errors.New("error getting"), "transaction repair document search failure", false)
		},
	}

	mockLogService := mockLogService{
		WarnStub: func(c context.Context, msg string, fields ...log.Field) {
			errExpected := gnierrors.Wrap(domain.BadGateway, errors.New("error getting"), "transaction repair document search failure", false)
			assert.Equal(t, "unrecovered transaction repair, by paymentID: 100", msg)
			assert.Equal(t, errExpected.Error(), fields[0].Interface.(error).Error())
		},
		ErrStub: func(err error) log.Field {
			return log.Field{Interface: err}
		},
	}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.NoError(t, err)
}

func Test_reparation_Reverse_Success_WithError_SaveTransactionReparation(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{
		SaveStub: func(ctx context.Context, transactionRepair storage.ReparationInput) error {
			return errors.New("error saving")
		},
	}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{Type: "authorization"},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{
		ReverseStub: func(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
			assert.NotNil(t, input.TransactionData)
			return &domain.ReversalOutput{
				ReverseID: "ReverseID",
			}, nil
		},
	}

	mockValidationService := mockValidationService{
		ExecuteValidationStub: func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
			return nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "adapter..", false)
		},
	}

	mockLogService := mockLogService{
		WarnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "could not insert transaction repair, by paymentID: 100", msg)
			assert.Equal(t, errors.New("error saving"), fields[0].Interface)
		},
		ErrStub: func(err error) log.Field {
			return log.Field{Interface: err}
		},
	}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.NoError(t, err)
}

func Test_reparation_Reverse_Error_ReverseTransactionReparation(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{Type: "authorization"},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{
		ReverseStub: func(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
			attrErr := gnierrors.Attr{
				Key:   "body_reversal",
				Value: `{"error":"unauthorized"}`,
			}
			return nil, gnierrors.Wrap(domain.Unauthorized,
				errors.New(`request unauthorized`), "request transaction reversal failed", false, attrErr)
		},
	}

	mockValidationService := mockValidationService{
		ExecuteValidationStub: func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
			return nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "adapter..", false)
		},
	}

	mockLogService := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			switch key {
			case string(domain.KeyFieldOriginalErrType):
				assert.Equal(t, value, "unauthorized")
			case string(domain.KeyFieldPaymentID):
				assert.Equal(t, value, "100")
			}

			return log.Field{Key: key, Interface: value}
		},
	}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`request transaction reversal failed : request unauthorized`), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_ReverseAuthorizationAlreadyReversed(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{Type: "authorization"},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{
		ReverseStub: func(ctx context.Context, input domain.ReversalInput) (*domain.ReversalOutput, error) {
			return nil, gnierrors.Wrap(domain.AuthorizationAlreadyReversed,
				errors.New(`{"error":{"result":{"status":"declined","status_detail":"authorization_already_reversed"}}}`),
				"request transaction reversal failed",
				false)
		},
	}

	mockValidationService := mockValidationService{
		ExecuteValidationStub: func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
			return nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "adapter..", false)
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.AuthorizationAlreadyReversed,
		errors.New(
			`request transaction reversal failed : {"error":{"result":{"status":"declined","status_detail":"authorization_already_reversed"}}}`),
		"invalid request",
		false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_ExecuteValidation_NotEligible(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{Type: "authorization"},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{}

	mockValidationService := mockValidationService{
		ExecuteValidationStub: func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
			return gnierrors.Wrap(domain.NotEligible,
				errors.New(`customer not eligible for reversal`), "validation result", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "adapter..", false)
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.NotEligible,
		errors.New(`customer not eligible for reversal`), "validation result", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_ExecuteValidation_Unknown(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{Type: "authorization"},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{}

	mockValidationService := mockValidationService{
		ExecuteValidationStub: func(ctx context.Context, input domain.ValidationInput, wrapFields *field.WrappedFields) error {
			return gnierrors.Wrap(domain.BadGateway,
				errors.New(`unknown`), "validation result", false)
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "adapter..", false)
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`unknown`), "validation result", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_GetTransaction(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{}, gnierrors.Wrap(domain.Unauthorized,
				errors.New(`unauthorized cards search`), "request transaction search failed", false)
		},
	}

	mockCardsTransactions := mockCardsTransactions{}
	mockValidationService := mockValidationService{}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, gnierrors.Wrap(domain.NotFound, errors.New("not found"), "adapter..", false)
		},
	}

	mockLogService := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			switch key {
			case string(domain.KeyFieldOriginalErrType):
				assert.Equal(t, value, "unauthorized")
			case string(domain.KeyFieldPaymentID):
				assert.Equal(t, value, "100")
			}

			return log.Field{Key: key, Interface: value}
		},
	}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`request transaction search failed : unauthorized cards search`), "request transaction search failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_GetUserKyc(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{}
	mockValidationService := mockValidationService{}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return &storage.ReparationOutput{}, nil
		},
	}

	mockLogService := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			switch key {
			case string(domain.KeyFieldOriginalErrType):
				assert.Equal(t, value, "unprocessable_entity")
			default:
				panic(key)
			}

			return log.Field{Key: key, Interface: value}
		},
	}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return nil, gnierrors.Wrap(domain.UnprocessableEntity,
				errors.New(`{"errors":[{"message":"timeout"}]`),
				"graphql-search-kyc get_user_kyc failure", false)
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`graphql-search-kyc get_user_kyc failure : {"errors":[{"message":"timeout"}]`),
		"graphql-search-kyc get_user_kyc failure", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_GetUserKyc_UnableUserNotAvailable(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{},
			}, nil
		},
	}

	mockCardsTransactions := mockCardsTransactions{}
	mockValidationService := mockValidationService{}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return &storage.ReparationOutput{}, nil
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return nil, gnierrors.Wrap(domain.UnavailableForLegalReasons,
				errors.New(`{"errors":[{"message":"The requested user is not available due to legal reasons"}]`),
				"graphql-search-kyc get_user_kyc failure", false)
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.UnableRequest,
		errors.New(
			`graphql-search-kyc get_user_kyc failure : {"errors":[{"message":"The requested user is not available due to legal reasons"}]`),
		"invalid request", false,
	)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_Repair_Already_Done(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}
	mockCardsTransactions := mockCardsTransactions{}
	mockValidationService := mockValidationService{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{},
			}, nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return &storage.ReparationOutput{}, nil
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.AlreadyRepaired,
		errors.New(`reverse already done`), "invalid request", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_PaymentIDEmpty(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}
	mockCardsTransactions := mockCardsTransactions{}
	mockValidationService := mockValidationService{}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return &storage.ReparationOutput{}, nil
		},
	}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{},
			}, nil
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: ""}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.BadRequest,
		errors.New(`required fields: PaymentID; UserID;`), "invalid request", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_KycIdentificationIsEmpty(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}
	mockCardsTransactions := mockCardsTransactions{}
	mockValidationService := mockValidationService{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{},
			}, nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, nil
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.UnableRequest,
		errors.New(`kyc_identification from userID:123 is empty`), "invalid request", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_Reverse_Error_TransactionTypeNotAuthorization(t *testing.T) {
	mockTransactionRepairRepo := mockTransactionRepairRepo{}
	mockCardsTransactions := mockCardsTransactions{}
	mockValidationService := mockValidationService{}

	mockSearchHub := mockSearchHub{
		GetPaymentTransactionStub: func(
			ctx context.Context,
			transactionInput domain.TransactionInput,
			wrapFields *field.WrappedFields) (domain.TransactionOutput, error) {
			return domain.TransactionOutput{
				TransactionData: &domain.TransactionData{Type: "capture"},
			}, nil
		},
	}

	mockReparationSearchRepo := mockReparationSearchRepo{
		GetByPaymentIDAndUserIDStub: func(paymentID string, userID int64) (*storage.ReparationOutput, error) {
			return nil, nil
		},
	}

	mockLogService := mockLogService{}

	mockSearchKyc := mockSearchKyc{
		GetUserKycStub: func(ctx context.Context, input domain.SearchKycInput) (*domain.SearchKycOutput, error) {
			return &domain.SearchKycOutput{KycIdentificationID: "abc", DateCreated: time.Date(2022, 12, 28, 0, 0, 0, 0, time.UTC)}, nil
		},
	}

	reparationMock := New(
		mockLogService,
		mockValidationService,
		mockTransactionRepairRepo,
		mockSearchHub,
		mockCardsTransactions,
		mockReparationSearchRepo,
		mockSearchKyc,
	)

	err := reparationMock.Reverse(context.TODO(), domain.ReverseTransactionInput{PaymentID: "100", UserID: 123}, field.NewWrappedFields())

	assert.Error(t, err)

	errorExpected := gnierrors.Wrap(domain.UnableRequest,
		errors.New(`type of operation cannot be:capture`), "invalid request", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_buildError_unmapped_error_type(t *testing.T) {
	gnierr := gnierrors.Wrap(domain.Locked, errors.New("locked"), "locked", false)

	mockLogService := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.Equal(t, key, string(domain.KeyFieldOriginalErrType))
			assert.Equal(t, value, "locked")
			return log.Field{Key: key, Interface: value}
		},
		WarnStub: func(c context.Context, msg string, fields ...log.Field) {
			assert.Equal(t, "unmapped error type: locked", msg)
		},
	}

	err := reBuildError(context.TODO(), mockLogService, domain.Locked.String(), gnierr)

	errorExpected := gnierrors.Wrap(domain.BadGateway, errors.New("locked : locked"), "locked", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_reparation_calculatePaymentAgeInDays(t *testing.T) {
	hoursOneDay := 24

	transactionData := &domain.TransactionData{
		Operation: domain.Operation{
			CreationDatetime: time.Date(2022, 10, 18, 0, 0, 0, 0, time.UTC),
			PaymentAge:       int(time.Now().UTC().Sub(time.Date(2022, 10, 18, 0, 0, 0, 0, time.UTC)).Hours()) / hoursOneDay,
		},
	}

	age := transactionData.Operation.PaymentAge

	againstProof := time.Now().UTC().Add(-time.Duration(age) * time.Hour * 24)

	assert.Equal(t, againstProof.Day(), transactionData.Operation.CreationDatetime.Day())
	assert.Equal(t, againstProof.Month(), transactionData.Operation.CreationDatetime.Month())
	assert.Equal(t, againstProof.Year(), transactionData.Operation.CreationDatetime.Year())
}
