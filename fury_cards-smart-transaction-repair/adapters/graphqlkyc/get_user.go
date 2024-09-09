package graphqlkyc

import (
	"context"
	"net/http"
	"strconv"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/sdk/graphql"
)

func (s searchKyc) GetUserKyc(
	ctx context.Context,
	input domain.SearchKycInput,
) (*domain.SearchKycOutput, error) {

	request := graphql.Request{
		Query: "query($userId: ID!) { getUser(id: $userId) { identification { id } date_created } }",
		Params: map[string]interface{}{
			"userId": input.UserID,
		},
	}

	var data response

	err := s.service.SetDataByQueryRequest(ctx, request, &data)
	if err != nil {
		if errCustom, castOk := err.(graphql.Error); castOk {
			statusUnavailableForLegalReasons := strconv.Itoa(http.StatusUnavailableForLegalReasons)
			for _, e := range errCustom.ErrorResponse {
				if e.Extensions != nil && e.Extensions.Code == statusUnavailableForLegalReasons {
					return nil, buildErrorUnavailableForLegalReasons(err, msgErrFailGetUser)
				}
			}
		}
		return nil, buildErrorBadGateway(err, msgErrFailGetUser)
	}

	return &domain.SearchKycOutput{
		KycIdentificationID: data.User.Identification.ID,
		DateCreated:         data.User.DateCreated,
	}, nil
}
