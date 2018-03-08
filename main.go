package main

import (
	"flag"
	"net/http"
	"regexp"
	"time"

	"github.com/alexebird/papertrail-exporter/papertrail"
	//"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var groupFitler = flag.String("group-filter", ".*", "A regex by which to whitelist group names.")
var systemFitler = flag.String("system-filter", ".*", "A regex by which to whitelist system names.")

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	groupFilterRegex := regexp.MustCompile(*groupFitler)
	systemFilterRegex := regexp.MustCompile(*systemFitler)
	var lastEventAt *prometheus.GaugeVec

	go func() {
		for {
			groups := papertrail.ListGroups()

			if lastEventAt != nil {
				prometheus.Unregister(lastEventAt)
			}

			lastEventAt = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: "papertrail",
					Subsystem: "system",
					Name:      "last_event_at_seconds",
					Help:      "Last event seen at for the system.",
				},
				[]string{"group", "name"},
			)
			// re-register so that old systems are purged from the metrics
			prometheus.MustRegister(lastEventAt)

			for _, group := range groups {
				for _, system := range group.Systems {
					if groupFilterRegex.MatchString(group.Name) && systemFilterRegex.MatchString(system.Name) {
						log.Debug(group.Name + " " + system.Name)
						lastEventAt.WithLabelValues(group.Name, system.Name).Set(float64(system.LastEventAt.Unix()))
					}
				}
			}

			log.Info("loop")
			time.Sleep(60 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
