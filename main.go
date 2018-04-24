package main

import (
	"flag"
	"net/http"
	"regexp"
	"time"

	"github.com/alexebird/papertrail-exporter/ec2"
	"github.com/alexebird/papertrail-exporter/papertrail"
	//"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var groupFitler = flag.String("group-filter", ".*", "A regex by which to whitelist group names.")
var systemFitler = flag.String("system-filter", ".*", "A regex by which to whitelist system names.")

// re-register so that old systems are purged from the metrics
func reRegisterGaugeVec(metric *prometheus.GaugeVec) *prometheus.GaugeVec {
	if metric != nil {
		prometheus.Unregister(metric)
	}

	metric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "papertrail",
			Subsystem: "system",
			Name:      "last_event_at_seconds",
			Help:      "Last event seen at for the system.",
		},
		[]string{"group", "name", "ec2_inst_present"},
	)

	prometheus.MustRegister(metric)

	return metric
}

func ec2InstancePresentLabelValue(instances map[string]bool, system papertrail.System) string {
	if instances[system.Name] {
		return "true"
	} else {
		return "false"
	}
}

func refreshPapertrailSystems(groupRegex *regexp.Regexp, sysRegex *regexp.Regexp) {
	var lastEventAtMetric *prometheus.GaugeVec

	for {
		lastEventAtMetric = reRegisterGaugeVec(lastEventAtMetric)

		systems, err := papertrail.FilterSystems(groupRegex, sysRegex)
		if err != nil {
			log.Fatal(err)
		}

		instances, err := ec2.InstanceNames()
		if err != nil {
			log.Fatal(err)
		}

		for _, system := range systems {
			log.Debug(system.GroupName + " " + system.Name)
			present := ec2InstancePresentLabelValue(instances, system)
			lastEventAtMetric.WithLabelValues(system.GroupName, system.Name, present).Set(float64(system.LastEventAt.Unix()))
		}

		log.Debug("loop")
		time.Sleep(60 * time.Second)
	}
}

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	groupFilterRegex := regexp.MustCompile(*groupFitler)
	systemFilterRegex := regexp.MustCompile(*systemFitler)

	ec2.Setup()
	go refreshPapertrailSystems(groupFilterRegex, systemFilterRegex)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
