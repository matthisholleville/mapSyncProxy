package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/mapsyncproxy/pkg/gcs"
	"github.com/matthisholleville/mapsyncproxy/pkg/haproxy"
	"github.com/matthisholleville/mapsyncproxy/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type SynchronizeRequestBody struct {
	MapName        string `json:"map_name"`
	BucketName     string `json:"bucket_name"`
	BucketFileName string `json:"bucket_file_name"`
}

func (r *SynchronizeRequestBody) AreFieldsNotEmpty() bool {
	return r.MapName != "" && r.BucketName != "" && r.BucketFileName != ""
}

func Synchronize(c echo.Context, h *haproxy.Client, g *gcs.GCSClientWrapper, m *metrics.ServerMetrics) (err error) {

	requestBody := SynchronizeRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse("Error reading JSON request body."))
	}

	if !requestBody.AreFieldsNotEmpty() {
		return c.JSON(http.StatusBadRequest, jsonResponse("All fields (map_name, bucket_name, bucket_file_name) are required in the request body."))
	}

	m.SynchronizationTotalCount.With(setMetricsStatusLabels("processed", requestBody.MapName)).Inc()

	// Get MapEntries file from GCS
	gcsEntries, err := getGCSJsonFile(g, requestBody.BucketName, requestBody.BucketFileName)
	if err != nil {
		m.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The GCS file could not be downloaded or interpreted."))
	}

	// Get HAProxy entries from map
	haproxyEntries, err := h.GetMapEntries(requestBody.MapName)
	if err != nil {
		m.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The entries from the HAProxy Map file could not be retrieved or interpreted."))
	}

	// If Not Exist CreateMap
	entriesToBeCreated := findDifference(*gcsEntries, *haproxyEntries)
	for _, entrie := range entriesToBeCreated {
		_, err = h.CreateMapEntrie(&entrie, requestBody.MapName)
		if err != nil {
			m.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be created.", entrie.Key)))
		}
		m.MapEntriesTotalCount.With(setMetricsStatusLabels("created", requestBody.MapName)).Inc()
	}

	// If Not Exist in gcs file DeleteMap
	entriesToBeDeleted := findDifference(*haproxyEntries, *gcsEntries)
	for _, entrie := range entriesToBeDeleted {
		_, err = h.DeleteMapEntrie(&entrie, requestBody.MapName)
		if err != nil {
			m.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be deleted.", entrie.Key)))
		}
		m.MapEntriesTotalCount.With(setMetricsStatusLabels("deleted", requestBody.MapName)).Inc()
	}

	// If Exist and not already processed UpdateMap
	entriesAlreadyProcessed := append(entriesToBeCreated, entriesToBeDeleted...)
	entriesToBeUpdated := findDifference(*gcsEntries, *&entriesAlreadyProcessed)
	for _, entrie := range entriesToBeUpdated {
		_, err = h.UpdateMapEntrie(&entrie, requestBody.MapName)
		if err != nil {
			m.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be updated.", entrie.Key)))
		}
		m.MapEntriesTotalCount.With(setMetricsStatusLabels("updated", requestBody.MapName)).Inc()
	}

	// Return success
	m.SynchronizationTotalCount.With(setMetricsStatusLabels("success", requestBody.MapName)).Inc()
	return c.JSON(http.StatusOK, jsonResponse("synchronization success."))
}

func getGCSJsonFile(g *gcs.GCSClientWrapper, bucketName, fileName string) (*[]haproxy.MapEntrie, error) {
	rc, err := g.DownloadFile(bucketName, fileName)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	var mapEntries []haproxy.MapEntrie
	err = json.Unmarshal(data, &mapEntries)
	if err != nil {
		return nil, err
	}
	return &mapEntries, nil

}

func findDifference(array1, array2 []haproxy.MapEntrie) []haproxy.MapEntrie {
	difference := []haproxy.MapEntrie{}

	map2 := make(map[string]haproxy.MapEntrie)
	for _, item := range array2 {
		map2[item.Key] = item
	}

	for _, item := range array1 {
		if _, exists := map2[item.Key]; !exists {
			difference = append(difference, item)
		}
	}

	return difference
}

func setMetricsStatusLabels(status, mapName string) prometheus.Labels {
	return prometheus.Labels{"status": status, "map_name": mapName}
}
