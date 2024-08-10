package common

import (
	"fmt"
)

type ErrorMap map[string]interface{}

type HttpJsonError struct {
	StatusCode int
	Message    string
	ErrorData  ErrorMap
}

func (e *HttpJsonError) Error() string {
	return fmt.Sprintf("StatusCode:%d, Message:%s", e.StatusCode, e.Message)
}
