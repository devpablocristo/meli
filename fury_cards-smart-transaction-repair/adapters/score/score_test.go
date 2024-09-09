package score

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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

func Test_Score_Success(t *testing.T) {
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{
		ReadAllStub: io.ReadAll,
	}
	mockNewJson := pkgjson.NewJSON()
	mockNewLog := mockLogService{}

	mockScoreService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := newMockScoreInput(t)
	mockCallScoreSuccess()()

	expectedOutput := &domain.ScoreOutput{
		Score: 1.2322323,
	}

	output, err := mockScoreService.GetScore(context.TODO(), input)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func Test_Score_Error_MarshalRequest(t *testing.T) {
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
				Key:    "body_score",
				String: value.(string),
			}
		},
	}

	mockScoreService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := newMockScoreInput(t)

	output, err := mockScoreService.GetScore(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.InternalServerError,
		errors.New("error marshal"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Empty(t, output)
}

func Test_Score_Error_Call_Post(t *testing.T) {
	mockNewJson := pkgjson.NewJSON()
	mockNewHttpClient := rest.NewHTTPWithTimeout(5 * time.Second)
	mockNewIoutil := ioutilMock{}
	mockNewLog := mockLogService{
		AnyStub: func(key string, value interface{}) log.Field {
			assert.NotEmpty(t, value)
			return log.Field{
				Key:    "body_score",
				String: value.(string),
			}
		},
	}

	mockScoreService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	mockCallScoreFail()()

	input := newMockScoreInput(t)

	output, err := mockScoreService.GetScore(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New("error read"), "request_score_failed", false)

	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
	assert.Empty(t, output)
}

func Test_Score_Error_ValidateResponse(t *testing.T) {
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
					Key:    "body_score",
					String: value.(string),
				}
			},
		}

		mockScoreService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

		input := newMockScoreInput(t)

		mockCallScoreSuccess()()

		output, err := mockScoreService.GetScore(context.TODO(), input)

		errorExpected := gnierrors.Wrap(domain.BadGateway,
			errors.New("error read"), "request_score_failed", false)

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
					Key:    "body_score",
					String: value.(string),
				}
			},
		}

		mockScoreService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

		mockCallScoreBad()()

		input := newMockScoreInput(t)

		output, err := mockScoreService.GetScore(context.TODO(), input)

		errorExpected := gnierrors.Wrap(domain.BadRequest,
			bodyResponseError(_scoreResponseBadRequest), "request_score_failed", false)

		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
		assert.Empty(t, output)
	})
}

func Test_Score_Error_Unmarshal(t *testing.T) {
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
				Key:    "body_score",
				String: value.(string),
			}
		},
	}

	mockScoreService := New(mockConfig(), mockNewHttpClient, mockNewIoutil, mockNewJson, mockNewLog)

	input := newMockScoreInput(t)

	mockCallScoreSuccess()()

	output, err := mockScoreService.GetScore(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.InternalServerError,
		errors.New("error unmarshal"), "request_score_failed", false)

	utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	assert.Empty(t, output)
}

func Test_score_buildScoreOutput(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockScoreService := score{}

		output, err := mockScoreService.buildScoreOutput(responseScore{
			Result: "1",
			Score:  "123",
		})

		assert.NoError(t, err)
		assert.Equal(t, 123.0, output.Score)
	})

	t.Run("error empty score", func(t *testing.T) {
		mockScoreService := score{}

		output, err := mockScoreService.buildScoreOutput(responseScore{
			Result: "1",
			Score:  "",
		})

		assert.Nil(t, output)

		errorExpected := gnierrors.Wrap(domain.BadGateway,
			errors.New(`empty score`), "request_score_failed", false)

		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	})

	t.Run("error parse score", func(t *testing.T) {
		mockScoreService := score{}

		output, err := mockScoreService.buildScoreOutput(responseScore{
			Result: "1",
			Score:  "-",
		})

		assert.Nil(t, output)

		errorExpected := gnierrors.Wrap(domain.BadGateway,
			errors.New(`strconv.ParseFloat: parsing "-": invalid syntax`), "request_score_failed", false)

		utilstest.AssertGnierrorsExpected(t, errorExpected, err)
	})
}
