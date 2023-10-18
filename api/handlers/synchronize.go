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
	"github.com/rs/zerolog/log"
)

type SynchronizeRequestBody struct {
	BucketName     string `json:"bucket_name" validate:"required,bucket_name"`
	BucketFileName string `json:"bucket_file_name" validate:"required,bucket_file_name"`
}

// Synchronize godoc
//
//	@Tags			Map
//	@Summary		Synchronize GCS file to an HAProxy map file.
//	@Description	Synchronize GCS file to an HAProxy map file.
//	@Accept			json
//	@Produce		json
//	@Param		_			body	SynchronizeRequestBody	true	"Data of the synchronisation endpoint"
//	@Param		map_name	path	string				true	"Map name"//
//
// @Success		200
// @Failure		500		"Internal Server Error"
// @Router			/v1/map/{map_name}/synchronize [post]
func Synchronize(c echo.Context) (err error) {

	mapName := c.Param("mapName")
	if mapName == "" {
		log.Debug().Err(err).Msg("'map_name' param cannot be empty.")
		return c.JSON(http.StatusInternalServerError, jsonResponse("'map_name' param cannot be empty."))

	}

	mapSyncContext := c.Get("mapSyncContext").(*client.MapSyncProxyAPI)
	requestBody := SynchronizeRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, jsonResponse("Error reading JSON request body."))
	}

	mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("processed", mapName)).Inc()

	gcsEntries := &[]haproxy.MapEntrie{}

	if requestBody.BucketFileName == "*" {
		log.Info().Msgf("Multiple GCS files from %s bucket will be downloaded.", requestBody.BucketName)
		// Get MapEntries files from GCS
		gcsEntries, err = downloadMultipleFiles(mapSyncContext.GCSClientWrapper, requestBody.BucketName)
		if err != nil {
			log.Debug().Err(err).Msg("The GCS files could not be listed.")
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse("The GCS files could not be listed."))
		}

	} else {
		log.Info().Msgf("The GCS file %s from the %s bucket will be downloaded", requestBody.BucketFileName, requestBody.BucketName)
		// Get MapEntries file from GCS
		gcsEntries, err = getGCSJsonFile(mapSyncContext.GCSClientWrapper, requestBody.BucketName, requestBody.BucketFileName)
		if err != nil {
			log.Debug().Err(err).Msg("The GCS file could not be downloaded or interpreted.")
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse("The GCS file could not be downloaded or interpreted."))
		}

	}

	// Check if duplicate keys
	if hasDuplicateKeys(*gcsEntries) {
		log.Debug().Msg("The GCS file could not be downloaded or interpreted.")
		mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The GCS file contains duplicate keys."))
	}

	// Get HAProxy entries from map
	haproxyEntries, err := mapSyncContext.HAProxyClient.GetMapEntries(mapName)
	if err != nil {
		log.Debug().Err(err).Msg("The entries from the HAProxy Map file could not be retrieved or interpreted.")
		mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The entries from the HAProxy Map file could not be retrieved or interpreted."))
	}

	// If Not Exist CreateMap
	entriesToBeCreated := findDifference(*gcsEntries, *haproxyEntries, "")
	for _, entrie := range entriesToBeCreated {
		_, err = mapSyncContext.HAProxyClient.CreateMapEntrie(&entrie, mapName)
		if err != nil {
			log.Debug().Err(err).Msgf("The '%s' entry could not be created.", entrie.Key)
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be created.", entrie.Key)))
		}
		mapSyncContext.ServerMetrics.MapEntriesTotalCount.With(setMetricsStatusLabels("created", mapName)).Inc()
	}

	// If Not Exist in gcs file DeleteMap
	entriesToBeDeleted := findDifference(*haproxyEntries, *gcsEntries, "")
	for _, entrie := range entriesToBeDeleted {
		_, err = mapSyncContext.HAProxyClient.DeleteMapEntrie(&entrie, mapName)
		if err != nil {
			log.Debug().Err(err).Msgf("The '%s' entry could not be deleted.", entrie.Key)
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be deleted.", entrie.Key)))
		}
		mapSyncContext.ServerMetrics.MapEntriesTotalCount.With(setMetricsStatusLabels("deleted", mapName)).Inc()
	}

	// If Exist and not already processed UpdateMap
	entriesAlreadyProcessed := append(entriesToBeCreated, entriesToBeDeleted...)
	entriesNotProcessed := findDifference(*gcsEntries, *&entriesAlreadyProcessed, "")
	entriesToBeUpdated := findDifference(entriesNotProcessed, *haproxyEntries, "full")
	for _, entrie := range entriesToBeUpdated {
		_, err = mapSyncContext.HAProxyClient.UpdateMapEntrie(&entrie, mapName)
		if err != nil {
			log.Debug().Err(err).Msgf("The '%s' entry could not be updated.", entrie.Key)
			mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
			return c.JSON(http.StatusInternalServerError, jsonResponse(fmt.Sprintf("The '%s' entry could not be updated.", entrie.Key)))
		}
		mapSyncContext.ServerMetrics.MapEntriesTotalCount.With(setMetricsStatusLabels("updated", mapName)).Inc()
	}

	// Return success
	log.Info().Msgf("Synchronization success. %d created - %d updated - %d deleted", len(entriesToBeCreated), len(entriesToBeUpdated), len(entriesToBeDeleted))
	mapSyncContext.ServerMetrics.SynchronizationTotalCount.With(setMetricsStatusLabels("success", mapName)).Inc()
	return c.JSON(http.StatusOK, jsonResponse("synchronization success."))
}

