package gcp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	admin "cloud.google.com/go/iam/admin/apiv1"
	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
)

// Do ...
func Do(tok *oauth2.Token, userData UserData) {

	ctx := context.Background() // Make new context so that server can return authentication success to user

	doData := DoData{
		UserData: userData,
	}

	tokSource := oauth2.StaticTokenSource(tok)

	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(tokSource))

	if err != nil {
		log.Fatal("Failed instantiating GCP Storage Client: ", err)
	}

	fillDoData(ctx, storageClient, &doData)

	serviceAccount, err := makeServiceAccount(ctx, tokSource, storageClient, doData)

	if err != nil {
		log.Fatal(err)
	}

	// Sleep avoids race condition where Google sometimes does not recognize service account that was just created
	time.Sleep(time.Second)

	if err := makeBucket(ctx, storageClient, serviceAccount, doData); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func fillDoData(ctx context.Context, storageClient *storage.Client, doData *DoData) {
	doData.Project = userInput("Enter your GCP project name")

	doData.Bucket = userInput(
		"Enter a name for your new GCS bucket",
		// Validator: Bucket with that name must not already exist
		func(bucketName string) error {
			bucket := storageClient.Bucket(bucketName)
			_, err := bucket.Attrs(ctx)

			// Bucket does not exist
			if err == storage.ErrBucketNotExist {
				return nil
			}

			// Bucket exists and user has access to it
			if err == nil {
				failMsg := "Bucket already exists, and you have permission to access it"
				return errors.New(failMsg)
			}

			// Bucket exists and user does not have access to it
			if errGoogleAPI, ok := err.(*googleapi.Error); ok && errGoogleAPI.Code == http.StatusForbidden {
				failMsg := "Bucket already exists, but you do not have permission to access it"
				return errors.New(failMsg)
			}

			// Something else went wrong (but the bucket does not exist, so validator should pass)
			return nil
		},
		// Validator: Name must not contain illegal characters
		func(bucketName string) error {

			pattern := "^[a-z1-9._-]+$"

			if match, _ := regexp.MatchString(pattern, bucketName); !match {
				failMsg := "Bucket names must contain only lowercase letters, numbers, dashes (-), underscores (_), and dots (.)"
				return errors.New(failMsg)
			}

			return nil
		},
		// Validator: Name must start and end with a letter
		func(bucketName string) error {

			startPattern := "^[a-z1-9]"

			failMsg := "Bucket names must start and end with a number or letter"

			if match, _ := regexp.MatchString(startPattern, bucketName); !match {
				return errors.New(failMsg)
			}

			endPattern := "[a-z1-9]$"

			if match, _ := regexp.MatchString(endPattern, bucketName); !match {
				return errors.New(failMsg)
			}

			return nil
		},
		// Validator: Name must be appropriate length
		func(bucketName string) error {

			if len(bucketName) < 3 || len(bucketName) > 63 {
				failMsg := "Bucket names must contain 3 to 63 characters"
				return errors.New(failMsg)
			}

			return nil
		},
		// Validator: Name must not have illegal substrings
		func(bucketName string) error {

			if strings.HasPrefix(bucketName, "goog") {
				failMsg := "Bucket names cannot begin with the 'goog' prefix"
				return errors.New(failMsg)
			}

			return nil
		},
	)
}

