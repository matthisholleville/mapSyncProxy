package haproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var controllerUrl = "services/haproxy/runtime/maps_entries"

func (c *Client) GetMapEntries(mapName string) (*[]MapEntrie, error) {
	url := fmt.Sprintf(
		"%s/%s?map=%s",
		c.serverIP,
		controllerUrl,
		mapName,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res := []MapEntrie{}
	if err := c.executeRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil

}

func (c *Client) CreateMapEntrie(entrie *MapEntrie, mapName string) (*MapEntrie, error) {
	url := fmt.Sprintf(
		"%s/%s?map=%s&force_sync=true",
		c.serverIP,
		controllerUrl,
		mapName,
	)
	bodyStr, _ := json.Marshal(entrie)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res := MapEntrie{}
	if err := c.executeRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) UpdateMapEntrie(entrie *MapEntrie, mapName string) (*MapEntrie, error) {
	url := fmt.Sprintf(
		"%s/%s/%s?map=%s&force_sync=true",
		c.serverIP,
		controllerUrl,
		encodeUrl(entrie.Key),
		mapName,
	)
	bodyStr, _ := json.Marshal(entrie)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res := MapEntrie{}
	if err := c.executeRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteMapEntrie(entrie *MapEntrie, mapName string) (*MapEntrie, error) {
	url := fmt.Sprintf(
		"%s/%s/%s?map=%s&force_sync=true",
		c.serverIP,
		controllerUrl,
		encodeUrl(entrie.Key),
		mapName,
	)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res := MapEntrie{}
	if err := c.executeRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
