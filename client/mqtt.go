package client

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"time"
)

type MqttClient struct {
	mqtt.Client
	Name string
	sensors []*MqttSensor
}

type MqttSensor struct {
	Config *SensorConfig
	Client *MqttClient
	StateTopic string
}

type SensorConfig struct {
	StateTopic string `json:"state_topic"`
	Device MqttDevice `json:"device"`
	DeviceClass string `json:"device_class"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	UniqueID string `json:"unique_id"`
	Name string `json:"name"`
}

type MqttDevice struct {
	Manufacturer string `json:"manufacturer"`
	Model string `json:"model"`
	Name string `json:"name"`
	Identifiers []string `json:"identifiers"`
}

func MqttConnect(name, host, username, password string) (*MqttClient, error) {
	client := &MqttClient{nil ,name, []*MqttSensor{}}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + host)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(name)
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("mqtt connected")
		err := client.AnnounceSensors()
		if err != nil {
			log.Printf("Failed to register sensor: %s", err)
		}
	})
	c := mqtt.NewClient(opts)
	client.Client = c
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	if token := client.Subscribe("homeassistant/status", 2, func(c mqtt.Client, message mqtt.Message) {
		if string(message.Payload()) == "online" {
			err := client.AnnounceSensors()
			if err != nil {
				log.Printf("Failed to register sensor: %s", err)
			}
		}
	}); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func (c * MqttClient) AnnounceSensors() error {
	for _, s := range c.sensors {
		err := s.announceSensor()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c * MqttClient) RegisterSensor(name string, unit string) (*MqttSensor, error) {
	fullName := c.Name + "_" + name
	stateTopic := fmt.Sprintf("homeassistant/sensor/%s/state", fullName)
	config := &SensorConfig{
		StateTopic: stateTopic,
		Device: MqttDevice{
			Manufacturer: "Kaco",
			Model: "Powador",
			Name: name,
			Identifiers: []string{c.Name},
		},
		DeviceClass: "power",
		UnitOfMeasurement: unit,
		UniqueID: fullName,
		Name: c.Name + " " + name,
	}

	sensor := &MqttSensor{config,
		c,
		stateTopic,
	}
	c.sensors = append(c.sensors, sensor)
	err := sensor.announceSensor()
	return sensor, err
}

func (s *MqttSensor) announceSensor() error {
	pl, err := json.Marshal(s.Config)
	if err != nil {
		return err
	}
	topic := fmt.Sprintf("homeassistant/sensor/%s/config", s.Config.UniqueID)
	token := s.Client.Publish(topic , 0, false, pl)
	for !token.WaitTimeout(3 * time.Second) {
	}
	return token.Error()
}

func (s *MqttSensor) Emit(value float64) error {
	token := s.Client.Publish(s.StateTopic, 0, false, strconv.FormatFloat(value, 'f', -1, 64 ))
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return err
	}
	return nil
}