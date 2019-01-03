package driver_test

import (
	"testing"

	. "github.com/athera-io/gcs-bucket-helper/internal/athera/validate/driver"

	"github.com/athera-io/gcs-bucket-helper/internal/test/errorcase"
)

func TestLegalChars(t *testing.T) {
	cases := []struct {
		errorcase.Error
		Name        string
		DriverNames []string
	}{
		{
			Name: "Empty name is OK",
			DriverNames: []string{
				"",
			},
		},
		{
			Name: "Name containing alphanumeric, dash, underscore, and period characters is OK",
			DriverNames: []string{
				"Foo.0-1_",
			},
		},
		{
			Name: "Name with illegal special character returns error",
			DriverNames: []string{
				"!",
				"\"",
				"£",
				"$",
				"%",
				"^",
				"&",
				"*",
				"(",
				")",
				"+",
				"=",
				"{",
				"}",
				"[",
				"]",
				",",
				"<",
				">",
				"/",
				"?",
				":",
				";",
				"@",
				"'",
				"~",
				"#",
			},
			Error: errorcase.NewErrorAny(),
		},
	}

	for _, td := range cases {

		for _, driverName := range td.DriverNames {
			err := LegalChars(driverName)

			_ = errorcase.Eval(t, err, td.Error) // No more test output to evaluate
		}
	}
}