func makeServiceAccount(ctx context.Context, tokSource oauth2.TokenSource, storageClient *storage.Client, doData DoData) (*adminpb.ServiceAccount, error) {

	iamClient, err := admin.NewIamClient(ctx, option.WithTokenSource(tokSource))

	if err != nil {
		return nil, errors.Wrap(err, "Failed instantiating GCP IAM Service Account Client")
	}

	defer func() {
		if err := iamClient.Close(); err != nil {
			log.Print("Error: Failed closing GCP IAM Service Account Client: ", err)
		}
	}()

	createServiceAccountReq := &adminpb.CreateServiceAccountRequest{
		Name:      "projects/" + doData.Project,
		AccountId: genServiceAccountID(doData),
	}

	serviceAccount, err := iamClient.CreateServiceAccount(ctx, createServiceAccountReq)

	if err != nil {
		return nil, errors.Wrap(err, "Failed creating IAM Service Account")
	}

	marshaler := &jsonpb.Marshaler{}

	serviceAccountJSONPath := userInput(
		"Enter a path on your local machine to save your service account credentials",
		// Validator: Directories in path must already exist
		func(s string) error {

			dir, filename := filepath.Split(s)

			if isDirectoryName(filename) {
				stat, err := os.Stat(s)

				if os.IsNotExist(err) {
					return errors.New("No such directory: " + s)
				}

				if !stat.IsDir() {
					failMsg := fmt.Sprintf("File named %s already exists", filename)
					return errors.New(failMsg)
				}
			} else {
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					// Custom error is more readable than error from package os
					return errors.New("No such directory: " + dir)
				}
			}

			return nil
		},
	)

	// Add filename to path if last element in path represents directory name
	_, serviceAccountJSONFilename := filepath.Split(serviceAccountJSONPath)
	if isDirectoryName(serviceAccountJSONFilename) {
		serviceAccountJSONPath = filepath.Join(serviceAccountJSONPath, createServiceAccountReq.AccountId+".json")
	}

	fmt.Println("Saving service account credentials at " + serviceAccountJSONPath)

	f, err := os.Create(serviceAccountJSONPath)

	if err != nil {
		return nil, errors.Wrap(err, "Failed creating Service Account JSON file")
	}

	defer f.Close()

	if err := marshaler.Marshal(f, serviceAccount); err != nil {
		return nil, errors.Wrap(err, "Failed marshaling Service Account to JSON")
	}

	return serviceAccount, nil
}

// isDirectoryName returns true if filename is a directory name.
// (N.B. We assume that if the filename has no extension, then the user intended last element in path to represent a directory)
func isDirectoryName(s string) bool {
	ls := strings.Split(s, ".")

	// No extension if there is no "." char
	return len(ls) == 1 ||
		// No extension if there is only one "." char, and that char is a prefix (because "." prefix represents hidden file, not extension)
		(len(ls) == 2 && strings.HasPrefix(s, "."))
}

func genServiceAccountID(doData DoData) string {
	maxServiceAccountIDLength := 40
	prefix := "athera-"
	bucketName := doData.Bucket

	for {
		// It is safe to measure bytes instead of runes because:
		// "Bucket names must contain only lowercase letters, numbers, dashes (-), underscores (_), and dots (.)".
		// Therefore, all legal characters have length of 1 byte
		// https://cloud.google.com/storage/docs/naming
		if len(prefix+bucketName) > maxServiceAccountIDLength {
			bucketName = bucketName[:len(bucketName)-1]
		} else {
			break
		}
	}

	return prefix + bucketName
}

func makeBucket(ctx context.Context, storageClient *storage.Client, serviceAccount *adminpb.ServiceAccount, doData DoData) error {

	bkt := storageClient.Bucket(doData.Bucket)

	bucketRegion := userInput("Enter a region for your new GCS bucket (e.g. US, EU)")

	serviceAccountEntity := storage.ACLEntity("user-" + serviceAccount.Email)

	ownerEntity := storage.ACLEntity("user-" + doData.UserData.Email)

	bktAttrs := &storage.BucketAttrs{
		Name:     doData.Bucket,
		Location: bucketRegion,
		ACL: []storage.ACLRule{
			storage.ACLRule{
				Entity: serviceAccountEntity,
				Role:   storage.RoleReader,
			},
			storage.ACLRule{
				Entity: serviceAccountEntity,
				Role:   storage.RoleWriter,
			},
			storage.ACLRule{
				Entity: ownerEntity,
				Role:   storage.RoleOwner,
			},
		},
	}

	if err := bkt.Create(ctx, serviceAccount.ProjectId, bktAttrs); err != nil {
		return errors.Wrap(err, "Failed creating GCS Bucket")
	}

	fmt.Printf("Created bucket: https://console.cloud.google.com/storage/browser/%s?project=%s\n", doData.Bucket, serviceAccount.ProjectId)

	return nil
}
