package haproxy

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	username   string
	password   string
	serverIP   string
	HTTPClient *resty.Client
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
