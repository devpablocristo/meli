package graphqltransaction

import (
	"testing"
	"time"

	rest "github.com/melisource/fury_cards-go-toolkit/pkg/http/v1"
	ioutilabs "github.com/melisource/fury_cards-go-toolkit/pkg/ioutil/v1"
	json "github.com/melisource/fury_cards-go-toolkit/pkg/json/v1"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	transactionSearch := New(NewConfig(mockBaseURL()), rest.NewHTTPWithTimeout(1*time.Second), ioutilabs.New(), json.NewJSON())
	assert.NotEmpty(t, transactionSearch)
}
