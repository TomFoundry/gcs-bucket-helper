package group_test

import (
	"testing"

	"github.com/athera-io/gcs-bucket-helper/internal/athera/models"
	. "github.com/athera-io/gcs-bucket-helper/internal/athera/validate/group"

	"github.com/athera-io/gcs-bucket-helper/internal/test/errorcase"
)

func newGroups(groups ...*models.Group) []*models.Group {
	return groups
}

func TestLegalIndex(t *testing.T) {

	org := &models.Group{}
	projectA := &models.Group{}
	projectB := &models.Group{}

	cases := []struct {
		errorcase.Error
		Name    string
		Input   string
		Groups  []*models.Group
		Lineage models.GroupLineage
	}{
		{
			Name:   "'1' is OK",
			Input:  "1",
			Groups: newGroups(org),
		},
		{
			Name:   "Max with no lineage is OK",
			Input:  "2",
			Groups: newGroups(projectA, projectB),
		},
		{
			Name:   "Greater than max with no lineage returns error",
			Input:  "3",
			Groups: newGroups(projectA, projectB),
			Error:  errorcase.NewErrorBehavior(InvalidIndexError, IsInvalidIndexError),
		},
		{
			Name:    "Max with lineage is OK",
			Input:   "3",
			Groups:  newGroups(projectA, projectB),
			Lineage: models.NewGroupLineage(org),
		},
		{
			Name:   "Greater than max with lineage returns error",
			Input:  "4",
			Groups: newGroups(projectA, projectB),
			Error:  errorcase.NewErrorBehavior(InvalidIndexError, IsInvalidIndexError),
		},
		{
			Name:   "'0' returns error",
			Input:  "0",
			Groups: newGroups(org),
			Error:  errorcase.NewErrorBehavior(InvalidIndexError, IsInvalidIndexError),
		},
		{
			Name:   "Non-integer returns error",
			Input:  "Hello World",
			Groups: newGroups(org),
			Error:  errorcase.NewErrorBehavior(NotIntegerError, IsNotIntegerError),
		},
	}

	for _, td := range cases {

		err := LegalIndex(td.Input, td.Groups, td.Lineage)

		_ = errorcase.Eval(t, err, td.Error) // No more test output to evaluate
	}
}
