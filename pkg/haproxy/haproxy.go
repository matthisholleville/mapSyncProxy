package haproxy

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
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
		HTTPClient: resty.New().
			SetBaseURL(fmt.Sprintf("%s://%s/v2", scheme, serverIP)).
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json; charset=utf-8").
			SetBasicAuth(username, password).
			SetRetryCount(retryCount).
			SetRetryWaitTime(retryWaitTime * time.Second),
		serverIP: fmt.Sprintf("%s://%s/v2", scheme, serverIP),
	}
}
