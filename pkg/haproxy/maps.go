package haproxy

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

var controllerUrl = "/services/haproxy/runtime/maps_entries"

func (c *Client) GetMapEntries(mapName string) (*[]MapEntrie, error) {
	mapEntrie := []MapEntrie{}
	url := fmt.Sprintf(
		"%s?map=%s",
		controllerUrl,
		mapName,
	)
	resp, err := c.HTTPClient.R().
		SetResult(&mapEntrie).
		Get(url)

	if err != nil {
		log.Debug().Err(err).Msg("Error while calling DataplaneAPI.")
		return &mapEntrie, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Debug().Msgf("Error while getting mapEntrie. Status code %d", resp.StatusCode())
		return &mapEntrie, fmt.Errorf("Error while getting mapEntrie: %s", resp.Status())
	}

	return &mapEntrie, nil

}

func (c *Client) CreateMapEntrie(entrie *MapEntrie, mapName string) (*MapEntrie, error) {
	url := fmt.Sprintf(
		"%s?map=%s&force_sync=true",
		controllerUrl,
		mapName,
	)
	mapEntrie := MapEntrie{}
	resp, err := c.HTTPClient.R().
		SetBody(entrie).
		SetResult(mapEntrie).
		Post(url)

	if err != nil {
		log.Debug().Err(err).Msg("Error while calling DataplaneAPI.")
		return &mapEntrie, err
	}

	if resp.StatusCode() != http.StatusCreated {
		log.Debug().Msgf("Error while creating mapEntrie. Status code %d", resp.StatusCode())
		return &mapEntrie, fmt.Errorf("Error while creating mapEntrie: %s", resp.Status())
	}

	return &mapEntrie, nil
}

func (c *Client) UpdateMapEntrie(entrie *MapEntrie, mapName string) (*MapEntrie, error) {
	url := fmt.Sprintf(
		"%s/%s?map=%s&force_sync=true",
		controllerUrl,
		encodeUrl(entrie.Key),
		mapName,
	)
	mapEntrie := MapEntrie{}
	resp, err := c.HTTPClient.R().
		SetBody(entrie).
		SetResult(mapEntrie).
		Put(url)

	if err != nil {
		log.Debug().Err(err).Msg("Error while calling DataplaneAPI.")
		return &mapEntrie, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Debug().Msgf("Error while updating mapEntrie. Status code %d", resp.StatusCode())
		return &mapEntrie, fmt.Errorf("Error while updating mapEntrie: %s", resp.Status())
	}

	return &mapEntrie, nil
}

func (c *Client) DeleteMapEntrie(entrie *MapEntrie, mapName string) (*MapEntrie, error) {
	url := fmt.Sprintf(
		"%s/%s?map=%s&force_sync=true",
		controllerUrl,
		encodeUrl(entrie.Key),
		mapName,
	)
	mapEntrie := MapEntrie{}
	resp, err := c.HTTPClient.R().
		SetBody(entrie).
		SetResult(mapEntrie).
		Delete(url)

	if err != nil {
		log.Debug().Err(err).Msg("Error while deleting DataplaneAPI.")
		return &mapEntrie, err
	}

	if resp.StatusCode() != http.StatusNoContent {
		log.Debug().Msgf("Error while deleting mapEntrie. Status code %d", resp.StatusCode())
		return &mapEntrie, fmt.Errorf("Error while deleting mapEntrie: %s", resp.Status())
	}

	return &mapEntrie, nil
}
