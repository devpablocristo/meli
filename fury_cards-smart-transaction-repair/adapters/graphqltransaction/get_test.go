package graphqltransaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	jsonpkg "github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func Test_transactionSearch_GetTransactionByPaymentID_Success(t *testing.T) {
	transactionSearchMock := newMockTest()

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	postMockFunc := func() func() {
		urlExpected := "https://urltest.com/cards/transactions/search"

		return apitest.NewMock().
			Post("https://urltest.com/cards/transactions/search").
			Header("Content-Type", "application/json").
			Body(string(_transactionRequest)).
			AddMatcher(func(r *http.Request, mr *apitest.MockRequest) error {
				if urlExpected != r.URL.String() {
					return fmt.Errorf("received url %s did not match url expected: %s", r.URL.String(), urlExpected)
				}
				return nil
			}).
			RespondWith().
			Body(string(_transactionResponse)).
			Status(http.StatusOK).
			EndStandalone()
	}

	postMockFunc()

	var outputExpected *domain.TransactionOutput
	errExpected := transactionSearchMock.json.Unmarshal(_transactionData, &outputExpected)
	outputExpected.TransactionData.Operation.PaymentAge = calculatePaymentAgeInDays(outputExpected.TransactionData.Operation.CreationDatetime)

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	assert.NoError(t, errExpected)
	assert.NoError(t, err)
	assert.Equal(t, outputExpected, output)
}

