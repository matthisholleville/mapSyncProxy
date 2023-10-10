// Package handlers for monitoring api call.
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health godoc
//
//	@Tags			Monitoring
//	@Summary		Health.
//	@Description	Health.
//	@Accept			*/*
//	@Produce		json
//	@Success		200
//	@Router			/healthz [get]
func Health(echo echo.Context) error {
	return echo.JSONPretty(http.StatusOK, "OK", "")
}

// Ready godoc
//
//	@Tags			Monitoring
//	@Summary		Ready.
//	@Description	Ready.
//	@Accept			*/*
//	@Produce		json
//	@Success		200
//	@Router			/readyz [get]
func Ready(echo echo.Context) error {
	return echo.JSONPretty(http.StatusOK, "OK", "")
}
