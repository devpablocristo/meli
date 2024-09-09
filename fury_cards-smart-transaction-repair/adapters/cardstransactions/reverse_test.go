package cardstransactions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	pkgjson "github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_cardsTransactions_Reverse_Success(t *testing.T) {
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{
		ReadAllStub: io.ReadAll,
	}
	mockNewJson := pkgjson.NewJSON()
	mockNewLog := mockLogService{}

	mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := newMockReversalInput(t)

	mockCallReversalSuccess(input.TransactionData.Provider.ID)()

	expectedOutput := &domain.ReversalOutput{
		ReverseID: "smart_reverse_auth_mlb_test_aotorres_FyLRqgiKVMw9Vb8y737ppvNBsYpvkf",
	}

	output, err := mockReversalService.Reverse(context.TODO(), input)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func Test_cardsTransactions_Reverse_Error_CheckResult_MarshalResponse(t *testing.T) {
	t.Run("should return marshal error", func(t *testing.T) {
		mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
		mockNewIoutil := ioutilMock{
			ReadAllStub: io.ReadAll,
		}
		mockNewJson := jsonMock{
			UnmarshalStub: func(_ []byte, v interface{}) error {
				value := v.(*responseReversal)
				err := json.Unmarshal(_reversalResponseNoApproved, value)
				assert.NoError(t, err)
				return nil
			},
			MarshalStub: func(v interface{}) ([]byte, error) {
				_, castOk := v.(responseReversal)
				if castOk {
					return nil, errors.New("error marshal")
				}
				return _reversalRequest, nil
			},
		}
		mockNewLog := mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		}

		mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

		input := newMockReversalInput(t)

		mockCallReversalSuccess(input.TransactionData.Provider.ID)()

		output, err := mockReversalService.Reverse(context.TODO(), input)

		errorExpected := gnierrors.Wrap(domain.Forbidden,
			errors.New(`reversal not approved. unfortunately it was not possible to convert the reason (error marshal)`),
			"request transaction reversal failed",
			false)

		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
		assert.Empty(t, output)
	})
}

func Test_cardsTransactions_Reverse_Error_CheckResult_NotApproved(t *testing.T) {
	t.Run("should return declined result", func(t *testing.T) {
		mockNewJson := jsonMock{
			UnmarshalStub: func(_ []byte, v interface{}) error {
				value := v.(*responseReversal)
				err := json.Unmarshal(_reversalResponseNoApproved, value)
				assert.NoError(t, err)
				return nil
			},
			MarshalStub: func(v interface{}) ([]byte, error) {
				return json.Marshal(v)
			},
		}
		mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
		mockNewIoutil := ioutilMock{
			ReadAllStub: io.ReadAll,
		}
		mockNewLog := mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		}

		mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

		input := newMockReversalInput(t)

		mockCallReversalSuccess(input.TransactionData.Provider.ID)()

		output, err := mockReversalService.Reverse(context.TODO(), input)

		errorExpected := gnierrors.Wrap(domain.Forbidden,
			bodyResponseError(_reversalResponseNoApproved), "request transaction reversal failed", false)

		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
		assert.Empty(t, output)
	})
}

