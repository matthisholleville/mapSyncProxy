package haproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func NewClient(
	username string,
	password string,
	serverIP string,
	insecure bool,
) *Client {
	scheme := "https"
	if insecure {
		scheme = "http"
	}

	return &Client{
		username: username,
		password: password,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		serverIP: fmt.Sprintf("%s://%s/v2", scheme, serverIP),
	}
}

func (c *Client) executeRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", generateBasicAuthHeader(c.username, c.password)))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("HTTP Error: Resource Not Found (Status %d)", res.StatusCode)
	}

	if res.StatusCode >= 300 {
		var errResponse errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errResponse); err == nil {
			errors.New(errResponse.Message)
		}
		return fmt.Errorf("HTTP Error: Unknown Error (Status %d)", res.StatusCode)
	}

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	if v == nil {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
