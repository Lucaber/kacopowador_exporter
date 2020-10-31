package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/lucaber/kacopowador_exporter/client"
	"github.com/lucaber/kacopowador_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Host         string
	MetricPort   int
	MqttHost     string
	MqttUser     string
	MqttPassword string
	MqttName     string
	MqttInterval int64
}

var (
	conf Config
)

func main() {
	err := envconfig.Process("KACOPOWADOR", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	if conf.MqttHost != "" {
		err = setupMqtt()
		if err != nil {
			log.Fatal(err)
		}
	}

	err = setupPrometheus()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.MetricPort), nil))
}

func setupMqtt() error {
	mqttClient, err := client.MqttConnect(conf.MqttName, conf.MqttHost, conf.MqttUser, conf.MqttPassword)
	if err != nil {
		log.Fatal(err)
	}
	power, err := mqttClient.RegisterSensor("Power", "W")
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Duration(conf.MqttInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				resc := make(chan client.KacoRealtimeState)
				go client.RequestRealtimeState(conf.Host, time.Now(), resc)
				res := <-resc
				err = power.Emit(res.PRealtime)
				if err != nil {
					log.Print(err)
				}

			}
		}
	}()
	return nil
}

func setupPrometheus() error {
	kacoCollector, err := collector.NewKacoCollector(conf.Host)
	if err != nil {
		return err
	}
	err = prometheus.Register(kacoCollector)
	if err != nil {
		return err
	}

	kacoRealtimeCollector, err := collector.NewKacoRealtimeCollector(conf.Host)
	if err != nil {
		return err
	}
	err = prometheus.Register(kacoRealtimeCollector)
	if err != nil {
		return err
	}
	return nil
}
