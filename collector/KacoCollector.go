package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.mirconited.de/lusu/kacopowador_exporter/client"
	"time"
)

type KacoCollector struct {
	host string
	kwh  *prometheus.Desc
	udc1 *prometheus.Desc
	idc1 *prometheus.Desc
	pdc1 *prometheus.Desc
	udc2 *prometheus.Desc
	idc2 *prometheus.Desc
	pdc2 *prometheus.Desc
	uac1 *prometheus.Desc
	iac1 *prometheus.Desc
	pac1 *prometheus.Desc
	uac2 *prometheus.Desc
	iac2 *prometheus.Desc
	pac2 *prometheus.Desc
	pdc  *prometheus.Desc
	pac  *prometheus.Desc
	tsys *prometheus.Desc
}

func NewKacoCollector(host string) (prometheus.Collector, error) {
	return &KacoCollector{
		host: host,
		kwh: prometheus.NewDesc("kwh",
			"Generated kwh today",
			nil, nil,
		),
		udc1: prometheus.NewDesc("udc1",
			"Current udc1",
			nil, nil,
		),
		idc1: prometheus.NewDesc("idc1",
			"Current idc1",
			nil, nil,
		),
		pdc1: prometheus.NewDesc("pdc1",
			"Current pdc1",
			nil, nil,
		),
		udc2: prometheus.NewDesc("udc2",
			"Current udc2",
			nil, nil,
		),
		idc2: prometheus.NewDesc("idc2",
			"Current idc2",
			nil, nil,
		),
		pdc2: prometheus.NewDesc("pdc2",
			"Current pdc2",
			nil, nil,
		),
		uac1: prometheus.NewDesc("uac1",
			"Current uac1",
			nil, nil,
		),
		iac1: prometheus.NewDesc("iac1",
			"Current iac1",
			nil, nil,
		),
		pac1: prometheus.NewDesc("pac1",
			"Current pac1",
			nil, nil,
		),
		uac2: prometheus.NewDesc("uac2",
			"Current uac2",
			nil, nil,
		),
		iac2: prometheus.NewDesc("iac2",
			"Current iac2",
			nil, nil,
		),
		pac2: prometheus.NewDesc("pac2",
			"Current pac2",
			nil, nil,
		),
		pdc: prometheus.NewDesc("pdc",
			"Current pdc",
			nil, nil,
		),
		pac: prometheus.NewDesc("pac",
			"Current pac",
			nil, nil,
		),
		tsys: prometheus.NewDesc("tsys",
			"Current system temerature",
			nil, nil,
		),
	}, nil
}

func (collector *KacoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.kwh
	ch <- collector.udc1
	ch <- collector.idc1
	ch <- collector.pdc1
	ch <- collector.udc2
	ch <- collector.idc2
	ch <- collector.pdc2
	ch <- collector.uac1
	ch <- collector.iac1
	ch <- collector.pac1
	ch <- collector.uac2
	ch <- collector.iac2
	ch <- collector.pac2
	ch <- collector.pdc
	ch <- collector.pac
	ch <- collector.tsys

}

func (collector *KacoCollector) Collect(ch chan<- prometheus.Metric) {
	resc := make(chan client.KacoState)

	go client.RequestState(collector.host, time.Now(), resc)
	res := <-resc

	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.kwh, prometheus.CounterValue, res.KWh))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.udc1, prometheus.GaugeValue, res.Udc1))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.idc1, prometheus.GaugeValue, res.Idc1))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.pdc1, prometheus.GaugeValue, res.Pdc1))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.udc2, prometheus.GaugeValue, res.Udc2))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.idc2, prometheus.GaugeValue, res.Idc2))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.pdc2, prometheus.GaugeValue, res.Pdc2))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.uac1, prometheus.GaugeValue, res.Uac1))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.iac1, prometheus.GaugeValue, res.Iac1))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.pac1, prometheus.GaugeValue, res.Pac1))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.uac2, prometheus.GaugeValue, res.Uac2))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.iac2, prometheus.GaugeValue, res.Iac2))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.pac2, prometheus.GaugeValue, res.Pac2))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.pdc, prometheus.GaugeValue, res.Pdc))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.pac, prometheus.GaugeValue, res.Pac))
	ch <- prometheus.NewMetricWithTimestamp(res.Time, prometheus.MustNewConstMetric(collector.tsys, prometheus.GaugeValue, res.Tsys))
}
