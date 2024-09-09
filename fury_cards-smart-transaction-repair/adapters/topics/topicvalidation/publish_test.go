package topicvalidation

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2/field"
	"github.com/melisource/fury_cards-go-toolkit/pkg/utilstest"
	"github.com/stretchr/testify/assert"
)

func Test_setValueDeclined(t *testing.T) {
	validationResult := &validationResultTopic{}
	qtyRepairs := 1
	validationResult.ValueDeclined.toString(qtyRepairs)
	assert.Equal(t, "1", string(validationResult.ValueDeclined))

	validationResult = &validationResultTopic{}
	blocked := true
	validationResult.ValueDeclined.toString(blocked)
	assert.Equal(t, "true", string(validationResult.ValueDeclined))

	validationResult = &validationResultTopic{}
	amount := 2500.59
	validationResult.ValueDeclined.toString(amount)
	assert.Equal(t, "2500.59", string(validationResult.ValueDeclined))

	validationResult = &validationResultTopic{}
	statusDetail := "cancel"
	validationResult.ValueDeclined.toString(statusDetail)
	assert.Equal(t, "cancel", string(validationResult.ValueDeclined))

	validationResult = &validationResultTopic{}
	restrictsFailed := []string{"xxx2", "xxx3"}
	validationResult.ValueDeclined.toString(restrictsFailed)
	assert.Equal(t, "[xxx2,xxx3]", string(validationResult.ValueDeclined))

	validationResult = &validationResultTopic{}
	dateCreated := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	validationResult.ValueDeclined.toString(dateCreated)
	assert.Equal(t, "2023-09-01T00:00:00Z", string(validationResult.ValueDeclined))
}

func Test_validationEvent_Publish_NotEligible_BI(t *testing.T) {
	ctx := context.TODO()
	wrapFields := field.NewWrappedFields()
	expectedMessagesTopic := mockExpectedMessagesTopic(t, _biMsgsTopic)

	m := mockPublisher{
		PublishStub: func(ctx context.Context, body interface{}, filters ...string) error {

			validationResult := body.(validationResultTopic)
			expectedValidation := expectedMessagesTopic[validationResult.Rule]

			// forces equality of identifiers for testing purposes.
			validationResult.ID = expectedValidation.ID
			validationResult.Reason = expectedValidation.Reason
			validationResult.ApprovedData = expectedValidation.ApprovedData

			jsonPublishedMessage, _ := json.Marshal(validationResult)
			jsonExpectedMessage, _ := json.Marshal(expectedValidation)

			assert.JSONEq(t, string(jsonExpectedMessage), string(jsonPublishedMessage), validationResult.Rule)

			return nil
		},
	}

	eventValidation := New(m)

	input := mockValidationResultNotEligible(t)
	input.Reason["qty_reparation_per_period"] =
		&domain.ReasonResult{Actual: 100, Accepted: domain.ReasonResultPerPeriod{Qty: 100, PeriodDays: 1000}}

	err := eventValidation.Publish(ctx, input, wrapFields)
	assert.NoError(t, err)
}

func Test_validationEvent_Publish_NotEligible_DS(t *testing.T) {
	ctx := context.TODO()
	wrapFields := field.NewWrappedFields()
	expectedMessagesTopic := mockExpectedMessagesTopic(t, _dsMsgsTopic)

	currentMessage := 0
	m := mockPublisher{
		PublishStub: func(ctx context.Context, body interface{}, filters ...string) error {
			if currentMessage > 0 {
				return nil
			}

			validationResult := body.(validationResultTopic)
			expectedValidation := expectedMessagesTopic["full"]

			for rule, reason := range validationResult.Reason {
				expectedReason := expectedValidation.Reason[rule]
				assert.Equal(t, reason, expectedReason)
			}

			validationResult.Reason = nil
			expectedValidation.Reason = nil

			// forces equality of identifiers for testing purposes.
			validationResult.ID = expectedValidation.ID
			validationResult.ValueDeclined = expectedValidation.ValueDeclined
			validationResult.AllowedValue = expectedValidation.AllowedValue
			validationResult.Rule = expectedValidation.Rule

			jsonPublishedMessage, _ := json.Marshal(validationResult)
			jsonExpectedMessage, _ := json.Marshal(expectedValidation)

			assert.JSONEq(t, string(jsonExpectedMessage), string(jsonPublishedMessage), validationResult.Rule)

			currentMessage++
			return nil
		},
	}

	eventValidation := New(m)

	input := mockValidationResultNotEligible(t)

	err := eventValidation.Publish(ctx, input, wrapFields)
	assert.NoError(t, err)
}

