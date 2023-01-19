package main

import (
	"resource-control-sample/pkg/utils"
	"resource-control-sample/pkg/collector"
	"flag"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	version           = "v0.01"
	listenAddress     = flag.String("web.listen-address", ":9905", "Address to listen on for telemetry")
	metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
)

//Readiness message
func Ready() string {
	return "Resource Control Sample is ready"
}

func Init() *collector.Exporter {
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Info(".env file absent, assume env variables are set.")
	}
  
	minPx := utils.GetInt64Env("PIXELS_MIN",1)
	maxPx := utils.GetInt64Env("PIXELS_MAX",1000000)//387*248 = 95,976
	stepPx := utils.GetInt64Env("PIXELS_STEP",1000)

	log.Info("Min Pixels   => ", minPx)
	log.Info("Max Pixels   => ", maxPx)
	log.Info("Step Pixels  => ", stepPx)

	//Registering Exporter
	exporter := collector.NewExporter(version, minPx, maxPx, stepPx)
	prometheus.MustRegister(exporter)

	return exporter
}

func main() {
	log.Info("Starting resource control sample")
	Init()

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Resource Control Sample ` + version + `</title></head>
             <body>
             <h1>` + Ready() + `'</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Info("Listening on port " + *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}


