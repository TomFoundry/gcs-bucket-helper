package athera

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/athera-io/gcs-bucket-helper/internal/athera/client"
	"github.com/athera-io/gcs-bucket-helper/internal/athera/models"
	"github.com/athera-io/gcs-bucket-helper/internal/executor"
	"github.com/athera-io/gcs-bucket-helper/internal/input"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Athera is responsible for executing Athera operations
type Athera struct {
	voyagerURL *url.URL
	// client is instantiated after we have a token
	client *client.Client
}

// New ...
func New(voyagerURL string) (*Athera, error) {
	u, err := url.Parse(voyagerURL)

	if err != nil {
		return nil, err
	}

	a := &Athera{
		voyagerURL: u,
	}

	return a, nil
}

// Execute ...
func (a *Athera) Execute(tok *oauth2.Token, data *executor.Data) error {

	a.client = client.New(a.voyagerURL, tok)

	orgs, err := a.client.GetOrgs()

	if err != nil {
		return errors.Wrap(err, "Failed getting orgs")
	}

	fmt.Println("") // Line break for legibility
	fmt.Println("- We are going to create a storage driver linked to your GCS bucket")
	fmt.Println("- First we need to select a group to link with the storage driver")

	org := selectGroupFromInput(orgs, nil)

	selectedGroup, err := a.selectGroup(org)

	if err != nil {
		return err
	}

	fmt.Println("- Selected group", selectedGroup.Name)

	return a.createDriver(data, selectedGroup)
}

func (a *Athera) createDriver(data *executor.Data, group *models.Group) error {

	driverName := input.Recv(
		"Please choose a name for your storage driver",
		// Validator: Only legal characters are alphanumeric, "_", "-", and "."
		func(s string) error {
			match, _ := regexp.MatchString("^[a-zA-Z0-9_.-]+$", s)

			if !match {
				return errors.New("Mount name may only contain alphanumeric, hyphen ('-'), underscore ('_'), and period ('.') characters")
			}

			return nil
		},
	)

	_, err := a.client.CreateStorageDriver(group.ID, driverName, data.GCP.Bucket, data.GCP.ServiceAccountPrivateData)

	if err != nil {
		return errors.Wrap(err, "Failed creating storage driver")
	}

	return nil
}
