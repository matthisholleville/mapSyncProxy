package haproxy

import (
	"encoding/base64"
	"net/url"
)

func generateBasicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func encodeUrl(s string) string {
	return url.QueryEscape(s)
}
