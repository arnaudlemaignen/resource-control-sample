package collector

import (
	"time"
	"os"
	"image"
	"image/png"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	round     = 0
	imageIn = "resources/fine_387x248.png"
	imageOut = "output/resized.png"
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
		"Measurement of the image resizing in pixels",
		[]string{}, nil,
	)
	metricMeasurementWrittenBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "measurement_written_bytes"),
		"Measurement of written bytes when resizing",
		[]string{}, nil,
	)
	metricMeasurementThroughWrittenBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "measurement_through_written_bytes_secs"),
		"Measurement of the IO write bytes throughput when resizing",
		[]string{}, nil,
	)
)

type Exporter struct {
	version         		string
	minPix, maxPix, stepPix int
	src 					image.Image
}

func NewExporter(version string, minPix,maxPix,stepPix int) *Exporter {
	startInit := time.Now()
	log.Info("Init Collector")
	input, _ := os.Open(imageIn)
	defer input.Close()

	src, _ := png.Decode(input)
	log.Info("Init Collector finished in ", time.Now().Sub(startInit))

	return &Exporter{
		version:    version,
		minPix:     minPix,
		maxPix:     maxPix,
		stepPix:    stepPix,
		src: 		src,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricConfig
	ch <- metricMeasurementRound
	ch <- metricMeasurementPixels
	ch <- metricMeasurementWrittenBytes
	ch <- metricMeasurementThroughWrittenBytes
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Info("Collect round ", round)
	e.CollectResizeMetrics(ch)
	log.Info("Collect round ", round, " finished")
	round++
}
