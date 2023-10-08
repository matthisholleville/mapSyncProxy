package server

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matthisholleville/mapsyncproxy/controllers"
	"github.com/matthisholleville/mapsyncproxy/pkg/gcs"
	"github.com/matthisholleville/mapsyncproxy/pkg/haproxy"
	"github.com/matthisholleville/mapsyncproxy/pkg/metrics"
	"github.com/spf13/viper"
)

type MapSyncProxyAPI struct {
	e *echo.Echo
	h *haproxy.Client
	g *gcs.GCSClientWrapper
	m *metrics.ServerMetrics
}

func New() *MapSyncProxyAPI {

	viper.AutomaticEnv()
	viper.SetEnvPrefix("MAPSYNCPROXY")
	viper.SetDefault("DATAPLANE_USERNAME", "admin")
	viper.SetDefault("DATAPLANE_PASSWORD", "adminpwd")
	viper.SetDefault("DATAPLANE_HOST", "127.0.0.1:5555")

	return &MapSyncProxyAPI{
		e: echo.New(),
		h: haproxy.NewClient(
			viper.GetString("DATAPLANE_USERNAME"),
			viper.GetString("DATAPLANE_PASSWORD"),
			viper.GetString("DATAPLANE_HOST"),
			true,
		),
		g: gcs.NewClient(),
		m: metrics.New(),
	}
}

// Start server functionality
func (s *MapSyncProxyAPI) Start(port string) {
	// logger
	s.e.Use(middleware.Logger())
	//CORS
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	//PROMETHEUS
	s.e.Use(echoprometheus.NewMiddleware("mapsyncproxy"))
	s.e.GET("/metrics", echoprometheus.NewHandler())

	// synchronize endpoint
	s.e.POST("/synchronize", func(c echo.Context) error {
		return controllers.Synchronize(c, s.h, s.g, s.m)
	})
	// Start Server
	s.e.Logger.Fatal(s.e.Start(port))
}

// Close server functionality
func (s *MapSyncProxyAPI) Close() {
	s.e.Close()
}
