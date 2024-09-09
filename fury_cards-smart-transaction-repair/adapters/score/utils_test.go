package score

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

//go:embed testdata/score-data.json
var _scoreData []byte

//go:embed testdata/score-request.json
var _scoreRequest []byte

//go:embed testdata/score-response-badrequest.json
var _scoreResponseBadRequest []byte

//go:embed testdata/score-response.json
var _scoreResponse []byte

type mockLogService struct {
	log.LogService
	AnyStub func(key string, value interface{}) log.Field
}

func (m mockLogService) Any(key string, value interface{}) log.Field {
	return m.AnyStub(key, value)
}

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

type ioutilMock struct {
	ReadAllStub func(r io.Reader) ([]byte, error)
}

func (i ioutilMock) ReadAll(r io.Reader) ([]byte, error) {
	return i.ReadAllStub(r)
}

func newMockScoreInput(t *testing.T) domain.ScoreInput {
	var input domain.ScoreInput
	errExpected := json.Unmarshal(_scoreData, &input)

	assert.NoError(t, errExpected)
	return input
}

func mockConfig() *Config {
	baseURL, err := url.Parse("https://urltest.com")
	if err != nil {
		panic(err)
	}
	configScoreApp := domain.ConfigCardsDataModels{
		BaseURL: *baseURL,
	}
	return NewConfig(configScoreApp)
}

func mockCallScoreSuccess() func() func() {
	urlExpected := "https://urltest.com/mla/capture/predict"
	urlPost := urlExpected

	return func() func() {
		return apitest.NewMock().
			Post(urlPost).
			Header("Content-Type", "application/json").
			Body(string(_scoreRequest)).
			AddMatcher(func(r *http.Request, mr *apitest.MockRequest) error {
				if urlExpected != r.URL.String() {
					return fmt.Errorf("received url %s did not match url expected: %s", r.URL.String(), urlExpected)
				}
				return nil
			}).
			RespondWith().
			Body(string(_scoreResponse)).
			Status(http.StatusOK).
			EndStandalone()
	}
}

func mockCallScoreFail() func() func() {
	return func() func() {
		return apitest.NewMock().
			Post("/noexist").
			Header("Content-Type", "application/json").
			RespondWith().
			EndStandalone()
	}
}

func mockCallScoreBad() func() func() {
	return func() func() {
		return apitest.NewMock().
			Post("/mla/capture/predict").
			Header("Content-Type", "application/json").
			Body(string(_scoreRequest)).
			RespondWith().
			Body(string(_scoreResponseBadRequest)).
			Status(http.StatusBadRequest).
			EndStandalone()
	}
}
