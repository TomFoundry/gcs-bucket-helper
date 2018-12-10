/*Package executor has an interface representing a facade for executing GCP and Athera operations
 */
package executor

import "golang.org/x/oauth2"

// Executor is a facade for executing GCP and Athera operations
type Executor interface {
	ExecuteGCP(tok *oauth2.Token, userEmail string, onCompleteURL string)
	ExecuteAthera(tok *oauth2.Token)
}

// Data stores information needed to execute operations.
// It is instantiated by the Executor implementation
type Data struct {
	GCP  *GCPData
	User *UserData
}

// NewData ...
func NewData() *Data {
	return &Data{
		GCP:  &GCPData{},
		User: &UserData{},
	}
}

// UserData contains information about GCP user
type UserData struct {
	Email string
}

// GCPData contains information about GCP project entities
type GCPData struct {
	// Project is the name of the existing GCP project
	Project string
	// Bucket is the name for the bucket to create (it must be globally unique)
	Bucket string
	// ServiceAccountPrivateData is the private key data for the service account
	ServiceAccountPrivateData []byte
}
