package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_round(t *testing.T) {
	tests := []struct {
		name        string
		amounts     []float64
		amountRound float64
	}{
		{
			name:        "round 15.800999999999998 to 15.801",
			amounts:     []float64{10, 5, 0.1, 0.1, 0.1, 0.1005, 0.1, 0.1, 0.1005, 0.1},
			amountRound: 15.801,
		},
		{
			name:        "round 21.490000000000002 to 21.49",
			amounts:     []float64{10.501, 10.989},
			amountRound: 21.49,
		},
	}
	for _, tt := range tests {
		var amountTotal float64

		for _, v := range tt.amounts {
			amountTotal += v
		}

		amountTotal = round(amountTotal, 8)
		assert.Equal(t, tt.amountRound, amountTotal)
	}
}
