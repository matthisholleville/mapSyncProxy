package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/mapsyncproxy/api/client"
	"github.com/rs/zerolog/log"
)

// Generate godoc
//
//	@Tags			Map
//	@Summary		Generate json file from map file.
//	@Description	Generate json file from map file.
//	@Accept			json
//	@Produce		json
//	@Param		map_name	path	string				true	"Map name"//
//
// @Success		200
// @Failure		500		"Internal Server Error"
// @Router			/v1/map/{map_name}/generate [get]
func GenerateJsonFromMap(c echo.Context) (err error) {
	mapSyncContext := c.Get("mapSyncContext").(*client.MapSyncProxyAPI)

	mapName := c.Param("mapName")
	if mapName == "" {
		log.Debug().Err(err).Msg("'map_name' param cannot be empty.")
		return c.JSON(http.StatusInternalServerError, jsonResponse("'map_name' param cannot be empty."))

	}

	mapSyncContext.ServerMetrics.GenerateJsonFromMapTotalCount.With(setMetricsStatusLabels("processed", mapName)).Inc()

	// Get HAProxy entries from map
	haproxyEntries, err := mapSyncContext.HAProxyClient.GetMapEntries(mapName)
	if err != nil {
		log.Debug().Err(err).Msg("The entries from the HAProxy Map file could not be retrieved or interpreted.")
		mapSyncContext.ServerMetrics.GenerateJsonFromMapTotalCount.With(setMetricsStatusLabels("error", mapName)).Inc()
		return c.JSON(http.StatusInternalServerError, jsonResponse("The entries from the HAProxy Map file could not be retrieved or interpreted."))
	}

	mapSyncContext.ServerMetrics.GenerateJsonFromMapTotalCount.With(setMetricsStatusLabels("success", mapName)).Inc()
	return c.JSON(http.StatusOK, haproxyEntries)

}
