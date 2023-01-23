package collector

import (
	"time"
	"os"
	"strconv"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"math"

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

	//if the tens is even no resizing
	//resizing if odd

	//even
	if(round/10 % 2 == 0){
		//do nothing but pushing 0 pixels metric
		ch <- prometheus.MustNewConstMetric(
			metricMeasurementPixels, prometheus.GaugeValue, float64(0), 
		)
		ch <- prometheus.MustNewConstMetric(
			metricMeasurementWrittenBytes, prometheus.GaugeValue, float64(0), 
		)
		ch <- prometheus.MustNewConstMetric(
			metricMeasurementThroughWrittenBytes, prometheus.GaugeValue, float64(0), 
		)
		log.Info("No Resize")
	} else {
		//from 0 to 9 resize the image by stepPix in both H&W
		mod := round % 10
		x := e.src.Bounds().Max.X+int(e.stepPix) * mod
		y := e.src.Bounds().Max.Y+int(e.stepPix) * mod
		e.ResizeImage(ch,x,y)
	}
}

//https://github.com/MariaLetta/free-gophers-pack
//https://stackoverflow.com/questions/22940724/go-resizing-images
func (e *Exporter) ResizeImage(ch chan<- prometheus.Metric, x , y int) {
	start := time.Now()

	//RENDERING IMAGE
	output, _ := os.Create(imageOut)
	defer output.Close()
	// Set the expected size
	dst := image.NewRGBA(image.Rect(0, 0, x, y))	
	// Resize with the best quality (so it should consume some resource)
	draw.CatmullRom.Scale(dst, dst.Rect, e.src, e.src.Bounds(), draw.Over, nil)	
	// Encode the output
	startEncoding := time.Now()
	png.Encode(output, dst)
	end := time.Now()

	//WRITING METRICS
	totalPixels := x * y
	ch <- prometheus.MustNewConstMetric(
		metricMeasurementPixels, prometheus.GaugeValue, float64(totalPixels), 
	)
	fi, err := os.Stat(imageOut)
	if err == nil {
		ch <- prometheus.MustNewConstMetric(
			metricMeasurementWrittenBytes, prometheus.GaugeValue, float64(fi.Size()), 
		)
		ch <- prometheus.MustNewConstMetric(
			metricMeasurementThroughWrittenBytes, prometheus.GaugeValue, float64(fi.Size())/(float64(end.Sub(startEncoding)/time.Microsecond)/1000000), 
		)
	} else {
		log.Error("Written Bytes error : ",err)
	}

	log.Info("Resized image is ", math.Floor(float64(totalPixels)/1000000*100)/100, " M pixels, ", 
	          math.Floor(float64(fi.Size())/(1024*1024)*100)/100, " MiB, ",
			  math.Floor((float64(fi.Size())/(float64(end.Sub(startEncoding)/time.Microsecond)/1000000))/(1024*1024)*100)/100, " MiB/s, ",
			  "overall duration ", end.Sub(start), 
			  " (resizing ",startEncoding.Sub(start)," / writing ",end.Sub(startEncoding), ")")
}