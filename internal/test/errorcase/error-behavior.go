package errorcase

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

type behaviorChecker func(error) bool

// errorBehavior is an Error implementation used when test cases expect an error to have a behavior
type errorBehavior struct {
	expectedTypeName string
	checker          behaviorChecker // If expectError is true and behaviorChecker is nil, then the test should succeed if any error is returned and fails if no error is returned
}

// NewErrorBehavior returns an errorBehavior testcase
func NewErrorBehavior(typeName string, checker behaviorChecker) Error {
	return errorBehavior{
		expectedTypeName: typeName,
		checker:          checker,
	}
}

// EvalError checks whether err value is expected
func (ec errorBehavior) EvalError(t *testing.T, err error) {

	cause := errors.Cause(err)

	if ok := ec.checker(cause); !ok {
		t.Fatalf("Expected %s error but got %v - (%v)", ec.expectedTypeName, reflect.TypeOf(cause), err)
	}
}
