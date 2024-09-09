package topicvalidation

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/publisher/v1"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/validation-result_domain.json
var _validationResultDomain []byte

//go:embed testdata/bi-validation-expected_message.json
var _biMsgsTopic []byte

//go:embed testdata/ds-validation-expected_message.json
var _dsMsgsTopic []byte

//go:embed testdata/validation-result-eligible_domain.json
var _validationResultEligibleDomain []byte

//go:embed testdata/validation-result-eligible-expected_message.json
var _msgEligibleTopic []byte

// Publisher.
type mockPublisher struct {
	publisher.Publisher
	PublishStub func(ctx context.Context, body interface{}, filters ...string) error
}

func (m mockPublisher) Publish(
	ctx context.Context,
	body interface{},
	filters ...string,
) error {
	return m.PublishStub(ctx, body, filters...)
}

func mockValidationResultNotEligible(t *testing.T) (validationResult storage.ValidationResult) {
	err := json.Unmarshal(_validationResultDomain, &validationResult)
	assert.NoError(t, err)
	return
}

func mockExpectedMessagesTopic(t *testing.T, msgsTopic []byte) (msgs map[string]validationResultTopic) {
	err := json.Unmarshal(msgsTopic, &msgs)
	assert.NoError(t, err)
	return
}

func mockValidationResultEligible(t *testing.T) (validationResult storage.ValidationResult) {
	err := json.Unmarshal(_validationResultEligibleDomain, &validationResult)
	assert.NoError(t, err)
	return
}

func mockExpectedEligibleMessageTopic(t *testing.T) (msg validationResultTopic) {
	err := json.Unmarshal(_msgEligibleTopic, &msg)
	assert.NoError(t, err)
	return
}