func Test_cardsTransactions_Reverse_Error_CheckResult_NotApproved_Authorization_Already_Reversed(t *testing.T) {
	var mockNewJson jsonMock
	var errorExpected error
	t.Run("should return declined result, and status_detail: authorization_already_reversed", func(t *testing.T) {
		mockNewJson = jsonMock{
			UnmarshalStub: func(_ []byte, v interface{}) error {
				value := v.(*responseReversal)
				err := json.Unmarshal(_reversalResponseNoApprovedWithAlreadyReversed, value)
				assert.NoError(t, err)
				return nil
			},
			MarshalStub: func(v interface{}) ([]byte, error) {
				return json.Marshal(v)
			},
		}

		errorExpected = gnierrors.Wrap(domain.AuthorizationAlreadyReversed,
			bodyResponseError(_reversalResponseNoApprovedWithAlreadyReversed), "request transaction reversal failed", false)

	})

	t.Run("should return declined result with error marshal in body, and status_detail: authorization_already_reversed", func(t *testing.T) {
		mockNewJson = jsonMock{
			UnmarshalStub: func(_ []byte, v interface{}) error {
				value := v.(*responseReversal)
				err := json.Unmarshal(_reversalResponseNoApprovedWithAlreadyReversed, value)
				assert.NoError(t, err)
				return nil
			},
			MarshalStub: func(v interface{}) ([]byte, error) {
				_, castOk := v.(responseReversal)
				if castOk {
					return nil, errors.New("{}")
				}
				return _reversalRequest, nil
			},
		}

		errorExpected = gnierrors.Wrap(domain.AuthorizationAlreadyReversed,
			bodyResponseError([]byte(string(`reversal not approved. unfortunately it was not possible to convert the reason ({})`))),
			"request transaction reversal failed", false)

	})

	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{
		ReadAllStub: io.ReadAll,
	}
	mockNewLog := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.NotEmpty(t, value)
			return log.Field{
				Key:    "body_reversal",
				String: value.(string),
			}
		},
	}

	mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := newMockReversalInput(t)

	mockCallReversalSuccess(input.TransactionData.Provider.ID)()

	output, err := mockReversalService.Reverse(context.TODO(), input)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Empty(t, output)
}

func Test_cardsTransactions_Reverse_Error_Unmarshal(t *testing.T) {
	mockNewJson := jsonMock{
		UnmarshalStub: func(_ []byte, v interface{}) error {
			return errors.New("error unmarshal")
		},
		MarshalStub: func(v interface{}) ([]byte, error) {
			return json.Marshal(v)
		},
	}
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{
		ReadAllStub: io.ReadAll,
	}
	mockNewLog := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.NotEmpty(t, value)
			return log.Field{
				Key:    "body_reversal",
				String: value.(string),
			}
		},
	}

	mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := newMockReversalInput(t)

	mockCallReversalSuccess(input.TransactionData.Provider.ID)()

	output, err := mockReversalService.Reverse(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.InternalServerError,
		errors.New("error unmarshal"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Empty(t, output)
}

func Test_cardsTransactions_Reverse_Error_ValidateResponse(t *testing.T) {
	t.Run("should return ioutil error", func(t *testing.T) {
		mockNewJson := jsonMock{
			UnmarshalStub: func(_ []byte, v interface{}) error {
				return errors.New("error unmarshal")
			},
			MarshalStub: func(v interface{}) ([]byte, error) {
				return json.Marshal(v)
			},
		}
		mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
		mockNewIoutil := ioutilMock{
			ReadAllStub: func(r io.Reader) ([]byte, error) {
				return nil, errors.New("error read")
			},
		}
		mockNewLog := mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		}

		mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

		input := newMockReversalInput(t)

		mockCallReversalSuccess(input.TransactionData.Provider.ID)()

		output, err := mockReversalService.Reverse(context.TODO(), input)

		errorExpected := gnierrors.Wrap(domain.BadGateway,
			errors.New("error read"), "request transaction reversal failed", false)

		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
		assert.Empty(t, output)
	})

	t.Run("should return response error", func(t *testing.T) {
		mockNewJson := pkgjson.NewJSON()
		mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
		mockNewIoutil := ioutilMock{
			ReadAllStub: io.ReadAll,
		}
		mockNewLog := mockLogService{
			AnyStub: func(key string, value interface{}) log.Field {
				assert.NotEmpty(t, value)
				return log.Field{
					Key:    "body_reversal",
					String: value.(string),
				}
			},
		}

		mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

		mockCallReversalBad()()

		input := newMockReversalInput(t)

		output, err := mockReversalService.Reverse(context.TODO(), input)

		errorExpected := gnierrors.Wrap(domain.BadRequest,
			bodyResponseError(_reversalResponseBadRequest), "request transaction reversal failed", false)

		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
		assert.Empty(t, output)
	})
}

func Test_cardsTransactions_Reverse_Error_Call_Post(t *testing.T) {
	mockNewJson := pkgjson.NewJSON()
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{}
	mockNewLog := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.NotEmpty(t, value)
			return log.Field{
				Key:    "body_reversal",
				String: value.(string),
			}
		},
	}

	mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	mockCallReversalFail()()

	input := newMockReversalInput(t)

	output, err := mockReversalService.Reverse(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("error read"), "request transaction reversal failed", false)

	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.Empty(t, output)
}

