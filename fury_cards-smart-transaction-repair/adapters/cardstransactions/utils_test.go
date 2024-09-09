package cardstransactions

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/transaction-data.json
var _transactionData []byte

//go:embed testdata/reversal-request.json
var _reversalRequest []byte

//go:embed testdata/reversal-response.json
var _reversalResponse []byte

//go:embed testdata/reversal-response-no-approved.json
var _reversalResponseNoApproved []byte

//go:embed testdata/reversal-response-no-approved_with_already_reversed.json
var _reversalResponseNoApprovedWithAlreadyReversed []byte

//go:embed testdata/reversal-response-badrequest.json
var _reversalResponseBadRequest []byte

// IoutilService
type ioutilMock struct {
	ReadAllStub func(r io.Reader) ([]byte, error)
}

func (i ioutilMock) ReadAll(r io.Reader) ([]byte, error) {
	return i.ReadAllStub(r)
}

// JSONService
type jsonMock struct {
	UnmarshalStub func(data []byte, v interface{}) error
	MarshalStub   func(v interface{}) ([]byte, error)
}

func (j jsonMock) Unmarshal(data []byte, v interface{}) error {
	return j.UnmarshalStub(data, v)
}

func (j jsonMock) Marshal(v interface{}) ([]byte, error) {
	return j.MarshalStub(v)
}

// LogService
type mockLogService struct {
	log.LogService
	AnyStub func(key string, value interface{}) log.Field
}

func (m mockLogService) Any(key string, value interface{}) log.Field {
	return m.AnyStub(key, value)
}

// Data mocks

func mockConfig() *Config {
	baseURL, err := url.Parse("https://urltest.com")
	if err != nil {
		panic(err)
	}
	baseURLLegacy, err := url.Parse("https://urltest.legacy.com")
	if err != nil {
		panic(err)
	}
	configTransactionsApp := domain.ConfigTransactionsApp{
		BetaScope:     "beta",
		BaseURL:       *baseURL,
		BaseURLLegacy: *baseURLLegacy,
		ProvidersWithLegacyURL: map[string]struct{}{
			"conductor": {},
		},
	}
	return NewConfig(configTransactionsApp)
}

func mockCallReversalSuccess(provider string) func() func() {
	urlExpected := "https://urltest.legacy.com/cards/transactions/debit/reversals"
	urlPost := urlExpected

	if provider != "conductor" {
		urlExpected = "https://urltest.com/cards/transactions/debit/reversals"
		urlPost = urlExpected
	}

	return func() func() {
		return apitest.NewMock().
			Post(urlPost).
			Header("Content-Type", "application/json").
			Body(string(_reversalRequest)).
			AddMatcher(func(r *http.Request, mr *apitest.MockRequest) error {
				if urlExpected != r.URL.String() {
					return fmt.Errorf("received url %s did not match url expected: %s", r.URL.String(), urlExpected)
				}
				return nil
			}).
			RespondWith().
			Body(string(_reversalResponse)).
			Status(http.StatusOK).
			EndStandalone()
	}
}

func mockCallReversalBad() func() func() {
	return func() func() {
		return apitest.NewMock().
			Post("/reversals").
			Header("Content-Type", "application/json").
			Body(string(_reversalRequest)).
			RespondWith().
			Body(string(_reversalResponseBadRequest)).
			Status(http.StatusBadRequest).
			EndStandalone()
	}
}

func mockCallReversalFail() func() func() {
	return func() func() {
		return apitest.NewMock().
			Post("/noexist").
			Header("Content-Type", "application/json").
			RespondWith().
			EndStandalone()
	}
}

func newMockReversalInput(t *testing.T) domain.ReversalInput {
	var input domain.ReversalInput
	errExpected := json.Unmarshal(_transactionData, &input)
	input.HeaderValueXClientID = "123456"
	assert.NoError(t, errExpected)
	return input
}
