package graphqltransaction

import (
	_ "embed"
	"io"
	"net/url"
	"time"

	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	jsonpkg "github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
)

//go:embed testdata/response-before-reversal.json
var _transactionResponse []byte

//go:embed testdata/request.json
var _transactionRequest []byte

//go:embed testdata/transaction-data.json
var _transactionData []byte

//go:embed testdata/response-err.json
var _responseErr []byte

type ioutilMock struct {
	ReadAllStub func(r io.Reader) ([]byte, error)
}

func (i ioutilMock) ReadAll(r io.Reader) ([]byte, error) {
	return i.ReadAllStub(r)
}

type jsonMock struct {
	UnmarshalStub func(data []byte, v interface{}) error
	MarshalStub   func(v interface{}) ([]byte, error)
}

func (j jsonMock) Marshal(v interface{}) ([]byte, error) {
	return j.MarshalStub(v)
}

func (j jsonMock) Unmarshal(data []byte, v interface{}) error {
	return j.UnmarshalStub(data, v)
}

func newMockTest() transactionSearch {
	return transactionSearch{
		config:     &Config{BaseURL: *mockBaseURL()},
		httpClient: rest.NewHTTPWithTimeout(10 * time.Second),
		ioutil: ioutilMock{
			ReadAllStub: io.ReadAll,
		},
		json: jsonpkg.NewJSON(),
	}
}

func mockBaseURL() *url.URL {
	baseURL, err := url.Parse("https://urltest.com")
	if err != nil {
		panic(err)
	}

	return baseURL
}