func Test_validationEvent_Publish_Eligible(t *testing.T) {
	ctx := context.TODO()
	wrapFields := field.NewWrappedFields()
	expectedEligibleMessageTopic := mockExpectedEligibleMessageTopic(t)

	m := mockPublisher{
		PublishStub: func(ctx context.Context, body interface{}, filters ...string) error {

			validationResult := body.(validationResultTopic)
			expectedValidation := expectedEligibleMessageTopic

			// forces equality of identifiers for testing purposes.
			validationResult.ID = expectedValidation.ID

			jsonPublishedMessage, _ := json.Marshal(validationResult)
			jsonExpectedMessage, _ := json.Marshal(expectedValidation)

			assert.JSONEq(t, string(jsonExpectedMessage), string(jsonPublishedMessage), validationResult.Rule)

			return nil
		},
	}

	eventValidation := New(m)

	input := mockValidationResultEligible(t)

	err := eventValidation.Publish(ctx, input, wrapFields)
	assert.NoError(t, err)
}

func Test_validationEvent_Publish_NotEligible_With_Error(t *testing.T) {
	ctx := context.TODO()
	wrapFields := field.NewWrappedFields()
	expectedMessagesTopic := mockExpectedMessagesTopic(t, _biMsgsTopic)

	msgsWithErrors := map[string]validationResultTopic{}
	currentMessage := 0
	m := mockPublisher{
		PublishStub: func(ctx context.Context, body interface{}, filters ...string) error {
			validationResult := body.(validationResultTopic)
			expectedValidation := expectedMessagesTopic[validationResult.Rule]

			if currentMessage == 0 || currentMessage == 2 {
				currentMessage++
				msgsWithErrors["message_body_from_id_"+validationResult.ID] = validationResult
				return errors.New("publication error")
			}

			// forces equality of identifiers for testing purposes.
			validationResult.ID = expectedValidation.ID
			validationResult.Reason = expectedValidation.Reason
			validationResult.ApprovedData = expectedValidation.ApprovedData

			jsonPublishedMessage, _ := json.Marshal(validationResult)
			jsonExpectedMessage, _ := json.Marshal(expectedValidation)

			assert.JSONEq(t, string(jsonExpectedMessage), string(jsonPublishedMessage), validationResult.Rule)

			currentMessage++

			return nil
		},
	}

	eventValidation := New(m)

	input := mockValidationResultNotEligible(t)
	input.Reason["qty_reparation_per_period"] =
		&domain.ReasonResult{Actual: 100, Accepted: domain.ReasonResultPerPeriod{Qty: 100, PeriodDays: 1000}}

	err := eventValidation.Publish(ctx, input, wrapFields)

	// Make sure you tried to send all messages.
	assert.Equal(t, currentMessage, len(expectedMessagesTopic))

	expectedErr := domain.BuildErrorBadGateway(errors.New("publication error"), "validation_result_event_publication_failure")
	utilstest.AssertGnierrorsExpected(t, expectedErr, err)

	// Ensures messages with errors will be logged.
	fields := wrapFields.Fields.ToLogField(log.InfoLevel)
	fieldsCast := fields.Interface.(map[string]interface{})
	for id, msg := range msgsWithErrors {
		validationWithError, found := fieldsCast[id]
		assert.True(t, found)

		msgValidationErr, _ := json.Marshal(msg)
		assert.Equal(t, string(msgValidationErr), validationWithError)
	}
}

func Test_validationEvent_Publish_Eligible_With_Error(t *testing.T) {
	ctx := context.TODO()
	wrapFields := field.NewWrappedFields()

	var validationResultFailed validationResultTopic
	m := mockPublisher{
		PublishStub: func(ctx context.Context, body interface{}, filters ...string) error {
			validationResultFailed = body.(validationResultTopic)
			return errors.New("publication error")
		},
	}

	eventValidation := New(m)

	input := mockValidationResultEligible(t)

	err := eventValidation.Publish(ctx, input, wrapFields)

	expectedErr := domain.BuildErrorBadGateway(errors.New("publication error"), "validation_result_event_publication_failure")

	utilstest.AssertGnierrorsExpected(t, expectedErr, err)

	// Ensures message with errors will be logged.
	fields := wrapFields.Fields.ToLogField(log.InfoLevel)
	fieldsCast := fields.Interface.(map[string]interface{})
	expectedValidationWithErr, _ := json.Marshal(validationResultFailed)
	assert.Equal(t, string(expectedValidationWithErr), fieldsCast["message_body_from_id_"+validationResultFailed.ID])
}
