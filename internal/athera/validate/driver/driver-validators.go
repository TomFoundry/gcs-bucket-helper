package driver

import (
	"errors"
	"regexp"
)

// LegalChars ensures that driver name only contains legal characters: alphanumeric, dashes (-), underscores (_), and periods (.)
func LegalChars(driverName string) error {

	if driverName == "" { // Use default value (from data.GCP.Bucket)
		return nil
	}

	match, _ := regexp.MatchString("^[a-zA-Z0-9_.-]+$", driverName)

	if !match {
		return errors.New("Mount name may only contain alphanumeric, dashes (-), underscores (_), and periods (.) characters")
	}

	return nil
}
