package collector

import (
	"time"
	"os"
	"strconv"
	"golang.org/x/image/draw"
	"image"
	"image/png"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (e *Exporter) CollectResizeMetrics(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metricConfig, prometheus.GaugeValue, 1, e.version, strconv.Itoa(int(e.minPix)), strconv.Itoa(int(e.maxPix)), strconv.Itoa(int(e.stepPix)),
	)
	ch <- prometheus.MustNewConstMetric(
		metricMeasurementRound, prometheus.GaugeValue, float64(round), 
	)

	e.ResizeImage(ch)
}

//https://github.com/MariaLetta/free-gophers-pack
//https://stackoverflow.com/questions/22940724/go-resizing-images
func (e *Exporter) ResizeImage(ch chan<- prometheus.Metric) {
	startRound := time.Now()
	log.Info("Begin measurement of Round : ", round)
	input, _ := os.Open("resources/fine_387x248.png")
	defer input.Close()
	
	output, _ := os.Create("resized.png")
	defer output.Close()
	
	// Decode the image (from PNG to image.Image):
	src, _ := png.Decode(input)
	
	//every 10 rounds we add 1 step
	tens := round % 10
	x := src.Bounds().Max.X+int(e.stepPix) * tens
	y := src.Bounds().Max.Y+int(e.stepPix) * tens
	totalPixels := x * y
	log.Info(" Round : ", round, " tens : ",tens, " pixels ",totalPixels)
	ch <- prometheus.MustNewConstMetric(
		metricMeasurementPixels, prometheus.GaugeValue, float64(totalPixels), 
	)

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, x, y))
	
	// Resize with the best quality (so it should consume some resource)
	draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	
	// Encode to `output`:      
	png.Encode(output, dst)
	log.Info("End measurement of Round : ", round, " with resizing ", totalPixels, " pixels in ", time.Since(startRound))
}