func downloadMultipleFiles(g *gcs.GCSClientWrapper, bucketName string) (*[]haproxy.MapEntrie, error) {
	gcsFiles, err := g.ListFiles(bucketName)
	if err != nil {
		return nil, err
	}
	result := []haproxy.MapEntrie{}
	for _, file := range *gcsFiles {
		if file.ContentType == "application/json" {
			gcsEntries, err := getGCSJsonFile(g, bucketName, file.Name)
			if err != nil {
				return nil, err
			}
			log.Info().Msgf("%s file downloaded successfully. %d entrie(s) found", file.Name, len(*gcsEntries))
			result = append(result, *gcsEntries...)
		}

	}
	return &result, nil
}

func getGCSJsonFile(g *gcs.GCSClientWrapper, bucketName, fileName string) (*[]haproxy.MapEntrie, error) {
	rc, err := g.DownloadFile(bucketName, fileName)
	if err != nil {
		log.Err(err).Msgf("Unable to download %s", fileName)
		return nil, err
	}
	data, err := io.ReadAll(rc)
	if err != nil {
		log.Err(err).Msgf("Unable to read %s", fileName)
		return nil, err
	}
	var mapEntries []haproxy.MapEntrie
	err = json.Unmarshal(data, &mapEntries)
	if err != nil {
		log.Err(err).Msgf("Unable to Unmarshal %s", fileName)
		return nil, err
	}
	return &mapEntries, nil

}

func findDifference(array1, array2 []haproxy.MapEntrie, diffType string) []haproxy.MapEntrie {
	difference := []haproxy.MapEntrie{}

	map2 := make(map[string]haproxy.MapEntrie)
	for _, item := range array2 {
		map2[item.Key] = item
	}

	if diffType == "full" {
		for _, item := range array1 {
			if _, exists := map2[item.Key]; exists {
				if map2[item.Key].Value != item.Value {
					difference = append(difference, item)
				}
			} else {
				difference = append(difference, item)
			}
		}
	} else {
		for _, item := range array1 {
			if _, exists := map2[item.Key]; !exists {
				difference = append(difference, item)
			}
		}

	}

	return difference
}

func setMetricsStatusLabels(status, mapName string) prometheus.Labels {
	return prometheus.Labels{"status": status, "map_name": mapName}
}
