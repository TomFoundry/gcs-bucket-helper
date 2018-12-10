package athera

import (
	"fmt"
	"strconv"

	"github.com/athera-io/gcs-bucket-helper/internal/athera/models"
	"github.com/athera-io/gcs-bucket-helper/internal/input"
	"github.com/pkg/errors"
)

func (a *Athera) selectGroup(org *models.Group) (*models.Group, error) {
	lineage := models.GroupLineage{
		org,
	}

	for {
		leaf := lineage.Leaf()

		children, err := a.client.GetGroupChildren(leaf.ID)

		if err != nil {
			return nil, errors.Wrap(err, "Failed getting groups")
		}

		if len(children) == 0 {
			return leaf, nil
		}

		selectedGroup := selectGroupFromInput(children, lineage)

		if selectedGroup == leaf {
			return leaf, nil
		}

		lineage = append(lineage, selectedGroup)
	}
}

func selectGroupFromInput(groups []*models.Group, lineage models.GroupLineage) *models.Group {

	s := input.Recv(
		buildSelectGroupMsg(groups, lineage),
		func(s string) error {
			// Validator: Must be integer greater than zero, and less than 1 + length of groups

			i, err := strconv.Atoi(s)

			if err != nil {
				return errors.New("Input must be an integer")
			}

			if i < 1 {
				return errors.New("Input must be greater than zero")
			}

			max := len(groups)

			if len(lineage) > 0 {
				max++
			}

			if i > max {
				return fmt.Errorf("Input must be less than or equal to %d", max)
			}

			return nil
		},
	)

	i, _ := strconv.Atoi(s) // Validator already checked error

	if i == len(groups)+1 {
		return lineage.Leaf()
	}

	// index = i - 1 (because input starts from 1, but idx starts from 0)

	return groups[i-1] // Validator already checked that index is legal
}

func buildSelectGroupMsg(groups []*models.Group, lineage models.GroupLineage) string {
	var msg string

	if len(lineage) > 0 {
		msg += "Selected Context: "

		for i, parent := range lineage {
			msg += parent.Name

			if i < len(lineage)-1 {
				msg += " | "
			}
		}

		msg += "\nSelect Group:\n"
	} else {
		msg += "Select Org:\n"
	}

	for i, group := range groups {
		msg += fmt.Sprintf("%d) %s", i+1, group.Name)

		if i < len(groups)-1 {
			msg += "\n"
		}
	}

	if len(lineage) > 0 {
		msg += fmt.Sprintf("\n%d) [Mount in %s]", len(groups)+1, lineage[0].Name)
	}

	return msg
}
