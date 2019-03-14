package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type KacoCollector struct {
	kwh *prometheus.Desc
}

func NewKacoCollector() (prometheus.Collector, error) {
	return &KacoCollector{
		kwh: prometheus.NewDesc("kwh",
			"kwh",
			nil, nil,
		),
	}, nil
}

func (collector *KacoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.kwh
}

func (collector *KacoCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.NewMetricWithTimestamp(
		time.Unix(200000, 0),
		prometheus.MustNewConstMetric(collector.kwh, prometheus.GaugeValue, 10),
	)
}
