package errors

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

type CustomError struct {
	Status  int
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Cause   interface{} `json:"cause"`
}

func (c *CustomError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Cause   interface{} `json:"cause"`
	}{
		Code:    c.Code,
		Message: c.Message,
		Cause:   c.Cause,
	})
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("%s: %s: %s:", c.Code, c.Message, c.Cause)
}

func (c *CustomError) StatusCode() int {
	return c.Status
}

type causeCustom struct {
	Reason           string        `json:"reason"`
	CreationDatetime time.Time     `json:"creation_datetime"`
	ReasonDetail     domain.Reason `json:"reason_detail"`
}
