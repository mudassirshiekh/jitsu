package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var sourcesPoolLabels = []string{"type"}

var (
	sourcesGoroutinesPoolSize *prometheus.GaugeVec
)

func initSourcesPool() {
	sourcesGoroutinesPoolSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "eventnative",
		Subsystem: "sources",
		Name:      "goroutines_pool",
	}, sourcesPoolLabels)
}

func FreeSourcesGoroutines(value int) {
	if Enabled {
		sourcesGoroutinesPoolSize.WithLabelValues("free").Set(float64(value))
	}
}

func RunningSourcesGoroutines(value int) {
	if Enabled {
		sourcesGoroutinesPoolSize.WithLabelValues("running").Set(float64(value))
	}
}
