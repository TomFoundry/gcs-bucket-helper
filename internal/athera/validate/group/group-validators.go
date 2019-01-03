package group

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/athera-io/gcs-bucket-helper/internal/athera/models"
)

// LegalIndex ensures that input is an integer greater than zero, and less than 1 + length of groups
func LegalIndex(s string, groups []*models.Group, lineage models.GroupLineage) error {

	i, err := strconv.Atoi(s)

	if err != nil {
		return newNotIntegerError()
	}

	if i < 1 {
		return newInvalidIndexError("Input must be greater than zero")
	}

	max := len(groups)

	if len(lineage) > 0 {
		max++
	}

	if i > max {
		return newInvalidIndexErrorF("Input must be less than or equal to %d", max)
	}

	return nil
}

// NotIntegerError is useful for test logs
const NotIntegerError = "NotInteger"

// notIntegerError gives more user friendly output than equivalent error from package os
type notIntegerError struct {
}

// newNotIntegerError ...
func newNotIntegerError() error {
	return notIntegerError{}
}

func (err notIntegerError) Error() string {
	return "Input must be an integer"
}

// isNotInteger method is called by IsNotInteger function
func (err notIntegerError) isNotIntegerError() bool {
	return true
}

// IsNotIntegerError returns true if err is tilde
func IsNotIntegerError(err error) bool {
	e, ok := errors.Cause(err).(notIntegerError)
	return ok && e.isNotIntegerError()
}

// InvalidIndexError is useful for test logs
const InvalidIndexError = "InvalidIndex"

// invalidIndexError gives more user friendly output than equivalent error from package os
type invalidIndexError struct {
	msg string
}

// newInvalidIndexError ...
func newInvalidIndexError(msg string) error {
	return invalidIndexError{
		msg: msg,
	}
}

// newInvalidIndexErrorF ...
func newInvalidIndexErrorF(msg string, args ...interface{}) error {
	return invalidIndexError{
		msg: fmt.Sprintf(msg, args),
	}
}

func (err invalidIndexError) Error() string {
	return err.msg
}

// isInvalidIndex method is called by IsInvalidIndex function
func (err invalidIndexError) isInvalidIndexError() bool {
	return true
}

// IsInvalidIndexError returns true if err is tilde
func IsInvalidIndexError(err error) bool {
	e, ok := errors.Cause(err).(invalidIndexError)
	return ok && e.isInvalidIndexError()
}
