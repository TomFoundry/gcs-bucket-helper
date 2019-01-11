package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/athera-io/gcs-bucket-helper/internal/athera/models"
	"github.com/pkg/errors"
)

// CreateStorageDriver ...
func (c *Client) CreateStorageDriver(groupID, driverName, bucketName string, serviceAccountKey []byte) (*models.StorageDriver, error) {

	relativeURL := &url.URL{Path: "/api/v1/storage/driver"}

	u := c.baseURL.ResolveReference(relativeURL)

	encodedServiceAccountKey := base64.StdEncoding.EncodeToString(serviceAccountKey)

	createReq := &createStorageDriverRequest{
		Name: driverName,
		GCS: &gcsConfig{
			BucketID:     bucketName,
			ClientSecret: encodedServiceAccountKey,
		},
		Option: &createStorageDriverOption{
			CreateDefaultMount: true,
		},
	}

	payload, _ := json.Marshal(createReq)

	reader := bytes.NewReader(payload)

	req, err := http.NewRequest("POST", u.String(), reader)

	if err != nil {
		return nil, errors.Wrap(err, "Failed creating request")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer: "+c.tok.AccessToken)

	req.Header.Set("active-group", groupID)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Failed sending request")
	}

	if !isSuccessStatus(res.StatusCode) {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	driver := &models.StorageDriver{}
	err = json.NewDecoder(res.Body).Decode(driver)

	if err != nil {
		return nil, errors.Wrap(err, "Failed unmarshaling response")
	}

	return driver, nil
}

type createStorageDriverRequest struct {
	Name   string                     `json:"name"`
	GCS    *gcsConfig                 `json:"gcs"`
	Option *createStorageDriverOption `json:"option"`
}

type gcsConfig struct {
	BucketID     string `json:"bucket_id"`
	ClientSecret string `json:"client_secret"`
}

type createStorageDriverOption struct {
	CreateDefaultMount bool `json:"create_default_mount"`
}
