package client

import (
	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/mapsyncproxy/pkg/gcs"
	"github.com/matthisholleville/mapsyncproxy/pkg/haproxy"
	"github.com/matthisholleville/mapsyncproxy/pkg/metrics"
	"github.com/spf13/viper"
)

type MapSyncProxyAPI struct {
	Echo             *echo.Echo
	HAProxyClient    *haproxy.Client
	GCSClientWrapper *gcs.GCSClientWrapper
	ServerMetrics    *metrics.ServerMetrics
}

func New() *MapSyncProxyAPI {

	viper.AutomaticEnv()
	viper.SetEnvPrefix("MAPSYNCPROXY")
	viper.SetDefault("DATAPLANE_USERNAME", "admin")
	viper.SetDefault("DATAPLANE_PASSWORD", "adminpwd")
	viper.SetDefault("DATAPLANE_HOST", "127.0.0.1:5555")

	return &MapSyncProxyAPI{
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
}
