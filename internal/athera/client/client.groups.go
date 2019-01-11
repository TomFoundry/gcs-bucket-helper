package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/athera-io/gcs-bucket-helper/internal/athera/models"
	"github.com/pkg/errors"
)

// groupsResponse ...
type groupsResponse struct {
	Groups []*models.Group `json:"groups,omitempty"`
}

// GetOrgs ...
func (c *Client) GetOrgs() ([]*models.Group, error) {

	relativeURL := &url.URL{Path: "/api/v1/orgs"}

	u := c.baseURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, errors.Wrap(err, "Failed creating request")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer: "+c.tok.AccessToken)

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Failed sending request")
	}

	if !isSuccessStatus(res.StatusCode) {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	var groupsRes groupsResponse
	err = json.NewDecoder(res.Body).Decode(&groupsRes)

	if err != nil {
		return nil, errors.Wrap(err, "Failed unmarshaling response")
	}

	return groupsRes.Groups, nil
}

// GetGroupChildren ...
func (c *Client) GetGroupChildren(groupID string) ([]*models.Group, error) {

	relativeURL := &url.URL{Path: fmt.Sprintf("/api/v1/groups/%s/children", groupID)}

	u := c.baseURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("GET", u.String(), nil)

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

	var groupsRes groupsResponse
	err = json.NewDecoder(res.Body).Decode(&groupsRes)

	if err != nil {
		return nil, errors.Wrap(err, "Failed unmarshaling response")
	}

	return groupsRes.Groups, nil
}
