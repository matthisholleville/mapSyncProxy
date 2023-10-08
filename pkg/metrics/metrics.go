package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type ServerMetrics struct {
	MapEntriesTotalCount      *prometheus.CounterVec
	SynchronizationTotalCount *prometheus.CounterVec
}

func New() *ServerMetrics {
	serverMetrics := &ServerMetrics{}

	serverMetrics.MapEntriesTotalCount = createAndRegisterCounter(
		"mapsyncproxy_haproxy_mapentries_total",
		"How many MapEntries processed, partitioned by status and map_name.",
		[]string{"status", "map_name"},
	)

	serverMetrics.SynchronizationTotalCount = createAndRegisterCounter(
		"mapsyncproxy_synchronization_total",
		"How many Synchronization processed, partitioned by status and map_name.",
		[]string{"status", "map_name"},
	)

	return serverMetrics

}

func createAndRegisterCounter(name, help string, labels []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)

	if err := prometheus.Register(counter); err != nil {
		log.Fatal(err)
	}

	return counter
}
