package serviceaccount

import (
	"os"

	"github.com/pkg/errors"

	"github.com/athera-io/gcs-bucket-helper/internal/tilde"
)

// PathExists ensures that path already exists
func PathExists(s string) error {

	s, err := tilde.ExpandPath(s)

	if err != nil {
		return err
	}

	stat, err := os.Stat(s)

	if os.IsNotExist(err) {
		return newNotExistError(s)
	}

	if !stat.IsDir() {
		return newNotDirError(s)
	}

	return nil
}

// NotExistError is useful for test logs
const NotExistError = "NotExist"

// notExistError gives more user friendly output than equivalent error from package os
type notExistError struct {
	path string
}

// newNotExistError ...
func newNotExistError(path string) error {
	return notExistError{
		path: path,
	}
}

func (err notExistError) Error() string {
	return "No such directory: " + err.path
}

// isNotExist method is called by IsNotExist function
func (err notExistError) isNotExistError() bool {
	return true
}

// IsNotExistError returns true if err is tilde
func IsNotExistError(err error) bool {
	e, ok := errors.Cause(err).(notExistError)
	return ok && e.isNotExistError()
}

// NotDirError is useful for test logs
const NotDirError = "NotDir"

type notDirError struct {
	path string
}

// newNotDirError ...
func newNotDirError(path string) error {
	return notDirError{
		path: path,
	}
}

func (err notDirError) Error() string {
	return "Path exists, but is not a directory: " + err.path
}

// isNotDir method is called by IsNotDir function
func (err notDirError) isNotDirError() bool {
	return true
}

// IsNotDirError returns true if err is tilde
func IsNotDirError(err error) bool {
	e, ok := errors.Cause(err).(notDirError)
	return ok && e.isNotDirError()
}
