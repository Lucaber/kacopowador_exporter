package client

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"time"
)

type MqttClient struct {
	mqtt.Client
	Name string
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
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + host)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(name)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return &MqttClient{client ,name}, nil
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
	pl, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	topic := fmt.Sprintf("homeassistant/sensor/%s/config", fullName)
	token := c.Client.Publish(topic , 0, false, pl)
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return &MqttSensor{config,
		c,
		stateTopic,
	}, nil
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