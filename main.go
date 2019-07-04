package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.mirconited.de/lusu/kacopowador_exporter/collector"
	"log"
	"net/http"
)

type Config struct {
	Host       string
	MetricPort int
}

var (
	conf Config
)

func main() {
	err := envconfig.Process("KACOPOWADOR", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	kacoCollector, err := collector.NewKacoCollector(conf.Host)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = prometheus.Register(kacoCollector)
	if err != nil {
		log.Fatal(err.Error())
	}

	kacoRealtimeCollector, err := collector.NewKacoRealtimeCollector(conf.Host)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = prometheus.Register(kacoRealtimeCollector)
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.MetricPort), nil))
}