func Test_transactionSearch_GetTransactionByPaymentID_Error_Response(t *testing.T) {
	transactionSearchMock := newMockTest()

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	postMockFunc := func() func() {
		urlExpected := "https://urltest.com/cards/transactions/search"

		return apitest.NewMock().
			Post("https://urltest.com/cards/transactions/search").
			Header("Content-Type", "application/json").
			Body(string(_transactionRequest)).
			AddMatcher(func(r *http.Request, mr *apitest.MockRequest) error {
				if urlExpected != r.URL.String() {
					return fmt.Errorf("received url %s did not match url expected: %s", r.URL.String(), urlExpected)
				}
				return nil
			}).
			RespondWith().
			Body(string(_responseErr)).
			Status(http.StatusOK).
			EndStandalone()
	}

	postMockFunc()

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`[{"message":"bad request when try to get payment id '52525109477A'","path":["payment"]}]`),
		"request transaction search failed", false)

	assert.Empty(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_transactionSearch_GetTransactionByPaymentID_Error_Marshal_Response(t *testing.T) {
	transactionSearchMock := newMockTest()
	transactionSearchMock.json = jsonMock{
		UnmarshalStub: func(data []byte, v interface{}) error {
			return json.Unmarshal(data, v)
		},
		MarshalStub: func(v interface{}) ([]byte, error) {
			_, ok := v.([]errResponse)
			if ok {
				return nil, errors.New("err marshal")
			}

			return json.Marshal(v)
		},
	}

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	postMockFunc := func() func() {
		urlExpected := "https://urltest.com/cards/transactions/search"

		return apitest.NewMock().
			Post("https://urltest.com/cards/transactions/search").
			Header("Content-Type", "application/json").
			Body(string(_transactionRequest)).
			AddMatcher(func(r *http.Request, mr *apitest.MockRequest) error {
				if urlExpected != r.URL.String() {
					return fmt.Errorf("received url %s did not match url expected: %s", r.URL.String(), urlExpected)
				}
				return nil
			}).
			RespondWith().
			Body(string(_responseErr)).
			Status(http.StatusOK).
			EndStandalone()
	}

	postMockFunc()

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	errorExpected := gnierrors.Wrap(domain.InternalServerError,
		errors.New(`err marshal`), "request transaction search failed", false)

	assert.Empty(t, output)
	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
}

func Test_transactionSearch_GetTransactionByPaymentID_Error_Body(t *testing.T) {
	transactionSearchMock := newMockTest()

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	postMockFunc := func() func() {
		return apitest.NewMock().
			Post("/search").
			Header("Content-Type", "application/json").
			Body(string(_transactionRequest)).
			RespondWith().
			Body(`{"errors":[{"message":"error"}]}`).
			Status(http.StatusBadRequest).
			EndStandalone()
	}

	postMockFunc()

	errorExpected := gnierrors.Wrap(domain.BadRequest, errors.New(`{"errors":[{"message":"error"}]}`), "request transaction search failed", false)

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.Equal(t, gnierrors.Cause(errorExpected).Error(), gnierrors.Cause(err).Error())
	assert.Empty(t, output)
}

func Test_transactionSearch_GetTransactionByPaymentID_ErrorReadResponse(t *testing.T) {
	// Forced error in RealAll
	transactionSearchMock := transactionSearch{
		config:     &Config{BaseURL: *mockBaseURL()},
		httpClient: rest.NewHTTPWithTimeout(10 * time.Second),
		ioutil: ioutilMock{
			ReadAllStub: func(r io.Reader) ([]byte, error) { return nil, errors.New("error read") },
		},
		json: jsonpkg.NewJSON(),
	}

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	postMockFunc := func() func() {
		return apitest.NewMock().
			Post("/search").
			Header("Content-Type", "application/json").
			Body(string(_transactionRequest)).
			RespondWith().
			Body(`{}`).
			Status(http.StatusOK).
			EndStandalone()
	}

	postMockFunc()

	errorExpected := gnierrors.Wrap(domain.BadGateway, errors.New(`error read`), "request transaction search failed", false)

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.Equal(t, gnierrors.Cause(errorExpected).Error(), gnierrors.Cause(err).Error())
	assert.Empty(t, output)
}

func Test_transactionSearch_GetTransactionByPaymentID_ErrorUnmarshalResponse(t *testing.T) {
	transactionSearchMock := newMockTest()

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	postMockFunc := func() func() {
		return apitest.NewMock().
			Post("/search").
			Header("Content-Type", "application/json").
			Body(string(_transactionRequest)).
			RespondWith().
			Body(`{`).
			Status(http.StatusOK).
			EndStandalone()
	}

	postMockFunc()

	errorExpected := gnierrors.Wrap(domain.InternalServerError, errors.New(`unexpected end of JSON input`), "request transaction search failed", false)

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.Equal(t, gnierrors.Cause(errorExpected).Error(), gnierrors.Cause(err).Error())
	assert.Empty(t, output)
}

func Test_transactionSearch_GetTransactionByPaymentID_ErrorHTTP(t *testing.T) {
	transactionSearchMock := newMockTest()

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	postMockFunc := func() func() {
		return apitest.NewMock().
			Post("/xxx").
			Header("Content-Type", "application/json").
			Body(string(_transactionRequest)).
			RespondWith().
			EndStandalone()
	}

	postMockFunc()

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	assert.Error(t, err)
	assert.Equal(t, domain.BadGateway, gnierrors.Type(err))
	assert.Equal(t, "request transaction search failed", gnierrors.Message(err))
	assert.Empty(t, output)
}

func Test_transactionSearch_GetTransactionByPaymentID_ErrorMarshal(t *testing.T) {
	// Forced error in Marshal
	transactionSearchMock := transactionSearch{
		json: jsonMock{
			MarshalStub: func(v interface{}) ([]byte, error) {
				return nil, errors.New("error marshal")
			},
		},
	}

	input := &domain.TransactionInput{
		PaymentID: "50844833020",
	}

	errorExpected := gnierrors.Wrap(domain.InternalServerError, errors.New(`error marshal`), "request transaction search failed", false)

	output, err := transactionSearchMock.GetTransactionByPaymentID(context.TODO(), input, field.NewWrappedFields())

	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.Equal(t, gnierrors.Cause(errorExpected).Error(), gnierrors.Cause(err).Error())
	assert.Empty(t, output)
}
