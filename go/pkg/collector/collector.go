package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	round     = 0
	namespace = "res_control"
	// Metrics
	metricConfig = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "config"),
		"Main configuration details",
		[]string{"version", "min_pixels", "max_pixels", "step_pixels"}, nil,
	)
	metricMeasurementRound = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "measurement_round"),
		"Measurement round",
		[]string{}, nil,
	)
	metricMeasurementPixels = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "measurement_pixels"),
		"Measurement pixels generation",
		[]string{}, nil,
	)
)

type Exporter struct {
	version         string
	minPix, maxPix, stepPix int64
}

func NewExporter(version string, minPix int64, maxPix int64, stepPix int64) *Exporter {
	return &Exporter{
		version:    version,
		minPix:     minPix,
		maxPix:     maxPix,
		stepPix:    stepPix,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricConfig
	ch <- metricMeasurementRound
	ch <- metricMeasurementPixels
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	startProm := time.Now()
	e.CollectResizeMetrics(ch)
	end := time.Now()

	log.Info("Collect round ", round, " finished in ", end.Sub(startProm))
	round++
}