func Test_cardsTransactions_Reverse_Error_MarshalRequest(t *testing.T) {
	mockNewJson := jsonMock{
		MarshalStub: func(v interface{}) ([]byte, error) {
			return nil, errors.New("error marshal")
		},
	}
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{}
	mockNewLog := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.Empty(t, value)
			return log.Field{
				Key:    "body_reversal",
				String: value.(string),
			}
		},
	}

	mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := newMockReversalInput(t)

	output, err := mockReversalService.Reverse(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.InternalServerError,
		errors.New("error marshal"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Empty(t, output)
}

func Test_cardsTransactions_Reverse_Error_InvalidHeader(t *testing.T) {
	mockNewJson := jsonMock{}
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{}
	mockNewLog := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.Empty(t, value)
			return log.Field{
				Key:    "body_reversal",
				String: value.(string),
			}
		},
	}

	mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := domain.ReversalInput{
		TransactionData: &domain.TransactionData{},
	}

	output, err := mockReversalService.Reverse(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.BadRequest,
		errors.New("required headers: X-Client-Id, X-Site-Id, X-Environment, X-Provider"), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Empty(t, output)
}

func Test_cardsTransactions_Reverse_Error_InvalidBody(t *testing.T) {
	mockNewJson := jsonMock{}
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{}
	mockNewLog := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.Empty(t, value)
			return log.Field{
				Key:    "body_reversal",
				String: value.(string),
			}
		},
	}

	mockReversalService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := domain.ReversalInput{
		TransactionData: &domain.TransactionData{
			Provider:    domain.Provider{ID: "conductor"},
			Environment: "test",
			SiteID:      "MLB",
		},
		HeaderValueXClientID: "123456",
	}

	output, err := mockReversalService.Reverse(context.TODO(), input)

	msgErrExpected := fmt.Sprint(
		`required fields: `,
		`operations: {transmission_datetime, installments, `,
		`authorization: {id, transmission_datetime}, `,
		`transaction: {currency}, `,
		`card: {number_id, country}}`)

	errorExpected := gnierrors.Wrap(domain.BadRequest,
		errors.New(msgErrExpected), "request transaction reversal failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Empty(t, output)
}

func Test_validateRequest_Error_RequestFullEmpty(t *testing.T) {
	err := validateRequest(requestReversal{})

	msgErrExpected := fmt.Sprint(
		`required fields: id, `,
		`operations: {transmission_datetime, installments, `,
		`authorization: {id, transmission_datetime}, `,
		`transaction: {currency}, `,
		`card: {number_id, country}}, `,
		`provider: {id}, `,
		`options: {reversal_type}`)

	errorExpected := errors.New(msgErrExpected)

	assert.Equal(t, errorExpected, err)
}

func Test_cardsTransactions_getBaseURL(t *testing.T) {
	baseURLTransactions := domain.ConfigTransactionsApp{
		BaseURL:       url.URL{Host: "http://melisystems.com"},
		BaseURLLegacy: url.URL{Host: "http://mercadopago.com"},
		ProvidersWithLegacyURL: map[string]struct{}{
			"globalprocessing": {},
			"evertec":          {},
		},
	}
	c := cardsTransactions{
		config: NewConfig(baseURLTransactions),
	}

	t.Run("url_legacy", func(t *testing.T) {
		transaction := domain.TransactionData{
			Provider: domain.Provider{ID: "globalprocessing"},
			SiteID:   "MLA",
		}

		baseURL := c.getBaseURL(transaction.Provider.ID)

		assert.Equal(t, "http://mercadopago.com", baseURL.Host)
	})

	t.Run("url", func(t *testing.T) {
		transaction := domain.TransactionData{
			Provider: domain.Provider{ID: "bari"},
			SiteID:   "MLA",
		}

		baseURL := c.getBaseURL(transaction.Provider.ID)

		assert.Equal(t, "http://melisystems.com", baseURL.Host)
	})
}
