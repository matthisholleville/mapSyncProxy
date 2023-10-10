package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/mapsyncproxy/api/client"
	"github.com/matthisholleville/mapsyncproxy/api/handlers"

	// Import docs for swagger.
	_ "github.com/matthisholleville/mapsyncproxy/docs"
)

// API route.
func API(s *client.MapSyncProxyAPI) {
	s.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mapSyncContext", s)
			return next(c)
		}
	})
	v1Api := s.Echo.Group("/v1")

	// map endpoints
	v1Api.POST("/map/:mapName/synchronize", handlers.Synchronize)
	v1Api.GET("/map/:mapName/generate", handlers.GenerateJsonFromMap)
}
