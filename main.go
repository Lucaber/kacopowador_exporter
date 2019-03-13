package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_model/go"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	Host       string
	MetricPort int
}

type Info struct {
	KWh float64
}

type Current struct {
	Time time.Time
	Udc1 float64
	Idc1 float64
	Pdc1 float64
	Udc2 float64
	Idc2 float64
	Pdc2 float64
	Uac1 float64
	Iac1 float64
	Pac1 float64
	Uac2 float64
	Iac2 float64
	Pac2 float64
	Pdc  float64
	Pac  float64
	Tsys float64
}

type Response struct {
	Info
	Current
}

var (
	conf         Config
	metrickwh = prometheus.NewMetricWithTimestamp(
		time.Unix(0,0),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"kwh",
				"kwh",
				nil, nil,
			),
			prometheus.GaugeValue,
			0,
		))
)

func main() {
	err := envconfig.Process("KACOPOWADOR", &conf)

	if err != nil {
		log.Fatal(err.Error())
	}

	ch := make(chan Response)

	c := cron.New()
	c.AddFunc("*/15 * * * * *", func() { request(ch) })
	c.Start()


	var met io_prometheus_client.Metric
	metrickwh.Write(&met)
	met.Gauge.


	go func() {
		for {
			res := <-ch

			fmt.Println(res)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", conf.MetricPort), nil)
}

func request(ch chan Response) {
	t := time.Now()
	resp, err := http.Get(fmt.Sprintf("http://%s/%d%02d%02d.CSV", conf.Host, t.Year(), t.Month(), t.Day()))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	rows := bytes.Split(b, []byte("\r"))

	//File actually contains two csv "files"

	info, err := parseInfo(rows[:2])
	if err != nil {
		log.Fatal(err.Error())
	}

	current, err := parseCurrent(rows[2:])
	if err != nil {
		log.Fatal(err.Error())
	}

	ch <- Response{*info, *current}
}

func parseInfo(rows [][]byte) (*Info, error) {
	infoReader := csv.NewReader(bytes.NewReader(bytes.Join(rows, []byte("\n"))))
	infoReader.Comma = ';'

	infoRows, err := infoReader.ReadAll()
	if err != nil {
		return nil, err
	}
	KWh, err := strconv.ParseFloat(infoRows[1][4], 64)
	if err != nil {
		return nil, err
	}
	return &Info{KWh}, nil
}

func parseCurrent(rows [][]byte) (*Current, error) {
	infoReader := csv.NewReader(bytes.NewReader(bytes.Join(rows, []byte("\n"))))
	infoReader.Comma = ';'

	infoRows, err := infoReader.ReadAll()
	if err != nil {
		return nil, err
	}

	row := infoRows[len(infoRows)-1]

	Time, err := time.Parse("15:04:05", row[0])
	if err != nil {
		return nil, err
	}
	Udc1, err := strconv.ParseFloat(row[1], 64)
	Idc1, err := strconv.ParseFloat(row[2], 64)
	Pdc1, err := strconv.ParseFloat(row[3], 64)
	Udc2, err := strconv.ParseFloat(row[4], 64)
	Idc2, err := strconv.ParseFloat(row[5], 64)
	Pdc2, err := strconv.ParseFloat(row[6], 64)
	Uac1, err := strconv.ParseFloat(row[7], 64)
	Iac1, err := strconv.ParseFloat(row[8], 64)
	Pac1, err := strconv.ParseFloat(row[9], 64)
	Uac2, err := strconv.ParseFloat(row[10], 64)
	Iac2, err := strconv.ParseFloat(row[11], 64)
	Pac2, err := strconv.ParseFloat(row[12], 64)
	Pdc, err := strconv.ParseFloat(row[13], 64)
	Pac, err := strconv.ParseFloat(row[14], 64)
	Tsys, err := strconv.ParseFloat(row[15], 64)
	if err != nil {
		return nil, err
	}
	return &Current{Time, Udc1, Idc1, Pdc1, Udc2, Idc2, Pdc2, Uac1, Iac1, Pac1, Uac2, Iac2, Pac2, Pdc, Pac, Tsys}, nil
}
