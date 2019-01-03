package serviceaccount_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	. "github.com/athera-io/gcs-bucket-helper/internal/gcp/validate/serviceaccount"

	"github.com/athera-io/gcs-bucket-helper/internal/test/errorcase"
)

func TestPathExists(t *testing.T) {

	workingDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Fatal("Failed getting working directory: ", err)
	}

	filePath := filepath.Join(workingDir, "foo.txt")

	f, err := os.Create(filePath)

	if err != nil {
		t.Fatal("Failed creating dummy file: ", err)
	}

	if err := f.Close(); err != nil {
		t.Fatal("Failed closing dummy file: ", err)
	}

	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Error("Failed deleting dummy file: ", err)
		}
	}()

	cases := []struct {
		errorcase.Error
		Name string
		Path string
	}{
		{
			Name: "Path exists OK",
			Path: workingDir,
		},
		{
			Name:  "Path does not exist returns error",
			Path:  path.Join(workingDir, "foo"),
			Error: errorcase.NewErrorBehavior(NotExistError, IsNotExistError),
		},
		{
			Name:  "Path exists but is not directory returns error",
			Path:  filePath,
			Error: errorcase.NewErrorBehavior(NotDirError, IsNotDirError),
		},
		{
			Name: "Tilde path exists OK",
			Path: "~/Documents",
		},
	}

	for _, td := range cases {

		err := PathExists(td.Path)

		_ = errorcase.Eval(t, err, td.Error) // No more test output to evaluate
	}
}
