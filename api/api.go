package api

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matthisholleville/mapsyncproxy/api/client"
	"github.com/matthisholleville/mapsyncproxy/api/handlers"
	v1 "github.com/matthisholleville/mapsyncproxy/api/v1"
	"github.com/matthisholleville/mapsyncproxy/pkg/gcs"
	"github.com/matthisholleville/mapsyncproxy/pkg/haproxy"
	"github.com/matthisholleville/mapsyncproxy/pkg/metrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Server(ctx context.Context, sigs chan os.Signal) {

	viper.AutomaticEnv()
	viper.SetEnvPrefix("MAPSYNCPROXY")
	viper.SetDefault("DATAPLANE_USERNAME", "admin")
	viper.SetDefault("DATAPLANE_PASSWORD", "adminpwd")
	viper.SetDefault("DATAPLANE_HOST", "127.0.0.1:5555")

	s := &client.MapSyncProxyAPI{
		Echo: echo.New(),
		HAProxyClient: haproxy.NewClient(
			viper.GetString("DATAPLANE_USERNAME"),
			viper.GetString("DATAPLANE_PASSWORD"),
			viper.GetString("DATAPLANE_HOST"),
			true,
		),
		GCSClientWrapper: gcs.NewClient(),
		ServerMetrics:    metrics.New(),
	}

	s.Echo.Use(middleware.Logger())
	//CORS
	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	//PROMETHEUS
	s.Echo.Use(echoprometheus.NewMiddleware("mapsyncproxy"))

	s.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.Echo.GET("/metrics", echoprometheus.NewHandler())
	s.Echo.GET("/healthz", handlers.Health)
	s.Echo.GET("/readyz", handlers.Ready)

	v1.API(s)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		err := s.Echo.Start(fmt.Sprintf(":%s", port))
		if err != nil {
			s.Echo.Logger.Fatal(err)
		}
	}()

	// wait for signals to shutdown
	<-sigs
	log.Info().Msg("shutting down the API server")

	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout*time.Second)

	defer cancel()

	if err := s.Echo.Shutdown(ctxTimeout); err != nil {
		s.Echo.Logger.Fatal(err)
	}

}
