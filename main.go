package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/alexebird/papertrail-exporter/papertrail"
	//"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

var lastEventAt = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "papertrail",
		Subsystem: "system",
		Name:      "last_event_at_seconds",
		Help:      "Last event seen at for the system.",
	},
	[]string{"group", "name"},
)

func main() {
	flag.Parse()

	prometheus.MustRegister(lastEventAt)

	go func() {
		for {
			groups := papertrail.ListGroups()

			for _, group := range groups {
				for _, system := range group.Systems {
					lastEventAt.WithLabelValues(group.Name, system.Name).Set(float64(system.LastEventAt.Unix()))
				}
			}

			log.Info("loop")
			time.Sleep(60 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
