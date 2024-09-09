package validation

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/mercadolibre/fury_gogenie/pkg/gnierrors/v2"
)

type loadedData struct {
	RulesSite                      domain.ConfigRulesSite
	Repairs                        []storage.ReparationOutput
	UserIsBlocked                  bool
	PolicyOutput                   *domain.PolicyOutput
	SearchHubCountChargebackOutput *domain.SearchHubCountChargebackOutput
	TransactionHistory             *domain.TransactionHistoryOutput
	CalculatedScore                *float64
}

func (v validationService) loadData(
	ctx context.Context,
	input domain.ValidationInput,
	wrapFields *field.WrappedFields,
) (*loadedData, error) {

	rulesSite, err := v.getParametersValidation(ctx, input, wrapFields)
	if err != nil {
		return nil, err
	}

	var errTransactionRepairs, errUserIsBlocked, errPolicy, errCountChargeback, errTransactionHistory, errGetScore error
	loadedData := &loadedData{
		RulesSite: *rulesSite,
	}

	wg := sync.WaitGroup{}

	if rulesSite.QtyReparationPerPeriodDays != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loadedData.Repairs, errTransactionRepairs = v.getTransactionRepairs(ctx, input, *rulesSite.QtyReparationPerPeriodDays, wrapFields)
		}()
	}

	if rulesSite.Blockeduser != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			input := domain.BlockedUserInput{KycIdentificationID: input.UserKycData.KycIdentificationID, SiteID: input.SiteID}
			loadedData.UserIsBlocked, errUserIsBlocked = v.blockedUserService.CheckUserIsBlocked(ctx, input, wrapFields)
		}()
	}

	if rulesSite.EvaluationUser != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loadedData.PolicyOutput, errPolicy = v.getPolicyAgentEvaluation(ctx, input, *rulesSite.EvaluationUser, wrapFields)
		}()
	}

	if rulesSite.QtyChargebackPerPeriodDays != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loadedData.SearchHubCountChargebackOutput, errCountChargeback =
				v.countChargebackFromUser(ctx, input, *rulesSite.QtyChargebackPerPeriodDays, wrapFields)
		}()
	}

	if rulesSite.MinTotalHistoricalAmount != nil || rulesSite.MinTransactionsQty != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loadedData.TransactionHistory, errTransactionHistory =
				v.getTransactionHistory(ctx, input, wrapFields)
		}()
	}

	if rulesSite.CheckApplyAnyScore() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loadedData.CalculatedScore, errGetScore =
				v.getScore(ctx, input, wrapFields)
		}()
	}

	wg.Wait()

	if errTransactionRepairs != nil {
		return nil, errTransactionRepairs
	}

	if errUserIsBlocked != nil {
		return nil, errUserIsBlocked
	}

	if errPolicy != nil {
		return nil, errPolicy
	}

	if errCountChargeback != nil {
		return nil, errCountChargeback
	}

	if errTransactionHistory != nil {
		return nil, errTransactionHistory
	}

	if errGetScore != nil {
		return nil, errGetScore
	}

	return loadedData, nil
}

