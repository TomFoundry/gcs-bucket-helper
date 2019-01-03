/*Package errorcase has types that can be embedded in (table-driven) test cases
 */
package errorcase

import (
	"testing"
)

// Error is embedded in test cases that require error handling and evaluation
type Error interface {
	EvalError(t *testing.T, err error)
}

// Eval ...
func Eval(t *testing.T, err error, errorCase Error) (complete bool) {
	if errorCase == nil {

		if err != nil {
			t.Fatalf("Unexpected error (%v)", err)
		}

		return // Test is not complete because calling function will now evaluate test output
	}

	errorCase.EvalError(t, err)

	return true // Test is complete because we have evaluated the error case
}
