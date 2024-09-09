package ds

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
	"github.com/melisource/cards-smart-transaction-repair/core/storage"
	"github.com/melisource/fury_cards-go-toolkit/pkg/dsclient/v1"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/validation-result_domain.json
var _validationResultDomain []byte

//go:embed testdata/validation-result_domain_eligible.json
var _validationResultDomainEligible []byte

//go:embed testdata/result-eligible.json
var _resultEligible []byte

//go:embed testdata/result-not-eligible-full.json
var _resultNotEligibleFull []byte

//go:embed testdata/result-not-eligible-status-not-allowed.json
var _resultNotEligibleStatusNotAllowed []byte

//go:embed testdata/result-not-eligible-max-amount.json
var _resultNotEligibleMaxAmount []byte

//go:embed testdata/result-not-eligible-qty-repair-period.json
var _resultNotEligibleQtyRepairPeriod []byte

//go:embed testdata/result-not-eligible-blocked.json
var _resultNotEligibleBlocked []byte

//go:embed testdata/result-not-eligible-restrictions.json
var _resultNotEligibleRestrictions []byte

//go:embed testdata/result-not-eligible-qty-chargeback-period.json
var _resultNotEligibleQtyChargebackPeriod []byte

//go:embed testdata/result-not-eligible-merchant-blockedlist.json
var _resultNotEligibleRuleMerchantBlockedlist []byte

//go:embed testdata/result-not-eligible-subtype-not-allowed.json
var _resultNotEligibleSubTypeNotAllowed []byte

//go:embed testdata/result-not-eligible-user-lifetime.json
var _resultNotEligibleUserLifetime []byte

//go:embed testdata/result-not-eligible-min-total-historical-amount.json
var _resultNotEligibleMinTotalHistoricalAmount []byte

//go:embed testdata/result-not-eligible-min-transactions-qty.json
var _resultNotEligibleMinTransactionsQty []byte

//go:embed testdata/result-not-eligible-score.json
var _resultNotEligibleScore []byte

// dsClientService
type mockDsClientService struct {
	dsclient.Service
	SaveDocumentWithContextStub func(ctx context.Context, key string, value interface{}) error
}

func (m mockDsClientService) SaveDocumentWithContext(ctx context.Context, id string, document interface{}) error {
	return m.SaveDocumentWithContextStub(ctx, id, document)
}

func mockValidationResultNotEligibleFull(t *testing.T) (validationResult storage.ValidationResult) {
	err := json.Unmarshal(_validationResultDomain, &validationResult)
	assert.NoError(t, err)
	mockResult(t, validationResult.Reason)
	return
}

func mockValidationResultEligible(t *testing.T) (validationResult storage.ValidationResult) {
	err := json.Unmarshal(_validationResultDomainEligible, &validationResult)
	assert.NoError(t, err)
	mockResult(t, validationResult.ApprovedData)
	return
}

func mockResult(t *testing.T, approvedData map[domain.RuleType]*domain.ReasonResult) {
	for k, v := range approvedData {
		switch accepted := v.Accepted.(type) {
		case map[string]interface{}:
			b, err := json.Marshal(accepted)
			assert.NoError(t, err)

			var perPeriod domain.ReasonResultPerPeriod
			err = json.Unmarshal(b, &perPeriod)
			assert.NoError(t, err)

			approvedData[k] = &domain.ReasonResult{
				Actual: v.Actual,
				Accepted: domain.ReasonResultPerPeriod{
					Qty:        perPeriod.Qty,
					PeriodDays: perPeriod.PeriodDays,
				},
			}
		default:
			continue
		}
	}
}
