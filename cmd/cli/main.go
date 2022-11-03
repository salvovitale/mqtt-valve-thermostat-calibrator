package main

import (
	"flag"
	"time"

	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/model"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/pubsub"
)

var calibration float64
var tempSensor float64
var tempThermostat float64
var topic string
var isThermostat bool
var brokerURL string
var qos int

func init() {
	flag.Float64Var(&calibration, "c", 0.0, "calibration value")
	flag.Float64Var(&tempSensor, "t", 0.0, "temperature sensor")
	flag.Float64Var(&tempThermostat, "th", 0.0, "temperature thermostat")
	flag.StringVar(&topic, "topic", "topic/sensor", "sensor topic")
	flag.BoolVar(&isThermostat, "isth", false, "device type")
	flag.StringVar(&brokerURL, "broker", "tcp://localhost:1883", "broker url")
	flag.IntVar(&qos, "q", 0, "qos")
}

func messageHandlerThermostat(client mqtt.Client, msg mqtt.Message) {
	log.Info().Msgf("message received: %s\n", msg.Payload())
}

func main() {
	flag.Parse()

	mqttClient, err := pubsub.New(brokerURL, topic, "cli-client", 1)
	if err != nil {
		panic(err)
	}
	var payloadByte []byte
	if !isThermostat {
		payload := model.SensorPayload{Temperature: tempSensor}
		payloadByte, err = json.Marshal(payload)
	} else {
		payload := model.ThermostatPayload{LocalTemperature: tempThermostat, LocalTemperatureCalibration: calibration}
		payloadByte, err = json.Marshal(payload)
	}
	if err != nil {
		panic(err)
	}

	err = mqttClient.Publish(payloadByte)
	if err != nil {
		panic(err)
	}
	err = mqttClient.Subscribe(messageHandlerThermostat)
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Payload %s published on topic %s\n", payloadByte, topic)
	time.Sleep(1 * time.Second)
}
