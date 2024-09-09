package graphqlkyc

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
	"github.com/stretchr/testify/assert"
)

func Test_searchKyc_GetUserKyc_Success(t *testing.T) {
	mockService := mockSearchKycClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			v := dataResponse.(*response)
			*v = response{
				User: user{Identification: identification{ID: "abc"}, DateCreated: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			}
			return nil
		},
	}

	service := New(mockService)

	input := domain.SearchKycInput{
		UserID: "123",
	}

	output, err := service.GetUserKyc(context.TODO(), input)

	assert.NoError(t, err)
	utilstest.AssertFieldsNoEmptyFromStruct(t, *output)
	assert.Equal(t, "abc", output.KycIdentificationID)
}

func Test_searchKyc_GetUserKyc_Error(t *testing.T) {
	mockService := mockSearchKycClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			return errors.New(`{
				"errors": [
				  {
					"message": "timeout",
					"path": [
					  "getUser"
					],
					"extensions": {
					  "code": "451",
					  "message": ""
					}
				  }
				],
				"status_code": 200
			  }`)
		},
	}

	service := New(mockService)

	input := domain.SearchKycInput{
		UserID: "123",
	}

	output, err := service.GetUserKyc(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.BadGateway,
		errors.New(`{
			"errors": [
			  {
				"message": "timeout",
				"path": [
				  "getUser"
				],
				"extensions": {
				  "code": "451",
				  "message": ""
				}
			  }
			],
			"status_code": 200
		  }`), "graphql-search-kyc get_user_kyc failure", false)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.JSONEq(t, gnierrors.Cause(err).Error(), gnierrors.Cause(errorExpected).Error())
	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
}

func Test_searchKyc_GetUserKyc_Error_UserNotAvailable(t *testing.T) {
	mockService := mockSearchKycClient{
		SetDataByQueryRequestStub: func(ctx context.Context, request graphql.Request, dataResponse interface{}) error {
			bErr := []byte(`{
				"errors": [
				  {
					"message": "The requested user is not available due to legal reasons",
					"path": [
					  "getUser"
					],
					"extensions": {
					  "code": "451",
					  "message": ""
					}
				  }
				],
				"status_code": 200
			  }`)

			var err graphql.Error
			assert.NoError(t, json.Unmarshal(bErr, &err))
			return err
		},
	}

	service := New(mockService)

	input := domain.SearchKycInput{
		UserID: "123",
	}

	output, err := service.GetUserKyc(context.TODO(), input)

	errorExpected := gnierrors.Wrap(domain.UnavailableForLegalReasons,
		errors.New(`{
			"errors": [
			  {
				"message": "The requested user is not available due to legal reasons",
				"path": [
				  "getUser"
				],
				"extensions": {
				  "code": "451",
				  "message": ""
				}
			  }
			],
			"status_code": 200
		  }`), "graphql-search-kyc get_user_kyc failure", false)

	assert.Nil(t, output)
	assert.JSONEq(t, gnierrors.Cause(err).Error(), gnierrors.Cause(errorExpected).Error())
	assert.Equal(t, gnierrors.Type(errorExpected), gnierrors.Type(err))
	assert.Equal(t, gnierrors.Message(errorExpected), gnierrors.Message(err))
}
