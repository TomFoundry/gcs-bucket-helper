package bucket_test

import (
	"strings"
	"testing"

	. "github.com/athera-io/gcs-bucket-helper/internal/gcp/validate/bucket"

	"github.com/athera-io/gcs-bucket-helper/internal/test/errorcase"
)

func TestLegalChars(t *testing.T) {
	cases := []struct {
		errorcase.Error
		Name        string
		BucketNames []string
	}{
		{
			Name: "Empty name is OK",
			BucketNames: []string{
				"",
			},
		},
		{
			Name: "Name containing lower-case letters, numbers, dashes, underscores, and periods is OK",
			BucketNames: []string{
				"foo.0-1_",
			},
		},
		{
			Name: "Name with upper-case letter returns error",
			BucketNames: []string{
				"Foo",
			},
			Error: errorcase.NewErrorAny(),
		},
		{
			Name: "Name with illegal special character returns error",
			BucketNames: []string{
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

		for _, bucketName := range td.BucketNames {
			err := LegalChars(bucketName)

			_ = errorcase.Eval(t, err, td.Error) // No more test output to evaluate
		}
	}
}

func TestStartAndEndAlphanumeric(t *testing.T) {
	cases := []struct {
		errorcase.Error
		Name       string
		BucketName string
	}{
		{
			Name:       "Empty name is OK",
			BucketName: "",
		},
		{
			Name:       "Starts and ends with letter is OK",
			BucketName: "hello",
		},
		{
			Name:       "Starts and ends with number is OK",
			BucketName: "0hello1",
		},
		{
			Name:       "Starts with period returns error",
			BucketName: ".foo",
			Error:      errorcase.NewErrorAny(),
		},
		{
			Name:       "Ends with period returns error",
			BucketName: "foo.",
			Error:      errorcase.NewErrorAny(),
		},
		{
			Name:       "Starts with hyphen returns error",
			BucketName: "-bar",
			Error:      errorcase.NewErrorAny(),
		},
		{
			Name:       "Ends with hyphen returns error",
			BucketName: "bar-",
			Error:      errorcase.NewErrorAny(),
		},
		{
			Name:       "Starts with underscore returns error",
			BucketName: "_baz",
			Error:      errorcase.NewErrorAny(),
		},
		{
			Name:       "Ends with underscore returns error",
			BucketName: "baz_",
			Error:      errorcase.NewErrorAny(),
		},
	}

	for _, td := range cases {

		err := StartAndEndAlphanumeric(td.BucketName)

		_ = errorcase.Eval(t, err, td.Error) // No more test output to evaluate
	}
}

func TestLength(t *testing.T) {
	cases := []struct {
		errorcase.Error
		Name       string
		BucketName string
	}{
		{
			Name:       "Empty name is OK",
			BucketName: "",
		},
		{
			Name:       "3 character name is OK",
			BucketName: "foo",
		},
		{
			Name:       "63 character name is OK",
			BucketName: strings.Repeat("a", 63),
		},
		{
			Name:       "2 character name returns error",
			BucketName: "aa",
			Error:      errorcase.NewErrorAny(),
		},
		{
			Name:       "64 character name is OK",
			BucketName: strings.Repeat("a", 64),
			Error:      errorcase.NewErrorAny(),
		},
	}

	for _, td := range cases {

		err := Length(td.BucketName)

		_ = errorcase.Eval(t, err, td.Error) // No more test output to evaluate
	}
}

func TestIllegalSubstrings(t *testing.T) {
	cases := []struct {
		errorcase.Error
		Name       string
		BucketName string
	}{
		{
			Name:       "Empty name is OK",
			BucketName: "",
		},
		{
			Name:       "Standard name is OK",
			BucketName: "foo",
		},
		{
			Name:       "Prefix 'goog' returns error",
			BucketName: "googfoo",
			Error:      errorcase.NewErrorAny(),
		},
	}

	for _, td := range cases {

		err := IllegalSubstrings(td.BucketName)

		_ = errorcase.Eval(t, err, td.Error) // No more test output to evaluate
	}
}
