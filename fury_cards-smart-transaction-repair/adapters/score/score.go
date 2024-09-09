package score

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
)

const (
	headerValueContentType = "application/json"
	pathURL                = "/%s/capture/predict"
	bitSize                = 64
)

func (s score) GetScore(ctx context.Context, input domain.ScoreInput) (*domain.ScoreOutput, error) {
	request := fillRequest(input.ScoreData)
	bodyRequest, err := s.json.Marshal(request)
	if err != nil {
		return nil, s.buildError(err, statusInternalServerError)
	}
	baseURL := s.getBaseURL()
	url := rest.CreateURL(baseURL, fmt.Sprintf(pathURL, strings.ToLower(input.ScoreData.SiteID)))
	response, err := s.httpClient.Do(ctx, "POST", url, headerValueContentType, bodyRequest, map[string][]string{})
	if err != nil {
		return nil, s.buildError(err, statusBadGateway)
	}
	body, err := s.readAndValidateResponse(response)
	if err != nil {
		return nil, err
	}
	var responseScore responseScore
	err = s.json.Unmarshal(body, &responseScore)
	if err != nil {
		return nil, s.buildError(err, statusInternalServerError)
	}

	return s.buildScoreOutput(responseScore)
}

func (s score) readAndValidateResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	body, err := s.ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, s.buildError(err, statusBadGateway)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, s.buildError(bodyResponseError(body), statusCodeScore(response.StatusCode))
	}

	return body, nil
}

func fillRequest(scoreData domain.TransactionData) requestScore {
	return requestScore{
		Environment:        scoreData.Environment,
		SiteID:             scoreData.SiteID,
		StatusDetail:       scoreData.StatusDetail,
		PayerID:            scoreData.PayerID,
		Type:               scoreData.Type,
		IsRecurringPayment: scoreData.IsRecurringPayment,
		Operation: operation{
			SubType:                scoreData.Operation.SubType,
			CreationDatetime:       scoreData.Operation.CreationDatetime,
			TransmissionDatetime:   scoreData.Operation.TransmissionDatetime,
			IsAdvice:               scoreData.Operation.IsAdvice,
			Installments:           scoreData.Operation.Installments,
			AcquirerCode:           scoreData.Operation.AcquirerCode,
			Stan:                   scoreData.Operation.Stan,
			PaymentAge:             scoreData.Operation.PaymentAge,
			AccountType:            scoreData.Operation.AccountType,
			AuthorizationIndicator: scoreData.Operation.AuthorizationIndicator,
			IsInternational:        scoreData.Operation.IsInternational,
			Authorization: authorization{
				ID:                   scoreData.Operation.Authorization.ID,
				Stan:                 scoreData.Operation.Authorization.Stan,
				TransmissionDatetime: scoreData.Operation.Authorization.TransmissionDatetime,
				CardAcceptor: cardAcceptor{
					Terminal:      scoreData.Operation.Authorization.CardAcceptor.Terminal,
					Name:          scoreData.Operation.Authorization.CardAcceptor.Name,
					EntryMode:     scoreData.Operation.Authorization.CardAcceptor.EntryMode,
					Type:          scoreData.Operation.Authorization.CardAcceptor.Type,
					TerminalMode:  scoreData.Operation.Authorization.CardAcceptor.TerminalMode,
					PinCapability: scoreData.Operation.Authorization.CardAcceptor.PinCapability,
				},
			},
			Transaction: transaction{
				Amount:        scoreData.Operation.Transaction.Amount,
				Currency:      scoreData.Operation.Transaction.Currency,
				DecimalDigits: scoreData.Operation.Transaction.DecimalDigits,
				TotalAmount:   scoreData.Operation.Transaction.TotalAmount,
			},
			Card: card{
				NumberId: scoreData.Operation.Card.NumberID,
				Country:  scoreData.Operation.Card.Country,
				IssuerAccount: issuerAccount{
					BusinessMode: scoreData.Operation.Card.IssuerAccount.CardBusinessMode,
				},
				TokenInfo: tokenInfo{
					TokenWalletID: scoreData.Operation.Card.TokenInfo.TokenWalletID,
				},
			},
			Billing: billing{
				Amount:        scoreData.Operation.Billing.Amount,
				Currency:      scoreData.Operation.Billing.Currency,
				DecimalDigits: scoreData.Operation.Billing.DecimalDigits,
				TotalAmount:   scoreData.Operation.Billing.TotalAmount,
				Conversion: conversion{
					Date:          scoreData.Operation.Billing.Conversion.Date,
					DecimalDigits: scoreData.Operation.Billing.Conversion.DecimalDigits,
					Rate:          scoreData.Operation.Billing.Conversion.Rate,
					From:          scoreData.Operation.Billing.Conversion.From,
				},
			},
			Settlement: settlement{
				Amount:        scoreData.Operation.Settlement.Amount,
				Currency:      scoreData.Operation.Settlement.Currency,
				DecimalDigits: scoreData.Operation.Settlement.DecimalDigits,
				TotalAmount:   scoreData.Operation.Settlement.TotalAmount,
				Conversion: conversion{
					Date:          scoreData.Operation.Settlement.Conversion.Date,
					DecimalDigits: scoreData.Operation.Settlement.Conversion.DecimalDigits,
					Rate:          scoreData.Operation.Settlement.Conversion.Rate,
					From:          scoreData.Operation.Settlement.Conversion.From,
				},
			},
		},
		Provider: provider{
			ID: scoreData.Provider.ID,
		},
	}
}

func (s score) getBaseURL() url.URL {
	return s.config.BaseURL
}

func (s score) buildScoreOutput(responseScore responseScore) (*domain.ScoreOutput, error) {
	if strings.TrimSpace(responseScore.Score) == "" {
		return nil, s.buildError(errors.New("empty score"), statusBadGateway)
	}

	score, err := strconv.ParseFloat(responseScore.Score, bitSize)
	if err != nil {
		return nil, s.buildError(err, statusBadGateway)
	}

	return &domain.ScoreOutput{
		Score: score,
	}, nil
}
