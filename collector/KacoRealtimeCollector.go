package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/lucaber/kacopowador_exporter/client"
	"time"
)

type KacoRealtimeCollector struct {
	host string
	p_realtime *prometheus.Desc
}

func NewKacoRealtimeCollector(host string) (prometheus.Collector, error) {
	return &KacoRealtimeCollector{
		host: host,
		p_realtime: prometheus.NewDesc("p_realtime",
			"Current realtime power in watt",
			nil, nil,
		),
	}, nil
}

func (collector *KacoRealtimeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.p_realtime

}

func (collector *KacoRealtimeCollector) Collect(ch chan<- prometheus.Metric) {
	resc := make(chan client.KacoRealtimeState)

	go client.RequestRealtimeState(collector.host, time.Now(), resc)
	res := <-resc

	ch <- prometheus.MustNewConstMetric(collector.p_realtime, prometheus.GaugeValue, res.PRealtime)
}