func (v validationService) getParametersValidation(
	ctx context.Context,
	input domain.ValidationInput,
	wrapFields *field.WrappedFields,
) (*domain.ConfigRulesSite, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerConfigurationsLoadParametersValidation)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerLoadParametersValidation(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	parameters, err := v.configurations.LoadParametersValidation(ctx)
	if err != nil {
		addMetricErr(ctx, err, "get_parameters_validation", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, buildErrorBadGateway(err, gnierrors.Message(err))
	}

	configRulesSite, found := parameters.RulesSite[input.SiteID]
	if !found {
		err = fmt.Errorf("rules not found for site(%s) - parameters not configured", input.SiteID)
		addMetricErr(ctx, err, "map_rules_site", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, buildErrorBadGateway(err, msgErrParametersNotConfigured)
	}

	return configRulesSite, nil
}

func (v validationService) getTransactionRepairs(
	ctx context.Context,
	input domain.ValidationInput,
	configQtyReparationPeriod domain.ConfigQtyReparationPeriod,
	wrapFields *field.WrappedFields,
) ([]storage.ReparationOutput, error) {

	startTimer := wrapFields.Timers.Start(domain.TimerDsRepairGetListByUserIDAndCreationPeriod)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerGetListByUserIDAndCreationPeriod(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	dateNow := time.Now().UTC()
	qtyDays := oneDay * time.Duration(configQtyReparationPeriod.PeriodDays)
	gteCreatedAt := dateNow.Add(-qtyDays)

	repairs, err := v.reparationSearchRepo.GetListByUserIDAndCreationPeriod(input.UserID, gteCreatedAt, dateNow)
	if err != nil {
		addMetricErr(ctx, err, "get_transaction_repairs", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, buildErrorBadGateway(err, gnierrors.Message(err))
	}

	return repairs, nil
}

func (v validationService) getPolicyAgentEvaluation(
	ctx context.Context,
	input domain.ValidationInput,
	configEvaluationUser domain.ConfigEvaluationUser,
	wrapFields *field.WrappedFields,
) (output *domain.PolicyOutput, err error) {

	startTimer := wrapFields.Timers.Start(domain.TimerPolicyAgentEvaluate)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerPolicyAgentEvaluate(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	outputPolicy, err := v.policy.EvaluateWithUser(ctx, input.UserID, configEvaluationUser.Tags...)
	if err != nil {
		addMetricErr(ctx, err, "get_policy_agent_evaluation", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, err
	}

	return outputPolicy, nil
}

func (v validationService) countChargebackFromUser(
	ctx context.Context,
	input domain.ValidationInput,
	configQtyChargebackPeriod domain.ConfigQtyChargebackPeriod,
	wrapFields *field.WrappedFields,
) (output *domain.SearchHubCountChargebackOutput, err error) {

	startTimer := wrapFields.Timers.Start(domain.TimerSearchHubCountChargebackFromUser)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerSearchHubCountChargebackFromUser(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	searchHubCountChargebackInput := domain.SearchHubCountChargebackInput{
		UserID:   strconv.Itoa(int(input.UserID)),
		LastDays: configQtyChargebackPeriod.PeriodDays,
	}

	searchHubCountChargebackOutput, err := v.searchHub.CountChargebackFromUser(ctx, searchHubCountChargebackInput)
	if err != nil {
		addMetricErr(ctx, err, "count_chargeback_from_user", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, err
	}

	return searchHubCountChargebackOutput, nil
}

func (v validationService) getTransactionHistory(
	ctx context.Context,
	input domain.ValidationInput,
	wrapFields *field.WrappedFields,
) (output *domain.TransactionHistoryOutput, err error) {

	startTimer := wrapFields.Timers.Start(domain.TimerSearchHubGetTransactionHistory)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerSearchHubGetTransactionHistory(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	transactionHistoryInput := domain.TransactionHistoryInput{
		UserID: strconv.Itoa(int(input.UserID)),
	}

	transactionHistoryOutput, err := v.searchHub.GetTransactionHistory(ctx, transactionHistoryInput)
	if err != nil {
		addMetricErr(ctx, err, "get_transaction_history", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, err
	}

	return transactionHistoryOutput, nil
}

func (v validationService) getScore(
	ctx context.Context,
	input domain.ValidationInput,
	wrapFields *field.WrappedFields,
) (minScore *float64, err error) {

	startTimer := wrapFields.Timers.Start(domain.TimerGetScore)
	start := time.Now()
	defer func() {
		startTimer.Stop()
		addMetricTimerGetScore(ctx, input.SiteID, input.FaqID, input.Requester.ClientApplication, time.Since(start))
	}()

	outputScore, err := v.scoreService.GetScore(ctx, domain.ScoreInput{ScoreData: input.TransactionData})
	if err != nil {
		addMetricErr(ctx, err, "get_score", input.SiteID, input.FaqID, input.Requester.ClientApplication)
		return nil, err
	}

	calculatedScore := round(outputScore.Score, precisionRoundScore)
	return &calculatedScore, nil
}
