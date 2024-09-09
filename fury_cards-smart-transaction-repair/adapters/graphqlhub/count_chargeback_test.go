package graphqlhub

import (
	"context"
	"errors"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql/searchhubclient/v1"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_searchHub_CountChargebackFromUser_Success(t *testing.T) {
	mockService := mockSearchHubClient{
		CountChargebacksInTheLastDaysStub: func(ctx context.Context, userID string, lastDays int) (*searchhubclient.OutputCountChargeback, error) {
			return &searchhubclient.OutputCountChargeback{Total: 10}, nil
		},
	}

	service := New(mockService)

	input := domain.SearchHubCountChargebackInput{
		UserID:   "123",
		LastDays: 30,
	}

	output, err := service.CountChargebackFromUser(context.TODO(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, output.Total, 10)
}

func Test_searchHub_CountChargebackFromUser_Error(t *testing.T) {
	mockService := mockSearchHubClient{
		CountChargebacksInTheLastDaysStub: func(ctx context.Context, userID string, lastDays int) (*searchhubclient.OutputCountChargeback, error) {
			return nil, errors.New(`{
				"error_response":
					[{
						"message":"error [not_found]: Not found wallet with user_id: 1117823961",
						"path":["wallet"],
						"extensions":{"code":"not_found","message":"Not found wallet with user_id: 1117823961"}
					}],
				"status_code":200}`)
		},
	}

	service := New(mockService)

	input := domain.SearchHubCountChargebackInput{
		UserID:   "123",
		LastDays: 30,
	}

	output, err := service.CountChargebackFromUser(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`{
			"error_response":
				[{
					"message":"error [not_found]: Not found wallet with user_id: 1117823961",
					"path":["wallet"],
					"extensions":{"code":"not_found","message":"Not found wallet with user_id: 1117823961"}
				}],
			"status_code":200}`), "search-hub failure", false)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.JSONEq(t, gnierrors.Cause(err).Error(), gnierrors.Cause(errorExpected).Error())
}
