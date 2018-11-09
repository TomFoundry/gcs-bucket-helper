package gcp

// UserData contains information about GCP user
type UserData struct {
	Email string `json:"email"`
}

// DoData contains information needed to create a service account and storage bucket
type DoData struct {
	UserData UserData
	// Project is the name of the existing GCP project
	Project string
	// Bucket is the name for the bucket to create (it must be globally unique)
	Bucket string
}
