package errorcase

import (
	"testing"
)

// errorAny is an Error implementation used when test cases expect a non-specific error
type errorAny struct {
}

// NewErrorAny returns an errorAny testcase
func NewErrorAny() Error {
	return errorAny{}
}

// EvalError checks whether err value is expected
func (ea errorAny) EvalError(t *testing.T, err error) {

	if err == nil {
		t.Fatal("Expected any error, but none occurred")
	}
}
