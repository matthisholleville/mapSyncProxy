package haproxy

import "net/http"

type Client struct {
	username   string
	password   string
	serverIP   string
	HTTPClient *http.Client
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MapEntrie struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
