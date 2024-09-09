package domain

import "testing"

func TestConfigRulesSite_CheckApplyAnyScore(t *testing.T) {
	tests := []struct {
		name      string
		want      bool
		rulesSite *ConfigRulesSite
	}{
		{
			name:      "not apply - 1",
			want:      false,
			rulesSite: &ConfigRulesSite{},
		},
		{
			name: "apply - 1",
			want: true,
			rulesSite: &ConfigRulesSite{
				QtyReparationPerPeriodDays: &ConfigQtyReparationPeriod{
					Score: &ConfigScoreQtyReparationPeriod{},
				},
			},
		},
		{
			name: "apply - 2",
			want: true,
			rulesSite: &ConfigRulesSite{
				MaxAmountReparation: &ConfigMaxAmountReparation{
					Score: &ConfigScoreMaxAmountReparation{},
				},
			},
		},
		{
			name: "apply - 3",
			want: true,
			rulesSite: &ConfigRulesSite{
				UserLifetime: &ConfigUserLifetime{
					Score: &ConfigScoreUserLifetime{},
				},
			},
		},
		{
			name: "apply - 4",
			want: true,
			rulesSite: &ConfigRulesSite{
				Score: &ConfigScore{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rulesSite.CheckApplyAnyScore(); got != tt.want {
				t.Errorf("CheckApplyAnyScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
