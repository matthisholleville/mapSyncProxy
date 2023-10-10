package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/mapsyncproxy/api/client"
	"github.com/matthisholleville/mapsyncproxy/pkg/gcs"
	"github.com/matthisholleville/mapsyncproxy/pkg/haproxy"
	"github.com/prometheus/client_golang/prometheus"
)

type SynchronizeRequestBody struct {
	MapName        string `json:"map_name"  validate:"required,map_name"`
	BucketName     string `json:"bucket_name" validate:"required,bucket_name"`
	BucketFileName string `json:"bucket_file_name" validate:"required,bucket_file_name"`
}

// Synchronize godoc
//
//	@Tags			Synchronization
//	@Summary		Synchronize GCS file to an HAProxy map file.
//	@Description	Synchronize GCS file to an HAProxy map file.
//	@Accept			json
//	@Produce		json
//	@Param		_			body	SynchronizeRequestBody	true	"Data of the synchronisation endpoint"
//
// @Success		200
// @Failure		500		"Internal Server Error"
// @Router			/v1/synchronize [post]
func Synchronize(c echo.Context) (err error) {

	mapSyncContext := c.Get("mapSyncContext").(*client.MapSyncProxyAPI)
	requestBody := SynchronizeRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse("Error reading JSON request body."))
	}

	mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("processed", requestBody.MapName)).Inc()

	// Get MapEntries file from GCS
	gcsEntries, err := getGCSJsonFile(mapSyncContext.GCSClientWrapper, requestBody.BucketName, requestBody.BucketFileName)
	if err != nil {
		mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The GCS file could not be downloaded or interpreted."))
	}

	// Check if duplicate keys
	if hasDuplicateKeys(*gcsEntries) {
		mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The GCS file contains duplicate keys."))
	}

	// Get HAProxy entries from map
	haproxyEntries, err := mapSyncContext.HAProxyClient.GetMapEntries(requestBody.MapName)
	if err != nil {
		mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The entries from the HAProxy Map file could not be retrieved or interpreted."))
	}

	// If Not Exist CreateMap
	entriesToBeCreated := findDifference(*gcsEntries, *haproxyEntries)
	for _, entrie := range entriesToBeCreated {
		_, err = mapSyncContext.HAProxyClient.CreateMapEntrie(&entrie, requestBody.MapName)
		if err != nil {
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be created.", entrie.Key)))
		}
		mapSyncContext.ServerMetrics.MapEntriesTotalCount.With(setMetricsStatusLabels("created", requestBody.MapName)).Inc()
	}

	// If Not Exist in gcs file DeleteMap
	entriesToBeDeleted := findDifference(*haproxyEntries, *gcsEntries)
	for _, entrie := range entriesToBeDeleted {
		_, err = mapSyncContext.HAProxyClient.DeleteMapEntrie(&entrie, requestBody.MapName)
		if err != nil {
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be deleted.", entrie.Key)))
		}
		mapSyncContext.ServerMetrics.MapEntriesTotalCount.With(setMetricsStatusLabels("deleted", requestBody.MapName)).Inc()
	}

	// If Exist and not already processed UpdateMap
	entriesAlreadyProcessed := append(entriesToBeCreated, entriesToBeDeleted...)
	entriesToBeUpdated := findDifference(*gcsEntries, *&entriesAlreadyProcessed)
	for _, entrie := range entriesToBeUpdated {
		_, err = mapSyncContext.HAProxyClient.UpdateMapEntrie(&entrie, requestBody.MapName)
		if err != nil {
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", requestBody.MapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be updated.", entrie.Key)))
		}
		mapSyncContext.ServerMetrics.MapEntriesTotalCount.With(setMetricsStatusLabels("updated", requestBody.MapName)).Inc()
	}

	// Return success
	mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("success", requestBody.MapName)).Inc()
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
