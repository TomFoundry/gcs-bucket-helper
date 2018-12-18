/*Package impl has an Executor implementation.
 */
package impl

import (
	"fmt"
	"log"
	"os"

	"github.com/athera-io/gcs-bucket-helper/internal/athera"
	"github.com/athera-io/gcs-bucket-helper/internal/executor"
	"github.com/athera-io/gcs-bucket-helper/internal/gcp"
	"golang.org/x/oauth2"
)

type executorImpl struct {
	data   *executor.Data
	athera *athera.Athera
	gcp    *gcp.GCP
}

// NewExecutor ...
func NewExecutor(atheraExecutor *athera.Athera, gcpExecutor *gcp.GCP) executor.Executor {
	e := &executorImpl{
		data:   executor.NewData(),
		athera: atheraExecutor,
		gcp:    gcpExecutor,
	}

	return e
}

// ExecuteGCP executes GCP operations: Create service account, create bucket, set IAM policy on bucket
func (e *executorImpl) ExecuteGCP(tok *oauth2.Token, userEmail string, onCompleteURL string) {
	e.data.User.Email = userEmail

	if err := e.gcp.Execute(tok, e.data); err != nil {
		log.Fatal("Failed executing GCP sequence: ", err)
	}

	fmt.Println("- We need to log in to Athera to get permission to connect the new bucket.")
	fmt.Println("- Please navigate to:")
	fmt.Printf(onCompleteURL)
}

// ExecuteAthera executes Athera operation: Create storage driver
func (e *executorImpl) ExecuteAthera(tok *oauth2.Token) {
	if err := e.athera.Execute(tok, e.data); err != nil {
		log.Fatal("Failed executing Athera sequence: ", err)
	}

	fmt.Println("- Process is complete. Your new bucket is now connected to Athera.")
	os.Exit(0)
}
